package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

// BasketDetail gets a basket by ID
// @Summary Get basket detail
// @Description Get a basket by ID
// @Tags baskets
// @Produce json
// @Param id path string true "Basket ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 404 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/{id} [get]
func BasketDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return api.Response().BadRequest("basket id is required").Send(c)
	}

	// Get userID from context (set by AuthMiddleware)
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	// Query basket by ID
	var basket models.Basket
	query := `
		SELECT id, name, description, info, hedge_config, created_by, updated_by, created_at, updated_at, status
		FROM execution.baskets
		WHERE id = ?
	`
	_, err := db.Postgres.Query(&basket, query, id)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Check if basket found
	if basket.ID == "" {
		return api.Response().NotFound("basket not found").Send(c)
	}

	return api.Response().
		Data(basket).
		Message("Basket retrieved successfully").
		Send(c)
}
