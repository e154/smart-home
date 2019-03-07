package models

import (
	"encoding/json"
	"time"
	"github.com/e154/smart-home/system/uuid"
)

type FlowConnectionModel struct {
	Uuid          uuid.UUID       `json:"uuid"`
	Name          string          `json:"name" valid:"MaxSize(254);Required"`
	ElementFrom   uuid.UUID       `json:"element_from" valid:"Required"`
	ElementTo     uuid.UUID       `json:"element_to" valid:"Required"`
	PointFrom     int64           `json:"point_from" valid:"Required"`
	PointTo       int64           `json:"point_to" valid:"Required"`
	FlowId        int64           `json:"flow_id" valid:"Required"`
	Direction     string          `json:"direction"`
	GraphSettings json.RawMessage `json:"graph_settings"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

