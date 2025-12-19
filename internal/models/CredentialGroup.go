package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

type CredentialGroup struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	UserID         string               `json:"user_id"`
	TotalCashLimit float64              `json:"total_cash_limit" pg:"-"`
	TotalBalance   float64              `json:"total_balance" pg:"-"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	Status         typing.AccountStatus `json:"status"`
}
