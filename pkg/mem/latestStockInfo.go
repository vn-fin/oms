package mem

import (
	"time"
)

type StockInfo struct {
	Symbol string    `json:"symbol"`
	TimeF  float64   `json:"time"`
	TimeT  time.Time `json:"_"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Avg    float64   `json:"avg"`
	Ceil   float64   `json:"ceil"`
	Floor  float64   `json:"floor"`
	Prior  float64   `json:"prior"`
}

func (c *StockInfo) Build() {
	seconds := int64(c.TimeF)
	nanos := int64((c.TimeF - float64(seconds)) * 1e9)
	c.TimeT = time.Unix(seconds, nanos)
}

var latestStockInfoMap = make(map[string]StockInfo)

func GetLatestStockInfoMap() map[string]StockInfo {
	Mutex.RLock()
	defer Mutex.RUnlock()
	// Create a copy
	cp := make(map[string]StockInfo, len(latestStockInfoMap))
	for k, v := range latestStockInfoMap {
		cp[k] = v
	}
	return cp
}

func SetLatestStockInfo(symbol string, tick StockInfo) bool {
	Mutex.Lock()
	defer Mutex.Unlock()

	// Check if the symbol already exists and if the existing tick is newer
	if prev, exists := latestStockInfoMap[symbol]; exists {
		if prev.TimeF >= tick.TimeF {
			return false
		}
	}
	latestStockInfoMap[symbol] = tick
	return true
}

func GetLastStockInfo(symbol string) *StockInfo {
	Mutex.RLock()
	defer Mutex.RUnlock()
	data, ok := latestStockInfoMap[symbol]
	if !ok {
		return nil
	}
	return &data
}
