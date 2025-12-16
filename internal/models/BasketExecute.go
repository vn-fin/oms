package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

// BasketExecuteRequest represents the request body for executing a basket
type BasketExecuteRequest struct {
	PriceLevel typing.PriceLevel `json:"price_level"`
	ActionType typing.ActionType `json:"action_type"`
	Weight     float64           `json:"weight"`
	FutureSize float64           `json:"future_size"`
}

// BasketExecuteSession represents a record in execution.basket_execute_sessions table
type BasketExecuteSession struct {
	ID            string             `json:"id"`
	BasketID      string             `json:"basket_id"`
	Weight        float64            `json:"weight"`
	PriceLevel    typing.PriceLevel  `json:"price_level"`
	ActionType    typing.ActionType  `json:"action_type"`
	FutureSize    float64            `json:"future_size"`
	EstimatedCash float64            `json:"estimated_cash"`
	MatchedCash   float64            `json:"matched_cash"`
	OrderStatus   typing.OrderStatus `json:"order_status"`
	CreatedBy     string             `json:"created_by"`
	CreatedAt     time.Time          `json:"created_at"`
}
