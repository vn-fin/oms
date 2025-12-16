package models

type BasketHedgeConfig struct {
	Symbol     string  `json:"symbol"`
	SymbolType float64 `json:"symbol_type"`
	Size       float64 `json:"size"`
	Direction  float64 `json:"direction"`
}
