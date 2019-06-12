package models

import (
	"encoding/json"
	"time"
	"github.com/e154/smart-home/system/uuid"
	"github.com/e154/smart-home/system/validation"
)

type Connection struct {
	Uuid          uuid.UUID       `json:"uuid"`
	Name          string          `json:"name" valid:"MaxSize(254)"`
	ElementFrom   uuid.UUID       `json:"element_from" valid:"Required"`
	ElementTo     uuid.UUID       `json:"element_to" valid:"Required"`
	PointFrom     int64           `json:"point_from" valid:"Required"`
	PointTo       int64           `json:"point_to" valid:"Required"`
	Flow          *Flow           `json:"flow"`
	FlowId        int64           `json:"flow_id" valid:"Required"`
	Direction     string          `json:"direction"`
	GraphSettings json.RawMessage `json:"graph_settings"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (d *Connection) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
