package models

import (
	"time"

	"github.com/vn-fin/oms/internal/typing"
)

type Basket struct {
	ID          string              `json:"id" pg:"id,pk"`
	GroupID     string              `json:"groupId" pg:"group_id,pk"`
	Name        string              `json:"name" pg:"name"`
	Description string              `json:"description" pg:"description"`
	Info        []BasketInfo        `json:"info" pg:"info,type:jsonb"`
	HedgeConfig []BasketHedgeConfig `json:"hedge_config" pg:"hedge_config,type:jsonb"`
	CreatedBy   string              `json:"created_by" pg:"created_by"`
	UpdatedBy   string              `json:"updated_by" pg:"updated_by"`
	CreatedAt   time.Time           `json:"created_at" pg:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" pg:"updated_at"`
	Status      typing.RecordStatus `json:"status" pg:"status"`
}
