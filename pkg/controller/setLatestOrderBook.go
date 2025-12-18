package controller

import (
	"github.com/vn-fin/oms/pkg/mem"
	"github.com/vn-fin/xpb/xpb/order"
)

func SetLatestOrderBook(symbol string, orderBook order.OrderBookInfo) {
	mem.Mutex.Lock()
	defer mem.Mutex.Unlock()
	mem.LatestOrderBookMap[symbol] = orderBook
}
