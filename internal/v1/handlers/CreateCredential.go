package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
)

type CredentialCreateRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Info        []models.CredentialInfo `json:"info"`
}

// CredentialCreate
// @Summary Create login credential
// @Description Create a new login credential (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body CredentialCreateRequest true "Credential create payload"
// @Success 201 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 403 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/credentials [post]
func CredentialCreate(c *fiber.Ctx) error {
	// Check admin permission

	userEmail := api.GetUserEmail(c)
	if userEmail != config.AdminEmails {
		return api.Response().Forbidden("You are not allowed to create credentials.").Send(c)
	}

	var req CredentialCreateRequest
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

	// Get userID from context

	now := time.Now().UTC()
	credential := models.Credential{
		CredentialID: uuid.NewString(),
		Name:         req.Name,
		Description:  req.Description,
		//Info:         req.Info,
		CreatedAt: now,
		UpdatedAt: now,
		Status:    typing.StatusActive,
	}

	// Insert into database
	query := `
		INSERT INTO users.login_credentials (id, name, description, info, created_at, updated_at, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Postgres.Exec(query,
		credential.CredentialID,
		credential.Name,
		credential.Description,
		//credential.Info,
		credential.CreatedAt,
		credential.UpdatedAt,
		credential.Status,
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	return api.Response().
		Data(credential).
		Message("Credential created successfully").
		Send(c)
}
