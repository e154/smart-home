// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package models

import (
	"encoding/json"
	"fmt"
	. "github.com/e154/smart-home/common"
	. "github.com/e154/smart-home/models/devices"
	"github.com/e154/smart-home/system/validation"
	"github.com/mitchellh/mapstructure"
	"time"
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
	case DevTypeModbusRtu:
		out = &DevModBusRtuConfig{}
	case DevTypeModbusTcp:
		out = &DevModBusTcpConfig{}
	case DevTypeSmartBus:
		out = &DevSmartBusConfig{}
	case DevTypeCommand:
		out = &DevCommandConfig{}
	case DevTypeMqtt:
		out = &DevMqttConfig{}
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
	case *DevModBusRtuConfig:
		ok, errs = v.Valid()
	case *DevModBusTcpConfig:
		ok, errs = v.Valid()
	case *DevSmartBusConfig:
		ok, errs = v.Valid()
	case *DevCommandConfig:
		ok, errs = v.Valid()
	case *DevMqttConfig:
		ok, errs = v.Valid()
	}

	return
}

func (d *Device) SetProperties(properties interface{}) (ok bool, errs []*validation.Error) {

	var dType DeviceType

	switch v := properties.(type) {
	case *DevModBusRtuConfig:
		dType = DevTypeModbusRtu
		ok, errs = v.Valid()
	case *DevModBusTcpConfig:
		dType = DevTypeModbusTcp
		ok, errs = v.Valid()
	case *DevSmartBusConfig:
		dType = DevTypeSmartBus
		ok, errs = v.Valid()
	case *DevMqttConfig:
		dType = DevTypeMqtt
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
	case DevTypeModbusRtu:
		out = &DevModBusRtuConfig{}
	case DevTypeModbusTcp:
		out = &DevModBusTcpConfig{}
	case DevTypeSmartBus:
		out = &DevSmartBusConfig{}
	case DevTypeCommand:
		out = &DevCommandConfig{}
	case DevTypeMqtt:
		out = &DevMqttConfig{}
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
