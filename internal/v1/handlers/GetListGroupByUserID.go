package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

// GetListGroupByUserID gets all credential groups for a specific user
// @Summary Get credential groups by user ID
// @Description Get all credential groups belonging to the authenticated user
// @Tags Groups
// @Produce json
// @Success 200 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credential-groups [get]
func GetListGroupByUserID(c *fiber.Ctx) error {
	// Get userID from context (set by AuthMiddleware)
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	// Query all credential groups for the user
	var groups []models.CredentialGroup
	query := `
		SELECT id, name, user_id, created_at, updated_at, status
		FROM users.credential_groups
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	_, err := db.Postgres.Query(&groups, query, userID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Return empty array if no groups found
	if groups == nil {
		groups = []models.CredentialGroup{}
	}

	return api.Response().
		Data(groups).
		Message("Credential groups retrieved successfully").
		Send(c)
}
