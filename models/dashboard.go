package models

import (
	"time"

	"github.com/e154/smart-home/common"
)

// Dashboard ...
type Dashboard struct {
	Id          int64                       `json:"id"`
	Name        string                      `json:"name" validate:"required"`
	Description string                      `json:"description"`
	Enabled     bool                        `json:"enabled"`
	AreaId      *int64                      `json:"area_id"`
	Area        *Area                       `json:"area"`
	Tabs        []*DashboardTab             `json:"tabs"`
	Entities    map[common.EntityId]*Entity `json:"entities"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}