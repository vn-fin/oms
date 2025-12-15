package models

// BotWeight represents the bot_id and weight information in the info jsonb field
type BotWeight struct {
	BotID  string  `json:"bot_id"`
	Weight float64 `json:"weight"`
}
