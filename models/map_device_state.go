package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type MapDeviceState struct {
	Id            int64     `json:"id"`
	DeviceStateId int64     `json:"device_state_id" valid:"Required"`
	MapDeviceId   int64     `json:"map_device_id" valid:"Required"`
	ImageId       int64     `json:"image_id" valid:"Required"`
	Style         string    `json:"style"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (m *MapDeviceState) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}
