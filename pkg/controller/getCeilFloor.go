package controller

import (
	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/db"
)

// HistoryStockInfo represents the latest stock info from xno_data.vn_market.history_stock_info
type HistoryStockInfo struct {
	Symbol string  `pg:"symbol"`
	Ceil   float64 `pg:"ceil"`
	Floor  float64 `pg:"floor"`
}

// GetCeilFloor returns the ceil and floor prices for a symbol from the database
func GetCeilFloor(symbol string) (ceil float64, floor float64) {
	if db.PostgresXnoData == nil {
		log.Warn().Msg("PostgresXnoData is not initialized")
		return 0, 0
	}

	var stockInfo HistoryStockInfo
	query := `
		SELECT symbol, ceil, floor
		FROM vn_market.snapshot_stock_info
		WHERE symbol = ?
		LIMIT 1
	`

	_, err := db.PostgresXnoData.QueryOne(&stockInfo, query, symbol)
	if err != nil {
		log.Debug().Err(err).Str("symbol", symbol).Msg("Failed to get ceil/floor from database")
		return 0, 0
	}

	return stockInfo.Ceil, stockInfo.Floor
}
