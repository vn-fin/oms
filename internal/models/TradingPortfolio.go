package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

// TradingPortfolio
type TradingPortfolio struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Info        []EnsembleWeight    `json:"info"`
	CreatedBy   string              `json:"created_by"`
	UpdatedBy   string              `json:"updated_by"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Status      typing.RecordStatus `json:"status"`
}
