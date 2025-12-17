package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

// BasketExecuteSessionsList
// @Summary List basket execute sessions
// @Description Get all execute sessions for a specific basket
// @Tags Execution
// @Accept json
// @Produce json
// @Param basket_id path string true "Basket ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/{basket_id}/execute-sessions [get]
func BasketExecuteSessionsList(c *fiber.Ctx) error {
	basketID := strings.TrimSpace(c.Params("basket_id"))
	if basketID == "" {
		return api.Response().BadRequest("basket_id is required").Send(c)
	}

	var sessions []models.BasketExecuteSession
	query := `
		SELECT id, basket_id, weight, price_level, action_type, future_size,
		       estimated_cash, matched_cash, order_status, created_by, created_at
		FROM execution.basket_execute_sessions
		WHERE basket_id = ?
		ORDER BY created_at DESC
	`
	_, err := db.Postgres.Query(&sessions, query, basketID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Return empty array if no sessions found
	if sessions == nil {
		sessions = []models.BasketExecuteSession{}
	}

	return api.Response().
		Data(sessions).
		Message("Basket execute sessions retrieved successfully").
		Send(c)
}
