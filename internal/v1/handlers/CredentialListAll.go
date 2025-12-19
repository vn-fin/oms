package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

// CredentialListAll
// @Summary Get all login credentials
// @Description Get all login credentials (Admin only)
// @Tags Admin
// @Produce json
// @Success 200 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 403 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credentials [get]
func CredentialListAll(c *fiber.Ctx) error {
	// Check admin permission
	userEmail := api.GetUserEmail(c)
	if userEmail != config.AdminEmails {
		return api.Response().Forbidden("You are not allowed to view all credentials.").Send(c)
	}

	// Query all credentials
	var credentials []models.Credential
	query := `
		SELECT id as credential_id, name, description, info, created_at, updated_at, status
		FROM users.login_credentials
		ORDER BY created_at DESC
	`
	_, err := db.Postgres.Query(&credentials, query)
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
