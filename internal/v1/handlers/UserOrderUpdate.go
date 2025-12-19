package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
	"github.com/vn-fin/oms/pkg/controller"
)

// UserOrderUpdatePrice updates the order_price based on price_level from orderbook
// @Summary Update order price by price level
// @Description Update order_price using bid01/bid02/bid03/ask01/ask02/ask03/mid/ceil/floor from latest orderbook
// @Tags Orders
// @Produce json
// @Param order_id path string true "Order ID"
// @Param price_level query string true "Price level: bid01, bid02, bid03, ask01, ask02, ask03, mid, ceil, floor"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 404 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/orders/{order_id}/update-price [put]
func UserOrderUpdatePrice(c *fiber.Ctx) error {
	orderID := strings.TrimSpace(c.Params("order_id"))
	if orderID == "" {
		return api.Response().BadRequest("order_id is required").Send(c)
	}

	priceLevel := typing.PriceLevel(strings.ToLower(strings.TrimSpace(c.Query("price_level"))))
	if priceLevel == "" {
		return api.Response().BadRequest("price_level is required (bid01, bid02, bid03, ask01, ask02, ask03, mid, ceil, floor)").Send(c)
	}

	// Validate price_level
	if !priceLevel.Valid() {
		return api.Response().BadRequest("invalid price_level. Valid: bid01, bid02, bid03, ask01, ask02, ask03, mid, ceil, floor").Send(c)
	}

	// Get userID from context
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	// Get the order first to find symbol
	var order models.UserOrder
	query := `
		SELECT id, credential_id, session_id, symbol, symbol_type, side, order_price, matched_price, quantity, filled_qty, remaining_qty, status, created_at, updated_at
		FROM execution.user_orders
		WHERE id = ?
	`
	_, err := db.Postgres.QueryOne(&order, query, orderID)
	if err != nil {
		return api.Response().NotFound("order not found").Send(c)
	}

	// Get price from latest message based on price_level
	newPrice := controller.GetPriceByLevel(order.Symbol, priceLevel)
	if newPrice <= 0 {
		return api.Response().BadRequest("price not available for " + string(priceLevel)).Send(c)
	}

	// Update order_price in database
	updateQuery := `
		UPDATE execution.user_orders
		SET order_price = ?, updated_at = ?
		WHERE id = ?
	`
	_, err = db.Postgres.Exec(updateQuery, newPrice, time.Now(), orderID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	return api.Response().
		Data(fiber.Map{
			"order_id":    orderID,
			"symbol":      order.Symbol,
			"price_level": priceLevel,
			"old_price":   order.OrderPrice,
			"new_price":   newPrice,
		}).
		Message("Order price updated successfully").
		Send(c)
}
