package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

// BasketList gets all baskets
// @Summary Get all baskets
// @Description Get all baskets from database
// @Tags baskets
// @Produce json
// @Success 200 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets [get]
func BasketList(c *fiber.Ctx) error {
	// Get userID from context (set by AuthMiddleware)
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	// Query all baskets
	var baskets []models.Basket
	query := `
		SELECT id, name, description, info, hedge_config, created_by, updated_by, created_at, updated_at, status
		FROM execution.baskets
		ORDER BY created_at DESC
	`
	_, err := db.Postgres.Query(&baskets, query)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Return empty array if no baskets found
	if baskets == nil {
		baskets = []models.Basket{}
	}

	return api.Response().
		Data(baskets).
		Message("Baskets retrieved successfully").
		Send(c)
}
