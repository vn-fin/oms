package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

type CredentialGroup struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	UserID    string               `json:"user_id"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	Status    typing.AccountStatus `json:"status"`
}
