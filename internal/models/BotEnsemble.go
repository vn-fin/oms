package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

// BotEnsemble represents the portfolio.bots_ensemble table
type BotEnsemble struct {
	ID          string           `json:"id" pg:"id,pk"`
	Name        string           `json:"name" pg:"name"`
	Description string           `json:"description" pg:"description"`
	Symbol      string           `json:"symbol" pg:"symbol"`
	BotType     typing.BotType   `json:"bot_type" pg:"bot_type"` // daily, intraday, hft
	Info        []BotWeight      `json:"info" pg:"info,type:jsonb"`
	CreatedBy   string           `json:"created_by" pg:"created_by"`
	UpdatedBy   string           `json:"updated_by" pg:"updated_by"`
	CreatedAt   time.Time        `json:"created_at" pg:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" pg:"updated_at"`
	BotStatus   typing.BotStatus `json:"bot_status" pg:"bot_status"` // active, disabled
}
