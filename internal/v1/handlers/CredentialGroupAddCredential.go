package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
)

type CredentialGroupAddCredentialRequest struct {
	CredentialID      string  `json:"credential_id"`
	CredentialGroupID string  `json:"credential_group_id"`
	CashLimit         float64 `json:"cash_limit"`
}

// CredentialGroupAddCredential
// @Summary Add credential to group
// @Description Add a credential to a credential group (Admin only)
// @Accept json
// @Produce json
// @Param request body CredentialGroupAddCredentialRequest true "Add credential to group payload"
// @Success 201 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 403 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credential-groups/add-credential [post]
func CredentialGroupAddCredential(c *fiber.Ctx) error {
	// Check admin permission
	userEmail := api.GetUserEmail(c)
	if userEmail != config.AdminEmails {
		return api.Response().Forbidden("You are not allowed to add credentials to groups.").Send(c)
	}

	var req CredentialGroupAddCredentialRequest
	if err := c.BodyParser(&req); err != nil {
		return api.Response().BadRequest("invalid JSON body").Send(c)
	}

	// Validate required fields
	if req.CredentialID == "" {
		return api.Response().BadRequest("credential_id is required").Send(c)
	}
	if req.CredentialGroupID == "" {
		return api.Response().BadRequest("credential_group_id is required").Send(c)
	}

	now := time.Now().UTC()
	detail := models.LoginCredentialGroupDetails{
		ID:                uuid.NewString(),
		CredentialID:      req.CredentialID,
		CredentialGroupID: req.CredentialGroupID,
		CashLimit:         req.CashLimit,
		Status:            typing.StatusActive,
		UpdatedAt:         now,
	}

	// Insert into database
	query := `
		INSERT INTO users.login_credential_group_details (id, credential_id, credential_group_id, cash_limit, status, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.Postgres.Exec(query,
		detail.ID,
		detail.CredentialID,
		detail.CredentialGroupID,
		detail.CashLimit,
		detail.Status,
		detail.UpdatedAt,
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	return api.Response().
		Data(detail).
		Message("Credential added to group successfully").
		Send(c)
}
