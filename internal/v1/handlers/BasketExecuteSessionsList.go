package handlers

import (
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
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/execute-sessions [get]
func BasketExecuteSessionsList(c *fiber.Ctx) error {
	userId := api.GetUserID(c)
	var sessions []models.BasketExecuteSession
	query := `
		SELECT
    bs.id,
    bs.basket_id,
    b.name AS basket_name,
    bs.weight,
    bs.price_level,
    bs.action_type,
    bs.future_size,
    bs.estimated_cash,
    bs.matched_cash,
    bs.order_status,
    bs.created_by,
    bs.created_at
FROM execution.basket_execute_sessions bs
         LEFT JOIN execution.baskets b
                   ON bs.basket_id = b.id
WHERE bs.created_by = ?
ORDER BY bs.created_at DESC;

	`
	_, err := db.Postgres.Query(&sessions, query, userId)
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
