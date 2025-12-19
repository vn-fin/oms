package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

type CredentialGroupDetails struct {
	ID                string               `json:"id"`
	CredentialID      string               `json:"credential_id"`
	CredentialGroupID string               `json:"credential_group_id"`
	CashLimit         float64              `json:"cash_limit"`
	Balance           float64              `json:"balance"`
	Status            typing.AccountStatus `json:"status"`
	UpdatedAt         time.Time            `json:"updated_at"`
}
