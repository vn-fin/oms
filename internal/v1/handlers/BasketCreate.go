package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
)

type BasketCreateRequest struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Info        []models.BasketInfo        `json:"info"`
	HedgeConfig []models.BasketHedgeConfig `json:"hedge_config"`
}

// BasketCreate
// @Summary Create basket
// @Description Create a new basket.
// @Accept json
// @Produce json
// @Param request body BasketCreateRequest true "Basket create payload"
// @Success 201 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets [post]
func BasketCreate(c *fiber.Ctx) error {
	var req BasketCreateRequest
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

	// Get userID from context (set 1by AuthMiddleware)
	userID := api.GetUserID(c)

	now := time.Now().UTC()
	basket := models.Basket{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		Info:        req.Info,
		HedgeConfig: req.HedgeConfig,
		CreatedBy:   userID,
		UpdatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
		Status:      typing.RecordStatusEnabled,
	}

	// Insert into database using raw SQL query
	query := `
		INSERT INTO execution.baskets (id, name, description, info,hedge_config, created_by, updated_by, created_at, updated_at, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,?)
	`
	_, err := db.Postgres.Exec(query,
		basket.ID,
		basket.Name,
		basket.Description,
		basket.Info,
		basket.HedgeConfig,
		basket.CreatedBy,
		basket.UpdatedBy,
		basket.CreatedAt,
		basket.UpdatedAt,
		basket.Status,
	)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	return api.Response().
		Status(fiber.StatusCreated).
		Data(basket).
		Message("Basket created successfully").
		Send(c)
}
