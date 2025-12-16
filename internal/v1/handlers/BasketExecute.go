package handlers

import (
	"io"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
)

// BasketExecute
// @Summary Execute basket
// @Description Execute a basket with specified parameters
// @Accept json
// @Produce json
// @Param basket_id path string true "Basket ID"
// @Param request body models.BasketExecuteRequest true "Basket execute payload"
// @Success 201 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/{basket_id}/execute [post]
func BasketExecute(c *fiber.Ctx) error {
	// Get basket_id from URL parameter
	basketID := strings.TrimSpace(c.Params("basket_id"))
	if basketID == "" {
		return api.Response().BadRequest("basket_id is required").Send(c)
	}

	var req models.BasketExecuteRequest
	if err := c.BodyParser(&req); err != nil {
		return api.Response().BadRequest("invalid JSON body").Send(c)
	}

	if !req.ActionType.Valid() {
		return api.Response().BadRequest("invalid action_type, must be 'B' or 'S'").Send(c)
	}

	userID := api.GetUserID(c)
	now := time.Now().UTC()

	// Get basket info and hedge_config
	log.Info().Msgf("Querying basket with ID: %s", basketID)
	var basket models.Basket
	basketQuery := `
		SELECT id, name, description, info, hedge_config, created_by, updated_by, created_at, updated_at, status
		FROM execution.baskets
		WHERE id = ?
	`
	_, err := db.Postgres.Query(&basket, basketQuery, basketID)
	log.Info().Msgf("Query error (if any): %v, EOF check: %v", err, err == io.EOF)
	if err != nil && err != io.EOF {
		log.Error().Err(err).Msg("Error querying basket")
		return api.Response().InternalError(err).Send(c)
	}
	log.Info().Msgf("Basket query result - ID: '%s', Name: '%s', Info count: %d, HedgeConfig count: %d",
		basket.ID, basket.Name, len(basket.Info), len(basket.HedgeConfig))
	if basket.ID == "" {
		return api.Response().NotFound("basket not found").Send(c)
	}

	// Calculate estimated_cash = sum(cash * size)
	// cash from info, size from hedge_config
	estimatedCash := 0.0
	for _, hedgeItem := range basket.HedgeConfig {
		// Find matching symbol in info
		for _, infoItem := range basket.Info {
			if infoItem.Symbol == hedgeItem.Symbol {
				estimatedCash += infoItem.Cash * hedgeItem.Size
				break
			}
		}
	}

	// Create basket execute session
	session := models.BasketExecuteSession{
		ID:            uuid.NewString(),
		BasketID:      basketID,
		Weight:        req.Weight,
		PriceLevel:    req.PriceLevel,
		ActionType:    req.ActionType,
		FutureSize:    req.FutureSize,
		EstimatedCash: estimatedCash,
		MatchedCash:   0,
		OrderStatus:   typing.OrderStatusPending,
		CreatedBy:     userID,
		CreatedAt:     now,
	}

	// Insert into database using raw SQL query
	query := `
		INSERT INTO execution.basket_execute_sessions
		(id, basket_id, weight, price_level, action_type, future_size, estimated_cash, matched_cash, order_status, created_by, created_at)
		VALUES (?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = db.Postgres.Exec(query,
		session.ID,
		session.BasketID,
		session.Weight,
		session.PriceLevel,
		session.ActionType,
		session.FutureSize,
		session.EstimatedCash,
		session.MatchedCash,
		session.OrderStatus,
		session.CreatedBy,
		session.CreatedAt,
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	return api.Response().
		Status(fiber.StatusCreated).
		Data(session).
		Message("Basket execution session created successfully").
		Send(c)
}
