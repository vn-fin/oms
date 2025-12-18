package controller

import (
	"github.com/vn-fin/oms/pkg/mem"
	"github.com/vn-fin/xpb/xpb/order"
)

func GetLatestTick(symbol string) *order.TickInfo {
	mem.Mutex.RLock()
	defer mem.Mutex.RUnlock()
	tick, ok := mem.LatestTickMap[symbol]
	if !ok {
		return nil
	}
	return &tick
}
