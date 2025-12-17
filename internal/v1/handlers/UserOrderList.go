package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

// UserOrderListBySession gets all user orders for a specific session
// @Summary Get user orders by execution session ID
// @Description Get all user orders belonging to a specific execution session
// @Tags Execution
// @Produce json
// @Param basket_id path string true "Basket ID"
// @Param execution_id path string true "Execution Session ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/{basket_id}/execute/{execution_id}/orders [get]
func UserOrderListBySession(c *fiber.Ctx) error {
	sessionID := strings.TrimSpace(c.Params("execution_id"))
	if sessionID == "" {
		return api.Response().BadRequest("session_id is required").Send(c)
	}

	// Get userID from context
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	// Query all user orders for the session
	var orders []models.UserOrder
	query := `
		SELECT id, credential_id, session_id, symbol, symbol_type, side, order_price, matched_price, quantity, filled_qty, remaining_qty, status, created_at, updated_at
		FROM execution.user_orders
		WHERE session_id = ?
		ORDER BY created_at DESC
	`
	_, err := db.Postgres.Query(&orders, query, sessionID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Return empty array if no orders found
	if orders == nil {
		orders = []models.UserOrder{}
	}

	return api.Response().
		Data(orders).
		Message("User orders retrieved successfully").
		Send(c)
}
