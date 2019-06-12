package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type DeviceState struct {
	Id          int64     `json:"id"`
	Description string    `json:"description"`
	SystemName  string    `json:"system_name" valid:"MaxSize(254);Required"`
	Device      *Device   `json:"device"`
	DeviceId    int64     `json:"device_id" valid:"Required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *DeviceState) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
