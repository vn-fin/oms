package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/typing"
)

// BasketCancel cancels a basket by ID
// @Summary Cancel basket
// @Description Cancel a basket by updating status to disabled
// @Tags baskets
// @Produce json
// @Param id path string true "Basket ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 404 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/{id}/cancel [post]
func BasketCancel(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return api.Response().BadRequest("basket id is required").Send(c)
	}

	// Get userID from context (set by AuthMiddleware)
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	now := time.Now().UTC()

	// Cancel: update status to disabled
	query := `
		UPDATE execution.baskets
		SET status = ?, updated_by = ?, updated_at = ?
		WHERE id = ? AND status = ?
	`
	result, err := db.Postgres.Exec(query,
		typing.RecordStatusDisabled,
		userID,
		now,
		id,
		typing.RecordStatusEnabled, // Only cancel if currently enabled
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Check if any row was affected
	if result.RowsAffected() == 0 {
		return api.Response().NotFound("basket not found or not in enabled status").Send(c)
	}

	return api.Response().
		Data(fiber.Map{
			"id":         id,
			"status":     typing.RecordStatusDisabled,
			"updated_by": userID,
			"updated_at": now,
		}).
		Message("Basket cancelled successfully").
		Send(c)
}
