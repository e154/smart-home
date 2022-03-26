package models

import "time"

// DashboardTab ...
type DashboardTab struct {
	Id          int64            `json:"id"`
	Name        string           `json:"name" validate:"required"`
	ColumnWidth int              `json:"column_width"`
	Gap         int              `json:"gap"`
	Background  string           `json:"background"`
	Icon        string           `json:"icon"`
	Enabled     bool             `json:"enabled"`
	Weight      int              `json:"weight"`
	DashboardId int64            `json:"dashboard_id" validate:"required"`
	Dashboard   *Dashboard       `json:"dashboard"`
	Cards       []*DashboardCard `json:"cards"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
