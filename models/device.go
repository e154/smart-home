package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type Device struct {
	Id          int64           `json:"id"`
	Device      *Device         `json:"device"`
	NodeId      *int64          `json:"node_id"`
	Address     *int            `json:"address"`
	Baud        int             `json:"baud"`
	Sleep       int64           `json:"sleep"`
	Description string          `json:"description" valid:"MaxSize(254)"`
	Name        string          `json:"name" valid:"MaxSize(254);Required"`
	Status      string          `json:"status" valid:"MaxSize(254)"`
	StopBite    int64           `json:"stop_bite"`
	Timeout     time.Duration   `json:"timeout"`
	Tty         string          `json:"tty" valid:"MaxSize(254)"`
	States      []*DeviceState  `json:"states"`
	Actions     []*DeviceAction `json:"actions"`
	IsGroup     bool            `json:"is_group"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (d *Device) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
