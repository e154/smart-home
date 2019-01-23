package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type MapDeviceAction struct {
	Id             int64     `json:"id"`
	DeviceActionId int64     `json:"device_action_id" valid:"Required"`
	MapDeviceId    int64     `json:"map_device_id" valid:"Required"`
	ImageId        int64     `json:"image_id" valid:"Required"`
	Type           string    `json:"type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (m *MapDeviceAction) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}
