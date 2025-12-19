package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
)

// CredentialListAll
// @Summary Get all login credentials
// @Description Get all login credentials with their groups (Admin only)
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

	// Define row struct for LEFT JOIN query
	type credentialRow struct {
		CredentialID string                  `pg:"credential_id"`
		Name         string                  `pg:"name"`
		Description  string                  `pg:"description"`
		Info         []models.CredentialInfo `pg:"info"`
		CreatedAt    time.Time               `pg:"created_at"`
		UpdatedAt    time.Time               `pg:"updated_at"`
		Status       typing.AccountStatus    `pg:"status"`
		GroupID      *string                 `pg:"group_id"`
		GroupName    *string                 `pg:"group_name"`
		CashLimit    *float64                `pg:"cash_limit"`
	}

	// Query all credentials with groups (LEFT JOIN)
	var rows []credentialRow
	query := `
		SELECT
			lc.id as credential_id,
			lc.name,
			lc.description,
			lc.info,
			lc.created_at,
			lc.updated_at,
			lc.status,
			cg.id as group_id,
			cg.name as group_name,
			lcgd.cash_limit
		FROM users.login_credentials lc
		LEFT JOIN users.login_credential_group_details lcgd ON lcgd.credential_id = lc.id
		LEFT JOIN users.credential_groups cg ON cg.id = lcgd.credential_group_id
		ORDER BY lc.created_at DESC, cg.name
	`
	_, err := db.Postgres.Query(&rows, query)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Group rows by credential_id
	credentialMap := make(map[string]*models.CredentialWithGroups)
	var order []string

	for _, row := range rows {
		cred, exists := credentialMap[row.CredentialID]
		if !exists {
			order = append(order, row.CredentialID)
			cred = &models.CredentialWithGroups{
				CredentialID: row.CredentialID,
				Name:         row.Name,
				Description:  row.Description,
				Info:         row.Info,
				Groups:       []models.CredentialGroupRef{},
				CreatedAt:    row.CreatedAt,
				UpdatedAt:    row.UpdatedAt,
				Status:       row.Status,
			}
			credentialMap[row.CredentialID] = cred
		}

		// Add group if exists
		if row.GroupID != nil && row.GroupName != nil {
			var limit float64
			if row.CashLimit != nil {
				limit = *row.CashLimit
			}
			cred.Groups = append(cred.Groups, models.CredentialGroupRef{
				GroupID:   *row.GroupID,
				Name:      *row.GroupName,
				CashLimit: limit,
			})
		}
	}

	// Build result array maintaining order
	credentials := make([]models.CredentialWithGroups, 0, len(order))
	for _, id := range order {
		credentials = append(credentials, *credentialMap[id])
	}

	return api.Response().
		Data(credentials).
		Message("Credentials retrieved successfully").
		Send(c)
}
