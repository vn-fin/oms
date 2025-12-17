package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/typing"
)

// CredentialDelete soft deletes a credential by ID
// @Summary Delete credential
// @Description Soft delete a credential by updating status to disabled (Admin only)
// @Tags credentials
// @Produce json
// @Param credential_id path string true "Credential ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 403 {object} models.DefaultResponseModel
// @Failure 404 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credentials/{credential_id} [delete]
func CredentialDelete(c *fiber.Ctx) error {
	// Check admin permission
	userEmail := api.GetUserEmail(c)
	if userEmail != config.AdminEmails {
		return api.Response().Forbidden("You are not allowed to delete credentials.").Send(c)
	}

	credentialID := strings.TrimSpace(c.Params("credential_id"))
	if credentialID == "" {
		return api.Response().BadRequest("credential_id is required").Send(c)
	}

	now := time.Now().UTC()

	// Soft delete: update status to disabled
	query := `
		UPDATE execution.login_credentials
		SET status = ?, updated_at = ?
		WHERE credential_id = ? AND status != ?
	`
	result, err := db.Postgres.Exec(query,
		typing.StatusDisabled,
		now,
		credentialID,
		typing.StatusDisabled, // Don't update if already disabled
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Check if any row was affected
	if result.RowsAffected() == 0 {
		return api.Response().NotFound("credential not found or already deleted").Send(c)
	}

	return api.Response().
		Data(fiber.Map{
			"credential_id": credentialID,
			"status":        typing.StatusDisabled,
			"updated_at":    now,
		}).
		Message("Credential deleted successfully").
		Send(c)
}
