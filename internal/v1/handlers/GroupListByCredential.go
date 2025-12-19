package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

// GroupListByCredential gets all groups that a credential is assigned to
// @Summary Get all groups by credential ID
// @Description Get all groups that a specific credential is assigned to (Admin only)
// @Tags Admin
// @Produce json
// @Param credential_id path string true "Credential ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 403 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credentials/{credential_id}/groups [get]
func GroupListByCredential(c *fiber.Ctx) error {
	// Check admin permission
	userEmail := api.GetUserEmail(c)
	if userEmail != config.AdminEmails {
		return api.Response().Forbidden("You are not allowed to view credential groups.").Send(c)
	}

	// Get credential_id from path
	credentialID := strings.TrimSpace(c.Params("credential_id"))
	if credentialID == "" {
		return api.Response().BadRequest("credential_id is required").Send(c)
	}

	// Query all groups for the credential
	var groups []models.CredentialGroup
	query := `
		SELECT DISTINCT
			cg.id,
			cg.name,
			cg.user_id,
			cg.created_at,
			cg.updated_at,
			cg.status
		FROM users.credential_groups cg
		INNER JOIN users.login_credential_group_details lcgd ON cg.id = lcgd.credential_group_id
		WHERE lcgd.credential_id = ?
		ORDER BY cg.created_at DESC
	`
	_, err := db.Postgres.Query(&groups, query, credentialID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Return empty array if no groups found
	if groups == nil {
		groups = []models.CredentialGroup{}
	}

	return api.Response().
		Data(groups).
		Message("Groups retrieved successfully").
		Send(c)
}
