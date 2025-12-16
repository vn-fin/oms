package models

import "github.com/vn-fin/oms/internal/typing"

type BasketHedgeConfig struct {
	Symbol     string            `json:"symbol"`
	SymbolType typing.SymbolType `json:"symbol_type"`
	Size       float64           `json:"size"`
	Direction  float64           `json:"direction"`
}
