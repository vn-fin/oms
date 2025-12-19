package mem

import "github.com/vn-fin/xpb/xpb/order"

type OrderBookInfo order.OrderBookInfo

var latestOrderBookMap = make(map[string]OrderBookInfo)

func GetLatestOrderBookMap() map[string]OrderBookInfo {
	Mutex.RLock()
	defer Mutex.RUnlock()
	// Create a copy
	cp := make(map[string]OrderBookInfo, len(latestOrderBookMap))
	for k, v := range latestOrderBookMap {
		cp[k] = v
	}
	return cp
}

func SetLatestOrderBook(symbol string, orderBook OrderBookInfo) bool {
	Mutex.Lock()
	defer Mutex.Unlock()

	// Check if the symbol already exists and if the existing tick is newer
	if prevOB, exists := latestOrderBookMap[symbol]; exists {
		if prevOB.TimeF >= orderBook.TimeF {
			return false
		}
	}

	latestOrderBookMap[symbol] = orderBook
	return true
}

func GetLatestOrderBook(symbol string) *OrderBookInfo {
	Mutex.RLock()
	defer Mutex.RUnlock()
	data, ok := latestOrderBookMap[symbol]
	if !ok {
		return nil
	}

	return &data
}
