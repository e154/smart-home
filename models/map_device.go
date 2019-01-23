package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type MapDevice struct {
	Id         int64              `json:"id"`
	SystemName string             `json:"system_name" valid:"Required"`
	DeviceId   int64              `json:"device_id" valid:"Required"`
	States     []*MapDeviceState  `json:"states"`
	Actions    []*MapDeviceAction `json:"actions"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

func (m *MapDevice) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}
