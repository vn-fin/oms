package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
	"github.com/vn-fin/oms/internal/utils"
)

// UserOrderBatchUpdatePrice updates order_price for all orders in an execution session
// @Summary Update all order prices in execution session by price level
// @Description Update order_price for all orders in session using bid01/bid02/bid03/ask01/ask02/ask03/mid from latest orderbook
// @Tags Execution
// @Produce json
// @Param basket_id path string true "Basket ID"
// @Param execution_id path string true "Execution Session ID"
// @Param price_level query string true "Price level: bid01, bid02, bid03, ask01, ask02, ask03, mid, ceil, floor"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/{basket_id}/execute/{execution_id}/update-price [put]
func UserOrderBatchUpdatePrice(c *fiber.Ctx) error {
	sessionID := strings.TrimSpace(c.Params("execution_id"))
	if sessionID == "" {
		return api.Response().BadRequest("execution_id is required").Send(c)
	}

	priceLevel := typing.PriceLevel(strings.ToLower(strings.TrimSpace(c.Query("price_level"))))
	// Validate price_level
	if !priceLevel.Valid() {
		return api.Response().BadRequest("invalid price_level. Valid: bid01, bid02, bid03, ask01, ask02, ask03, mid, ceil, floor").Send(c)
	}

	// Get userID from context
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	// Query all user orders for the session
	var orders []models.UserOrder
	query := `
		SELECT id, credential_id, session_id, symbol, symbol_type, side, order_price, matched_price, quantity, filled_qty, remaining_qty, status, created_at, updated_at
		FROM execution.user_orders
		WHERE session_id = ?
		ORDER BY created_at DESC
	`
	_, err := db.Postgres.Query(&orders, query, sessionID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	if len(orders) == 0 {
		return api.Response().NotFound("no orders found in this session").Send(c)
	}

	// Update each order with price from its symbol's orderbook
	type UpdateResult struct {
		OrderID  string  `json:"order_id"`
		Symbol   string  `json:"symbol"`
		OldPrice float64 `json:"old_price"`
		NewPrice float64 `json:"new_price"`
		Status   string  `json:"status"`
		Error    string  `json:"error,omitempty"`
	}

	results := make([]UpdateResult, 0, len(orders))
	successCount := 0
	failCount := 0

	for _, order := range orders {
		result := UpdateResult{
			OrderID:  order.ID,
			Symbol:   order.Symbol,
			OldPrice: order.OrderPrice,
		}

		// Get price from database based on price_level
		newPrice, err := utils.GetPriceByLevel(order.Symbol, priceLevel)
		if err != nil {
			return api.Response().InternalError(err).Send(c)
		}
		if newPrice <= 0 {
			return api.Response().BadRequest("price not available for symbol " + order.Symbol + " at price level " + string(priceLevel)).Send(c)
		}

		// Update order_price in database
		updateQuery := `
			UPDATE execution.user_orders
			SET order_price = ?, updated_at = ?
			WHERE id = ?
		`
		_, err = db.Postgres.Exec(updateQuery, newPrice, time.Now(), order.ID)
		if err != nil {
			return api.Response().InternalError(err).Send(c)
		}

		result.NewPrice = newPrice
		result.Status = "success"
		successCount++
		results = append(results, result)
	}

	return api.Response().
		Data(fiber.Map{
			"session_id":    sessionID,
			"price_level":   priceLevel,
			"total":         len(orders),
			"success_count": successCount,
			"fail_count":    failCount,
			"results":       results,
		}).
		Message("Batch update completed").
		Send(c)
}
