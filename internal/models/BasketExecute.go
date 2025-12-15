package models

import "github.com/vn-fin/oms/internal/typing"

type BasketExecute struct {
	ActionType typing.ActionType `json:"action_type"`
	PriceLevel typing.PriceLevel `json:"price_level"`
	AutoFuture bool              `json:"auto_future"`
}
