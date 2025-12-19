package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

type Credential struct {
	CredentialID string               `json:"credential_id" pg:"credential_id"`
	Name         string               `json:"name" pg:"name"`
	Description  string               `json:"description" pg:"description"`
	Info         []CredentialInfo     `json:"info" pg:"info"`
	CreatedAt    time.Time            `json:"created_at" pg:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at" pg:"updated_at"`
	Status       typing.AccountStatus `json:"status" pg:"status"`
}

type CredentialByGroup struct {
	CredentialID string               `json:"credential_id"`
	Name         string               `json:"name"`
	Status       typing.AccountStatus `json:"status"`
}
