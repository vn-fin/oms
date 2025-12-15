package models

import "time"

type BotSignal struct {
	BotID      string    `json:"botID" pg:"bot_id"`
	Name       string    `json:"name" pg:"recall_name"`
	Symbol     string    `json:"symbol" pg:"symbol"`
	SymbolType string    `json:"symbol_type" pg:"symbol_type"`
	Candle     time.Time `json:"candle" pg:"candle"`
	Price      float64   `json:"price" pg:"price"`
	Position   int64     `json:"position" pg:"position"`
	Timeframe  string    `json:"timeframe" pg:"timeframe"`
	Display    bool      `json:"display" pg:"display"`
	Status     string    `json:"status" pg:"status"`
}
