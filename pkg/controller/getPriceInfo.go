package controller

import (
	"github.com/vn-fin/oms/pkg/mem"
)

// PriceInfo contains bid, ask, mid, ceil, floor prices
type PriceInfo struct {
	Bid1  float64 `json:"bid1"`
	Bid2  float64 `json:"bid2"`
	Bid3  float64 `json:"bid3"`
	Ask1  float64 `json:"ask1"`
	Ask2  float64 `json:"ask2"`
	Ask3  float64 `json:"ask3"`
	Mid   float64 `json:"mid"`
	Ceil  float64 `json:"ceil"`
	Floor float64 `json:"floor"`
}

// GetPriceInfo returns 3 bid prices, 3 ask prices, mid, ceil, floor for a symbol
func GetPriceInfo(symbol string) *PriceInfo {
	ob := mem.GetLatestOrderBook(symbol)
	if ob == nil {
		return nil
	}

	info := &PriceInfo{}

	// Get top 3 bid prices
	if len(ob.BidPrices) > 0 {
		info.Bid1 = ob.BidPrices[0]
	}
	if len(ob.BidPrices) > 1 {
		info.Bid2 = ob.BidPrices[1]
	}
	if len(ob.BidPrices) > 2 {
		info.Bid3 = ob.BidPrices[2]
	}

	// Get top 3 ask prices
	if len(ob.AskPrices) > 0 {
		info.Ask1 = ob.AskPrices[0]
	}
	if len(ob.AskPrices) > 1 {
		info.Ask2 = ob.AskPrices[1]
	}
	if len(ob.AskPrices) > 2 {
		info.Ask3 = ob.AskPrices[2]
	}

	// Calculate mid price
	if info.Bid1 > 0 && info.Ask1 > 0 {
		info.Mid = (info.Bid1 + info.Ask1) / 2
	}

	// Get ceil and floor from StockInfo
	stockInfo := mem.GetLastStockInfo(symbol)
	if stockInfo != nil {
		info.Ceil = stockInfo.Ceil
		info.Floor = stockInfo.Floor
	}

	return info
}
