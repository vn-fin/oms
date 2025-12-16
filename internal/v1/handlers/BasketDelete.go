package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/typing"
)

// BasketDelete soft deletes a basket by ID
// @Summary Delete basket
// @Description Soft delete a basket by updating status to removed
// @Tags baskets
// @Produce json
// @Param id path string true "Basket ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 404 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/{id} [delete]
func BasketDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return api.Response().BadRequest("basket id is required").Send(c)
	}

	// Get userID from context (set by AuthMiddleware)
	userID := api.GetUserID(c)

	now := time.Now().UTC()

	// Soft delete: update status to removed
	query := `
		UPDATE execution.baskets
		SET status = ?, updated_by = ?, updated_at = ?
		WHERE id = ? AND status != ?
	`
	result, err := db.Postgres.Exec(query,
		typing.RecordStatusRemoved,
		userID,
		now,
		id,
		typing.RecordStatusRemoved, // Don't update if already removed
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Check if any row was affected
	if result.RowsAffected() == 0 {
		return api.Response().NotFound("basket not found or already deleted").Send(c)
	}

	return api.Response().
		Data(fiber.Map{
			"id":         id,
			"status":     typing.RecordStatusRemoved,
			"updated_by": userID,
			"updated_at": now,
		}).
		Message("Basket deleted successfully").
		Send(c)
}
