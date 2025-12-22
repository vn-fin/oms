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

	basketID := strings.TrimSpace(c.Params("basket_id"))
	if basketID == "" {
		return api.Response().BadRequest("basket_id is required").Send(c)
	}

	// Query all user orders for the session with join to verify user ownership
	var orders []models.UserOrder
	query := `
		SELECT
			uo.id,
			uo.credential_id,
			uo.session_id,
			uo.symbol,
			uo.symbol_type,
			uo.side,
			uo.order_type,
			uo.order_price,
			uo.matched_price,
			uo.quantity,
			uo.filled_qty,
			uo.remaining_qty,
			uo.status,
			uo.created_at,
			uo.updated_at
		FROM users.credential_groups cg
		INNER JOIN execution.baskets b ON cg.id = b.group_id
		INNER JOIN execution.basket_execute_sessions bes ON b.id = bes.basket_id
		INNER JOIN execution.user_orders uo ON bes.id = uo.session_id
		WHERE cg.user_id = ?
			AND bes.basket_id = ?
			AND bes.id = ?
		ORDER BY uo.created_at DESC
	`
	_, err := db.Postgres.Query(&orders, query, userID, basketID, sessionID)
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
