package handlers

import (
	"io"
	"math"
	"math/rand"
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
// @Tags Execution
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

	if req.CredentialID == "" {
		return api.Response().BadRequest("credential_id is required").Send(c)
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
		EstimatedCash: 100000000, //Gia su :v (chua tinh)
		MatchedCash:   0,
		OrderStatus:   typing.OrderStatusCreated,
		CreatedBy:     userID,
		CreatedAt:     now,
	}

	// Insert session into database
	sessionQuery := `
		INSERT INTO execution.basket_execute_sessions
		(id, basket_id, weight, price_level, action_type, future_size, estimated_cash, matched_cash, order_status, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = db.Postgres.Exec(sessionQuery,
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

	// Create user_orders for each symbol in hedge_config
	var userOrders []models.UserOrder
	for _, infoItem := range basket.Info {

		// Find matching symbol in info to get cash
		var cash float64
		for j, infoItem := range basket.Info {
			log.Info().Msgf("Info[%d]: Symbol=%s, Cash=%.2f", j, infoItem.Symbol, infoItem.Cash)
			if infoItem.Symbol == infoItem.Symbol {
				cash = infoItem.Cash
				break
			}
		}

		if cash == 0 {
			log.Info().Msgf("Skipping symbol %s because cash is 0", infoItem.Symbol)
			continue
		}

		matchedPrice := 10.0 + rand.Float64()*90.0
		matchedPrice = math.Round(matchedPrice*100) / 100 // Round to 2 decimal places

		quantity := math.Floor(cash / matchedPrice)
		if quantity == 0 {
			continue
		}

		userOrder := models.UserOrder{
			ID:           uuid.NewString(),
			CredentialID: req.CredentialID,
			SessionID:    session.ID,
			Symbol:       infoItem.Symbol,
			SymbolType:   "VnStock",
			Side:         req.ActionType,
			OrderPrice:   matchedPrice, // Use matched_price as order_price for mock
			MatchedPrice: matchedPrice,
			Quantity:     quantity,
			FilledQty:    0,
			RemainingQty: quantity,
			Status:       typing.OrderStatusCreated,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// Insert user_order into database
		orderQuery := `
			INSERT INTO execution.user_orders
			(id, credential_id, session_id, symbol, symbol_type, side, order_price, matched_price, quantity, filled_qty, remaining_qty, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		log.Info().Msgf("Inserting user_order - ID: %s, Symbol: %s, SymbolType: %s, Quantity: %.0f",
			userOrder.ID, userOrder.Symbol, userOrder.SymbolType, userOrder.Quantity)
		_, err = db.Postgres.Exec(orderQuery,
			userOrder.ID,
			userOrder.CredentialID,
			userOrder.SessionID,
			userOrder.Symbol,
			userOrder.SymbolType,
			userOrder.Side,
			userOrder.OrderPrice,
			userOrder.MatchedPrice,
			userOrder.Quantity,
			userOrder.FilledQty,
			userOrder.RemainingQty,
			userOrder.Status,
			userOrder.CreatedAt,
			userOrder.UpdatedAt,
		)
		if err != nil {
			return api.Response().InternalError(err).Send(c)
		}

		userOrders = append(userOrders, userOrder)
	}

	return api.Response().
		Status(fiber.StatusCreated).
		Data(fiber.Map{
			"session":     session,
			"user_orders": userOrders,
		}).
		Message("Basket execution session created successfully").
		Send(c)
}
