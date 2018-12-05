package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
	"encoding/json"
	. "github.com/e154/smart-home/common"
	. "github.com/e154/smart-home/models/devices"
	"fmt"
)

type Device struct {
	Id          int64           `json:"id"`
	Name        string          `json:"name" valid:"MaxSize(254);Required"`
	Description string          `json:"description" valid:"MaxSize(254)"`
	Status      string          `json:"status" valid:"MaxSize(254)"`
	Device      *Device         `json:"device"`
	Node        *Node           `json:"node"`
	Type        DeviceType      `json:"type"`
	Properties  json.RawMessage `json:"properties" valid:"Required"`
	States      []*DeviceState  `json:"states"`
	Actions     []*DeviceAction `json:"actions"`
	Devices     []*Device       `json:"devices"`
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

func (d *Device) SetProperties(properties interface{}) (ok bool, errs []*validation.Error) {

	var dType DeviceType

	switch v := properties.(type) {
	case *DevSmartBusConfig:
		dType = DevTypeSmartBus
		ok, errs = v.Valid()
	case *DevCommandConfig:
		dType = DevTypeCommand
		ok, errs = v.Valid()
	default:
		fmt.Printf("unknown device config %v", v)
		return
	}

	if !ok || len(errs) > 0 {
		return
	}

	d.Type = dType
	if data, err := json.Marshal(properties); err == nil {
		d.Properties.UnmarshalJSON(data)
	}

	return
}
