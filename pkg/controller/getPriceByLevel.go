package controller

import (
	"github.com/vn-fin/oms/internal/typing"
)

// GetPriceByLevel returns the price for a symbol at the specified price_level
func GetPriceByLevel(symbol string, priceLevel typing.PriceLevel) float64 {
	priceInfo := GetPriceInfo(symbol)
	if priceInfo == nil {
		return 0
	}

	switch priceLevel {
	case typing.PriceLevelBid01:
		return priceInfo.Bid1
	case typing.PriceLevelBid02:
		return priceInfo.Bid2
	case typing.PriceLevelBid03:
		return priceInfo.Bid3
	case typing.PriceLevelAsk01:
		return priceInfo.Ask1
	case typing.PriceLevelAsk02:
		return priceInfo.Ask2
	case typing.PriceLevelAsk03:
		return priceInfo.Ask3
	case typing.PriceLevelMid:
		return priceInfo.Mid
	case typing.PriceLevelCeil:
		return priceInfo.Ceil
	case typing.PriceLevelFloor:
		return priceInfo.Floor
	default:
		return 0
	}
}
