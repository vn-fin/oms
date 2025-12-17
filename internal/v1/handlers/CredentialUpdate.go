package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
)

type CredentialUpdateRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Info        []models.CredentialInfo `json:"info"`
}

// CredentialUpdate
// @Summary Update login credential
// @Description Update an existing login credential (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Param credential_id path string true "Credential ID"
// @Param request body CredentialUpdateRequest true "Credential update payload"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 403 {object} models.DefaultResponseModel
// @Failure 404 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credentials/{credential_id} [put]
func CredentialUpdate(c *fiber.Ctx) error {
	// Check admin permission
	userEmail := api.GetUserEmail(c)
	if userEmail != config.AdminEmails {
		return api.Response().Forbidden("You are not allowed to update credentials.").Send(c)
	}

	credentialID := strings.TrimSpace(c.Params("credential_id"))
	if credentialID == "" {
		return api.Response().BadRequest("credential_id is required").Send(c)
	}

	var req CredentialUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return api.Response().BadRequest("invalid JSON body").Send(c)
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	if req.Name == "" {
		return api.Response().BadRequest("name is required").Send(c)
	}

	// Info is required
	if req.Info == nil || len(req.Info) == 0 {
		return api.Response().BadRequest("info is required").Send(c)
	}

	now := time.Now().UTC()

	// Update credential
	query := `
		UPDATE execution.login_credentials
		SET name = ?, description = ?, info = ?, updated_at = ?
		WHERE credential_id = ?
	`
	result, err := db.Postgres.Exec(query,
		req.Name,
		req.Description,
		req.Info,
		now,
		credentialID,
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Check if any row was affected
	if result.RowsAffected() == 0 {
		return api.Response().NotFound("credential not found").Send(c)
	}

	return api.Response().
		Data(fiber.Map{
			"credential_id": credentialID,
			"name":          req.Name,
			"description":   req.Description,
			"info":          req.Info,
			"updated_at":    now,
		}).
		Message("Credential updated successfully").
		Send(c)
}
