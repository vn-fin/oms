package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

// GetListGroupByUserID gets all credential groups for a specific user
// @Summary Get credential groups by user ID
// @Description Get all credential groups belonging to the authenticated user
// @Tags Admin
// @Produce json
// @Success 200 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credential-groups [get]
func GetListGroup(c *fiber.Ctx) error {
	// Get userID from context (set by AuthMiddleware)
	userEmail := api.GetUserEmail(c)
	if userEmail != config.AdminEmails {
		return api.Response().Forbidden("You are not allowed to view all credentials.").Send(c)
	}
	// Query all credential groups for the user
	var groups []models.CredentialGroup
	query := `
		SELECT 
			cg.id, 
			cg.name, 
			cg.user_id, 
			cg.created_at, 
			cg.updated_at, 
			cg.status,
			COALESCE(SUM(lcgd.cash_limit), 0) as total_cash_limit,
			COALESCE(SUM(lcgd.balance), 0) as total_balance
		FROM users.credential_groups cg
		LEFT JOIN users.login_credential_group_details lcgd ON cg.id = lcgd.credential_group_id
		GROUP BY cg.id, cg.name, cg.user_id, cg.created_at, cg.updated_at, cg.status
		ORDER BY cg.created_at DESC
	`
	_, err := db.Postgres.Query(&groups, query)
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
