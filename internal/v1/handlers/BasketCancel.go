package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
)

// BasketExecutionCancel
// @Summary Cancel basket execution session
// @Description Cancel a specific execution session by setting status to canceled
// @Accept json
// @Produce json
// @Param basket_id path string true "Basket ID"
// @Param execution_id path string true "Execution ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 404 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/:basket_id/executions/{execution_id}/cancel [post]
func BasketExecutionCancel(c *fiber.Ctx) error {
	basketID := strings.TrimSpace(c.Params("basket_id"))
	if basketID == "" {
		return api.Response().BadRequest("basket_id is required").Send(c)
	}

	executionID := strings.TrimSpace(c.Params("execution_id"))
	if executionID == "" {
		return api.Response().BadRequest("execution_id is required").Send(c)
	}
	var session models.BasketExecuteSession
	checkQuery := `
		SELECT id, basket_id, weight, price_level, action_type, future_size,
		       estimated_cash, matched_cash, order_status, created_by, created_at, updated_at
		FROM execution.basket_execute_sessions
		WHERE id = $1 AND basket_id = $2
	`
	_, err := db.Postgres.Query(&session, checkQuery, executionID, basketID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Check if session found

	if session.OrderStatus == typing.OrderStatusCanceled {
		return api.Response().BadRequest("execution session already canceled").Send(c)
	}

	now := time.Now().UTC()

	updateQuery := `
		UPDATE execution.basket_execute_sessions
		SET order_status = $1, updated_at = $2
		WHERE id = $3 AND basket_id = $4
	`
	_, err = db.Postgres.Exec(updateQuery, typing.OrderStatusCanceled, now, executionID, basketID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Update session object for response
	session.OrderStatus = typing.OrderStatusCanceled
	session.UpdatedAt = now

	return api.Response().
		Data(session).
		Message("Execution session canceled successfully").
		Send(c)
}
