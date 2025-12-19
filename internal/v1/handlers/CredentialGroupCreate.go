package handlers

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
)

type CredentialItem struct {
	CredentialID string  `json:"credential_id"`
	CashLimit    float64 `json:"cash_limit"`
}

type CredentialGroupRequest struct {
	Name        string           `json:"group_name"`
	Email       string           `json:"email"`
	Credentials []CredentialItem `json:"credentials"`
}

// CredentialGroupCreate
// @Summary Create credential group
// @Description Create a new credential group with multiple credentials (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body CredentialGroupRequest true "create credential group payload"
// @Success 201 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 403 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credential-groups [post]
func CredentialGroupCreate(c *fiber.Ctx) error {
	// Check admin permission
	userEmail := api.GetUserEmail(c)
	if userEmail != config.AdminEmails {
		return api.Response().Forbidden("You are not allowed to create credential groups.").Send(c)
	}

	var req CredentialGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return api.Response().BadRequest("invalid JSON body").Send(c)
	}

	// Validate email
	if req.Email == "" {
		return api.Response().BadRequest("email is required").Send(c)
	}

	// Validate credentials
	if len(req.Credentials) == 0 {
		return api.Response().BadRequest("at least one credential is required").Send(c)
	}

	// Query user_id from xno_ai_data.users.users
	var userID string
	getUserQuery := `SELECT user_id FROM users.users WHERE email = ?`
	_, err := db.PostgresUserDB.QueryOne(pg.Scan(&userID), getUserQuery, req.Email)
	if err != nil {
		return api.Response().BadRequest("user not found with this email").Send(c)
	}

	now := time.Now().UTC()
	credentialGroup := models.CredentialGroup{
		ID:        uuid.NewString(),
		Name:      req.Name,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
		Status:    typing.StatusActive,
	}

	// Insert credential group
	groupQuery := `
		INSERT INTO users.credential_groups (id, name, user_id, created_at, updated_at, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = db.Postgres.Exec(groupQuery,
		credentialGroup.ID,
		credentialGroup.Name,
		credentialGroup.UserID,
		credentialGroup.CreatedAt,
		credentialGroup.UpdatedAt,
		credentialGroup.Status,
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Insert credential group details for each credential
	var credentialDetails []models.CredentialGroupDetails
	for _, cred := range req.Credentials {
		detail := models.CredentialGroupDetails{
			ID:                uuid.NewString(),
			CredentialID:      cred.CredentialID,
			CredentialGroupID: credentialGroup.ID,
			CashLimit:         cred.CashLimit,
			UpdatedAt:         now,
			Status:            typing.StatusActive,
		}

		detailQuery := `
			INSERT INTO users.login_credential_group_details (id, credential_id, credential_group_id, cash_limit, updated_at, status)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err = db.Postgres.Exec(detailQuery,
			detail.ID,
			detail.CredentialID,
			detail.CredentialGroupID,
			detail.CashLimit,
			detail.UpdatedAt,
			detail.Status,
		)
		if err != nil {
			return api.Response().InternalError(err).Send(c)
		}

		credentialDetails = append(credentialDetails, detail)
	}

	return api.Response().
		Status(fiber.StatusCreated).
		Data(fiber.Map{
			"group":       credentialGroup,
			"credentials": credentialDetails,
		}).
		Message("Credential group created successfully").
		Send(c)
}
