package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type DeviceAction struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name" valid:"MaxSize(254);Required"`
	Description string    `json:"description"`
	Device      *Device   `json:"device"`
	DeviceId    int64     `json:"device_id" valid:"Required"`
	Script      *Script   `json:"script"`
	ScriptId    int64     `json:"script_id" valid:"Required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *DeviceAction) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
