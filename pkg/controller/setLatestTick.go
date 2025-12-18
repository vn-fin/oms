package controller

import (
	"github.com/vn-fin/oms/pkg/mem"
	"github.com/vn-fin/xpb/xpb/order"
)

func SetLatestTick(symbol string, tick order.TickInfo) {
	mem.Mutex.Lock()
	defer mem.Mutex.Unlock()
	mem.LatestTickMap[symbol] = tick
}
