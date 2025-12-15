package models

import "time"

type BotsSignal struct {
	BotID     string    `json:"botID" pg:"bot_id"`
	Name      string    `json:"name" pg:"recall_name"`
	Display   bool      `json:"display" pg:"display"`
	Symbol    string    `pg:"symbol"`
	Candle    time.Time `pg:"candle"`
	Price     float64   `pg:"price"`
	Position  int64     `pg:"position"`
	Timeframe string    `pg:"timeframe"`
}
