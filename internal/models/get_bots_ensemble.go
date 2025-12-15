package models

import "time"

// BotWeight represents the bot_id and weight information in the info jsonb field
type BotWeight struct {
	BotID  string  `json:"bot_id"`
	Weight float64 `json:"weight"`
}

// BotsEnsemble represents the portfolio.bots_ensemble table
type BotsEnsemble struct {
	ID          string      `json:"id" pg:"id,pk"`
	Name        string      `json:"name" pg:"name"`
	Description string      `json:"description" pg:"description"`
	Symbol      string      `json:"symbol" pg:"symbol"`
	Type        string      `json:"type" pg:"type"` // daily, intraday, hft
	Info        []BotWeight `json:"info" pg:"info,type:jsonb"`
	CreatedBy   string      `json:"created_by" pg:"created_by"`
	UpdatedBy   string      `json:"updated_by" pg:"updated_by"`
	CreatedAt   time.Time   `json:"created_at" pg:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" pg:"updated_at"`
	Status      string      `json:"status" pg:"status"` // active, disabled
}
