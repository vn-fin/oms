package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

type BasketUpdateRequest struct {
	Name        *string                     `json:"name"`
	Description *string                     `json:"description"`
	Info        *[]models.BasketInfo        `json:"info"`
	HedgeConfig *[]models.BasketHedgeConfig `json:"hedge_config"`
}

// BasketUpdate
// @Summary Update basket
// @Description Update an existing basket. All fields are optional, only provided fields will be updated.
// @Tags basket
// @Accept json
// @Produce json
// @Param id path string true "Basket ID"
// @Param request body BasketUpdateRequest true "Basket update payload"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 404 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/{id} [put]
func BasketUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return api.Response().BadRequest("basket id is required").Send(c)
	}

	// Get userID from context
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	// Parse request body
	var req BasketUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return api.Response().BadRequest("invalid JSON body").Send(c)
	}

	// Fetch existing basket
	var existingBasket models.Basket
	querySelect := `
		SELECT id, name, description, info, hedge_config, created_by, updated_by, created_at, updated_at, status
		FROM execution.baskets
		WHERE id = ?
	`
	_, err := db.Postgres.Query(&existingBasket, querySelect, id)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Check if basket exists
	if existingBasket.ID == "" {
		return api.Response().NotFound("basket not found").Send(c)
	}

	// Apply partial updates
	if req.Name != nil {
		trimmedName := strings.TrimSpace(*req.Name)
		if trimmedName == "" {
			return api.Response().BadRequest("name cannot be empty").Send(c)
		}
		existingBasket.Name = trimmedName
	}

	if req.Description != nil {
		existingBasket.Description = strings.TrimSpace(*req.Description)
	}

	if req.Info != nil {
		existingBasket.Info = *req.Info
	}

	if req.HedgeConfig != nil {
		existingBasket.HedgeConfig = *req.HedgeConfig
	}

	// Update timestamp and user
	now := time.Now().UTC()
	existingBasket.UpdatedAt = now
	existingBasket.UpdatedBy = userID

	// Update in database
	queryUpdate := `
		UPDATE execution.baskets
		SET name = ?, description = ?, info = ?, hedge_config = ?, updated_by = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := db.Postgres.Exec(queryUpdate,
		existingBasket.Name,
		existingBasket.Description,
		existingBasket.Info,
		existingBasket.HedgeConfig,
		existingBasket.UpdatedBy,
		existingBasket.UpdatedAt,
		id,
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Check if any row was affected
	if result.RowsAffected() == 0 {
		return api.Response().NotFound("basket not found").Send(c)
	}

	return api.Response().
		Data(existingBasket).
		Message("Basket updated successfully").
		Send(c)
}
