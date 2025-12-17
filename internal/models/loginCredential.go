package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

type Credential struct {
	CredentialID string               `json:"credential_id"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Info         []CredentialInfo     `json:"info"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	Status       typing.AccountStatus `json:"status"`
}
