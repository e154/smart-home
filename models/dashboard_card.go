package models

import (
	"encoding/json"
	"time"
)

// DashboardCard ...
type DashboardCard struct {
	Id             int64           `json:"id"`
	Title          string          `json:"title" validate:"required"`
	Height         int             `json:"height" validate:"required"`
	Width          int             `json:"width" validate:"required"`
	Background     string          `json:"background" validate:"required"`
	Weight         int             `json:"weight"`
	Enabled        bool            `json:"enabled"`
	DashboardTabId int64           `json:"dashboard_tab_id" validate:"required"`
	DashboardTab   *DashboardTab   `json:"dashboard_tab"`
	Payload        json.RawMessage `json:"payload"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
