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

package core

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/models/devices"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/zigbee2mqtt"
)

// Javascript Binding
//
// Device
//	.GetName()
//	.GetModel()
//	.GetDescription()
//	.RunCommand(command []string)
//	.SmartBus(command []byte)
//	.ModBus(func string, address, count int64, command []byte)
//	.Zigbee2mqtt(path, payload)
//	.Mqtt(path, payload)
//	.Send(params interface)
//
type DeviceBind struct {
	model    *m.Device
	node     *Node
	mqtt     *mqtt.Mqtt
	adaptors *adaptors.Adaptors
	device   *Device
}

// NewDeviceBind ...
func NewDeviceBind(model *m.Device, node *Node, mqtt *mqtt.Mqtt, adaptors *adaptors.Adaptors, zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt) *DeviceBind {
	return &DeviceBind{
		model: model,
		node:  node,
		mqtt:  mqtt, adaptors: adaptors,
		device: NewDevice(model, node, mqtt, adaptors, zigbee2mqtt),
	}
}

// GetName ...
func (d *DeviceBind) GetName() string {
	return d.model.Name
}

// GetModel ...
func (d *DeviceBind) GetModel() *m.Device {
	return d.model
}

// GetDescription ...
func (d *DeviceBind) GetDescription() string {
	return d.model.Description
}

// RunCommand ...
func (d *DeviceBind) RunCommand(name string, args []string) (result DevCommandResponse) {
	result = d.device.RunCommand(name, args)
	return
}

// SmartBus ...
func (d *DeviceBind) SmartBus(command []byte) (result DevSmartBusResponse) {
	result = d.device.SmartBus(command)
	return
}

// ModBus ...
func (d *DeviceBind) ModBus(f string, address, count uint16, command []uint16) (result DevModBusResponse) {
	result = d.device.ModBus(f, address, count, command)
	return
}

// Zigbee2mqtt ...
func (d *DeviceBind) Zigbee2mqtt(path string, payload []byte) (result DevZigbee2mqttResponse) {
	result = d.device.Zigbee2mqtt(path, payload)
	return
}

// Mqtt ...
func (d *DeviceBind) Mqtt(path string, payload []byte) (result DevMqttResponse) {
	result = d.device.Mqtt(path, payload)
	return
}

// Send ...
func (d *DeviceBind) Send(payload interface{}) (result interface{}) {

	switch v := payload.(type) {
	case RunCommandBind:
		result = d.device.RunCommand(v.Name, v.Args)
	case ModBusBind:
		result = d.device.ModBus(v.F, v.Address, v.Count, v.Command)
	case Zigbee2mqttBind:
		result = d.device.Zigbee2mqtt(v.Path, v.Payload)
	case MqttBind:
		result = d.device.Mqtt(v.Path, v.Payload)
	default:
		fmt.Printf("unknown command type %v\n", v)
	}

	return
}
