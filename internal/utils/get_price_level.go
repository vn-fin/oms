package utils

import (
	"encoding/json"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/vn-fin/oms/internal/db"
	"github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/typing"
)

// GetPriceInfo fetches the latest price information for a symbol from the database
func GetPriceInfo(symbol string) (*models.PriceInfo, error) {
	var snapshot models.SnapshotStockTopPrice

	// Query the latest snapshot for this symbol
	query := `
		SELECT time, symbol, source, bp, bq, ap, aq, total_bid, total_ask
		FROM vn_market.snapshot_stock_top_price
		WHERE symbol = ?
		ORDER BY time DESC
		LIMIT 1
	`

	_, err := db.PostgresXnoData.QueryOne(&snapshot, query, symbol)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil // No data found
		}
		return nil, fmt.Errorf("error fetching price info: %w", err)
	}

	// Parse JSONB data to arrays
	var bp, bq, ap, aq []float64

	if err := json.Unmarshal(snapshot.BP, &bp); err != nil {
		return nil, fmt.Errorf("error parsing bp: %w", err)
	}
	if err := json.Unmarshal(snapshot.BQ, &bq); err != nil {
		return nil, fmt.Errorf("error parsing bq: %w", err)
	}
	if err := json.Unmarshal(snapshot.AP, &ap); err != nil {
		return nil, fmt.Errorf("error parsing ap: %w", err)
	}
	if err := json.Unmarshal(snapshot.AQ, &aq); err != nil {
		return nil, fmt.Errorf("error parsing aq: %w", err)
	}

	// Build PriceInfo
	priceInfo := &models.PriceInfo{
		Symbol: symbol,
	}

	// Extract bid prices (top 3)
	if len(bp) > 0 {
		priceInfo.Bid1 = bp[0]
	}
	if len(bp) > 1 {
		priceInfo.Bid2 = bp[1]
	}
	if len(bp) > 2 {
		priceInfo.Bid3 = bp[2]
	}

	// Extract ask prices (top 3)
	if len(ap) > 0 {
		priceInfo.Ask1 = ap[0]
	}
	if len(ap) > 1 {
		priceInfo.Ask2 = ap[1]
	}
	if len(ap) > 2 {
		priceInfo.Ask3 = ap[2]
	}

	// Calculate mid price
	if priceInfo.Bid1 > 0 && priceInfo.Ask1 > 0 {
		priceInfo.Mid = (priceInfo.Bid1 + priceInfo.Ask1) / 2
	}

	// Fetch Ceil and Floor from snapshot_stock_info
	var stockInfo models.SnapshotStockInfo
	stockInfoQuery := `
		SELECT symbol, ceil, floor
		FROM vn_market.snapshot_stock_info
		WHERE symbol = ?
		ORDER BY time DESC
		LIMIT 1
	`
	_, err = db.PostgresXnoData.QueryOne(&stockInfo, stockInfoQuery, symbol)
	if err != nil {
		if !errors.Is(err, pg.ErrNoRows) {
			return nil, fmt.Errorf("error fetching stock info: %w", err)
		}
		// If no stock info found, set ceil and floor to 0
		priceInfo.Ceil = 0
		priceInfo.Floor = 0
	} else {
		priceInfo.Ceil = stockInfo.Ceil
		priceInfo.Floor = stockInfo.Floor
	}

	return priceInfo, nil
}

// GetPriceByLevel returns the price for a symbol at the specified price_level
func GetPriceByLevel(symbol string, priceLevel typing.PriceLevel) (float64, error) {
	priceInfo, err := GetPriceInfo(symbol)
	if err != nil {
		return 0, err
	}
	if priceInfo == nil {
		return 0, nil
	}

	switch priceLevel {
	case typing.PriceLevelBid01:
		return priceInfo.Bid1, nil
	case typing.PriceLevelBid02:
		return priceInfo.Bid2, nil
	case typing.PriceLevelBid03:
		return priceInfo.Bid3, nil
	case typing.PriceLevelAsk01:
		return priceInfo.Ask1, nil
	case typing.PriceLevelAsk02:
		return priceInfo.Ask2, nil
	case typing.PriceLevelAsk03:
		return priceInfo.Ask3, nil
	case typing.PriceLevelMid:
		return priceInfo.Mid, nil
	case typing.PriceLevelCeil:
		return priceInfo.Ceil, nil
	case typing.PriceLevelFloor:
		return priceInfo.Floor, nil
	default:
		return 0, nil
	}
}
