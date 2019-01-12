package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
	"encoding/json"
	. "github.com/e154/smart-home/common"
	. "github.com/e154/smart-home/models/devices"
	"github.com/mitchellh/mapstructure"
	"fmt"
)

type Device struct {
	Id          int64           `json:"id"`
	Name        string          `json:"name" valid:"MaxSize(254);Required"`
	Description string          `json:"description" valid:"MaxSize(254)"`
	Status      string          `json:"status" valid:"MaxSize(254)"`
	Device      *Device         `json:"device"`
	DeviceId    *int64          `json:"device_id"`
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
		return
	}

	var out interface{}
	switch d.Type {
	case DevTypeModbus:
		out = &DevModBusConfig{}
	case DevTypeSmartBus:
		out = &DevSmartBusConfig{}
	case DevTypeCommand:
		out = &DevCommandConfig{}
	case DevTypeDefault:

	default:
		log.Warningf("unknown device config %v", d.Type)
		return
	}

	var data []byte
	var err error
	data, err = d.Properties.MarshalJSON()
	if err = json.Unmarshal(data, &out); err != nil {
		return
	}

	switch v := out.(type) {
	case *DevSmartBusConfig:
		ok, errs = v.Valid()
	case *DevCommandConfig:
		ok, errs = v.Valid()
	}

	return
}

func (d *Device) SetProperties(properties interface{}) (ok bool, errs []*validation.Error) {

	var dType DeviceType

	switch v := properties.(type) {
	case *DevModBusConfig:
		dType = DevTypeModbus
		ok, errs = v.Valid()
	case *DevSmartBusConfig:
		dType = DevTypeSmartBus
		ok, errs = v.Valid()
	case *DevCommandConfig:
		dType = DevTypeCommand
		ok, errs = v.Valid()
	default:
		dType = DevTypeDefault
		ok = true
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

func (d *Device) SetPropertiesFromMap(properties map[string]interface{}) (ok bool, errs []*validation.Error, err error) {

	var out interface{}
	switch d.Type {
	case DevTypeModbus:
		out = &DevModBusConfig{}
	case DevTypeSmartBus:
		out = &DevSmartBusConfig{}
	case DevTypeCommand:
		out = &DevCommandConfig{}
	default:
		log.Warningf("unknown device config %v", d.Type)
		err = fmt.Errorf("unknown device config %v", d.Type)
		return
	}

	if err = mapstructure.Decode(properties, out); err != nil {
		return
	}

	ok, errs = d.SetProperties(out)

	return
}