package models

import (
	"encoding/json"
	"time"
)

// SnapshotStockTopPrice represents the snapshot_stock_top_price table in xno_data.vn_market
type SnapshotStockTopPrice struct {
	Time     time.Time       `pg:"time" json:"time"`
	Symbol   string          `pg:"symbol" json:"symbol"`
	Source   string          `pg:"source" json:"source"`
	BP       json.RawMessage `pg:"bp" json:"bp"`
	BQ       json.RawMessage `pg:"bq" json:"bq"`
	AP       json.RawMessage `pg:"ap" json:"ap"`
	AQ       json.RawMessage `pg:"aq" json:"aq"`
	TotalBid float64         `pg:"total_bid" json:"total_bid"`
	TotalAsk float64         `pg:"total_ask" json:"total_ask"`
}

// SnapshotStockInfo represents the snapshot_stock_info table in xno_data.vn_market
type SnapshotStockInfo struct {
	Symbol string  `pg:"symbol" json:"symbol"`
	Ceil   float64 `pg:"ceil" json:"ceil"`
	Floor  float64 `pg:"floor" json:"floor"`
}

// PriceInfo contains processed price information for a symbol
type PriceInfo struct {
	Symbol string  `json:"symbol"`
	Bid1   float64 `json:"bid1"`
	Bid2   float64 `json:"bid2"`
	Bid3   float64 `json:"bid3"`
	Ask1   float64 `json:"ask1"`
	Ask2   float64 `json:"ask2"`
	Ask3   float64 `json:"ask3"`
	Mid    float64 `json:"mid"`
	Ceil   float64 `json:"ceil"`
	Floor  float64 `json:"floor"`
}
