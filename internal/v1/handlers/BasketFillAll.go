package handlers

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
)

// BasketFillAll
// @Summary Fill all basket execution
// @Description Fill all basket execution with match data
// @Tags Execution
// @Produce json
// @Param basket_id path string true "Basket ID"
// @Param execution_id path string true "Execution Session ID"
// @Success 200 {object} models.DefaultResponseModel
// @Failure 400 {object} models.DefaultResponseModel
// @Failure 401 {object} models.DefaultResponseModel
// @Failure 500 {object} models.DefaultResponseModel
// @Security BearerAuth
// @Router /oms/v1/baskets/{basket_id}/executions/{execution_id}/match [post]
func BasketFillAll(c *fiber.Ctx) error {
	sessionID := strings.TrimSpace(c.Params("execution_id"))
	if sessionID == "" {
		return api.Response().BadRequest("execution_id is required").Send(c)
	}

	basketID := strings.TrimSpace(c.Params("basket_id"))
	if basketID == "" {
		return api.Response().BadRequest("basket_id is required").Send(c)
	}

	// Get userID from context
	userID := api.GetUserID(c)
	if userID == "" {
		return api.Response().Unauthorized("user not authenticated").Send(c)
	}

	// Query all user orders for the session with join to verify user ownership
	var orders []models.UserOrder
	query := `
		SELECT
			uo.id,
			uo.credential_id,
			uo.session_id,
			uo.symbol,
			uo.symbol_type,
			uo.side,
			uo.order_type,
			uo.order_price,
			uo.matched_price,
			uo.quantity,
			uo.filled_qty,
			uo.remaining_qty,
			uo.status,
			uo.created_at,
			uo.updated_at
		FROM users.credential_groups cg
		INNER JOIN execution.baskets b ON cg.id = b.group_id
		INNER JOIN execution.basket_execute_sessions bes ON b.id = bes.basket_id
		INNER JOIN execution.user_orders uo ON bes.id = uo.session_id
		WHERE cg.user_id = ?
			AND bes.basket_id = ?
			AND bes.id = ?
		ORDER BY uo.created_at DESC
	`
	_, err := db.Postgres.Query(&orders, query, userID, basketID, sessionID)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	if len(orders) == 0 {
		return api.Response().BadRequest("No orders found for this session").Send(c)
	}

	// Update all orders from "created" to "pending"
	updateQuery := `
		UPDATE execution.user_orders
		SET status = ?, updated_at = ?
		WHERE session_id = ? AND status = ?
	`
	_, err = db.Postgres.Exec(updateQuery, typing.OrderStatusPending, time.Now(), sessionID, typing.OrderStatusCreated)
	if err != nil {
		return api.Response().InternalError(err).Send(c)
	}

	// Start mock matching process in background
	go mockMatchingProcess(sessionID)

	return api.Response().
		Message("Orders updated to pending and matching process started").
		Send(c)
}

// mockMatchingProcess simulates order matching: 20% quantity every 5 seconds
func mockMatchingProcess(sessionID string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Get all pending/partial-filled orders
		var orders []models.UserOrder
		query := `
			SELECT id, quantity, filled_qty, remaining_qty, status, order_price
			FROM execution.user_orders
			WHERE session_id = ? AND status IN (?, ?)
		`
		_, err := db.Postgres.Query(&orders, query, sessionID, typing.OrderStatusPending, typing.OrderStatusPartialFilled)
		if err != nil {
			log.Printf("Error querying orders for matching: %v", err)
			return
		}

		// If no orders left to fill, stop the process
		if len(orders) == 0 {
			log.Printf("All orders filled for session %s", sessionID)
			return
		}

		// Process each order
		for _, order := range orders {
			// Calculate 20% of total quantity
			fillAmount := order.Quantity * 0.2
			newFilledQty := order.FilledQty + fillAmount
			newRemainingQty := order.Quantity - newFilledQty

			// Ensure we don't overfill
			if newFilledQty >= order.Quantity {
				newFilledQty = order.Quantity
				newRemainingQty = 0
			}

			// Determine new status
			var newStatus typing.OrderStatus
			if newFilledQty >= order.Quantity {
				newStatus = typing.OrderStatusFilled
			} else if newFilledQty > 0 {
				newStatus = typing.OrderStatusPartialFilled
			} else {
				newStatus = typing.OrderStatusPending
			}

			// Update order
			updateQuery := `
				UPDATE execution.user_orders
				SET filled_qty = ?,
					remaining_qty = ?,
					matched_price = ?,
					status = ?,
					updated_at = ?
				WHERE id = ?
			`
			_, err := db.Postgres.Exec(updateQuery, newFilledQty, newRemainingQty, order.OrderPrice, newStatus, time.Now(), order.ID)
			if err != nil {
				log.Printf("Error updating order %s: %v", order.ID, err)
			} else {
				log.Printf("Order %s: filled %.2f/%.2f, status: %s", order.ID, newFilledQty, order.Quantity, newStatus)
			}
		}
	}
}
