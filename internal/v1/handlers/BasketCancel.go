package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
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
// @Router /oms/v1/baskets/{basket_id}/executions/{execution_id}/cancel [post]
func BasketExecutionCancel(c *fiber.Ctx) error {
	basketID := strings.TrimSpace(c.Params("basket_id"))
	if basketID == "" {
		return api.Response().BadRequest("basket_id is required").Send(c)
	}

	executionID := strings.TrimSpace(c.Params("execution_id"))
	if executionID == "" {
		return api.Response().BadRequest("execution_id is required").Send(c)
	}

	updateQuery := `
          UPDATE execution.basket_execute_sessions
          SET order_status = ?
          WHERE id = ? AND basket_id = ?
      `
	_, err := db.Postgres.Exec(updateQuery, typing.OrderStatusCanceled, executionID, basketID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	return api.Response().
		Message("Execution session canceled successfully").
		Send(c)
}
