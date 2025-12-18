package controller

import (
	"github.com/vn-fin/oms/pkg/mem"
	"github.com/vn-fin/xpb/xpb/order"
)

func GetLatestOrderBook(symbol string) *order.OrderBookInfo {
	mem.Mutex.RLock()
	defer mem.Mutex.RUnlock()
	ob, ok := mem.LatestOrderBookMap[symbol]
	if !ok {
		return nil
	}
	return &ob
}

// PriceInfo contains bid, ask and mid prices
type PriceInfo struct {
	Bid1 float64 `json:"bid1"`
	Bid2 float64 `json:"bid2"`
	Bid3 float64 `json:"bid3"`
	Ask1 float64 `json:"ask1"`
	Ask2 float64 `json:"ask2"`
	Ask3 float64 `json:"ask3"`
	Mid  float64 `json:"mid"`
}

// GetPriceInfo returns 3 bid prices, 3 ask prices and mid price
func GetPriceInfo(symbol string) *PriceInfo {
	mem.Mutex.RLock()
	defer mem.Mutex.RUnlock()

	ob, ok := mem.LatestOrderBookMap[symbol]
	if !ok {
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

	return info
}
