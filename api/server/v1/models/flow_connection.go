package models

import (
	"time"
)

// swagger:model
type FlowConnection struct {
	Uuid          string    `json:"uuid"`
	Name          string    `json:"name" valid:"MaxSize(254);Required"`
	ElementFrom   string    `json:"element_from" valid:"Required"`
	ElementTo     string    `json:"element_to" valid:"Required"`
	PointFrom     int64     `json:"point_from" valid:"Required"`
	PointTo       int64     `json:"point_to" valid:"Required"`
	FlowId        int64     `json:"flow_id" valid:"Required"`
	Direction     string    `json:"direction"`
	GraphSettings string    `json:"graph_settings"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
