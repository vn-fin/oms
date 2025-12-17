package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

type UserOrder struct {
	ID           string             `json:"id"`
	CredentialID string             `json:"credential_id"`
	SessionID    string             `json:"session_id"`
	Symbol       string             `json:"symbol"`
	SymbolType   typing.SymbolType  `json:"symbol_type"`
	Side         typing.ActionType  `json:"side"`
	OrderPrice   float64            `json:"order_price"`
	MatchedPrice float64            `json:"matched_price"`
	Quantity     float64            `json:"quantity"`
	FilledQty    float64            `json:"filled_qty"`
	RemainingQty float64            `json:"remaining_qty"`
	Status       typing.OrderStatus `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}
