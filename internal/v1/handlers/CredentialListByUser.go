package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

// CredentialListByGroup gets all credentials for a specific group
// @Summary Get all credentials by group ID
// @Description Get all credentials belonging to a specific credential group
// @Tags groups
// @Produce json
// @Param group_id path string true "Credential Group ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credential-groups/{group_id}/credentials [get]
func CredentialListByGroup(c *fiber.Ctx) error {
	// Get userID from context
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	// Get group_id from path
	groupID := strings.TrimSpace(c.Params("group_id"))
	if groupID == "" {
		return api.Response().BadRequest("group_id is required").Send(c)
	}

	// Query all credentials for the group (with user_id check for security)
	var credentials []models.Credential
	query := `
		SELECT DISTINCT
			lc.id as credential_id,
			lc.name,
			lc.description,
			lc.info,
			lc.created_at,
			lc.updated_at,
			lc.status
		FROM users.credential_groups cg
		INNER JOIN users.login_credential_group_details lcgd ON cg.id = lcgd.credential_group_id
		INNER JOIN users.login_credentials lc ON lcgd.credential_id = lc.id
		WHERE cg.user_id = ? AND cg.id = ?
		ORDER BY lc.created_at DESC
	`
	_, err := db.Postgres.Query(&credentials, query, userID, groupID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Return empty array if no credentials found
	if credentials == nil {
		credentials = []models.Credential{}
	}

	return api.Response().
		Data(credentials).
		Message("Credentials retrieved successfully").
		Send(c)
}
