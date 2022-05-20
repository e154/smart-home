package models

import (
	"encoding/json"
	"time"

	"github.com/e154/smart-home/common"
)

// DashboardCardItem ...
type DashboardCardItem struct {
	Id              int64            `json:"id"`
	Title           string           `json:"title" validate:"required"`
	Type            string           `json:"type" validate:"required"`
	Weight          int              `json:"weight"`
	Enabled         bool             `json:"enabled"`
	DashboardCardId int64            `json:"dashboard_card_id" validate:"required"`
	DashboardCard   *DashboardCard   `json:"dashboard_card"`
	EntityId        *common.EntityId `json:"entity_id"`
	Entity          *Entity          `json:"entity"`
	Payload         json.RawMessage  `json:"payload"`
	Hidden          bool             `json:"hidden"`
	Frozen          bool             `json:"frozen"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}
