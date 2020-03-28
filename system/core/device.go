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
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/models/devices"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/zigbee2mqtt"
)

type Device struct {
	dev         *m.Device
	node        *Node
	mqtt        *mqtt.Mqtt
	adaptors    *adaptors.Adaptors
	zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt
}

func NewDevice(dev *m.Device, node *Node, mqtt *mqtt.Mqtt,
	adaptors *adaptors.Adaptors, zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt) *Device {
	return &Device{
		dev:         dev,
		node:        node,
		mqtt:        mqtt,
		adaptors:    adaptors,
		zigbee2mqtt: zigbee2mqtt,
	}
}

// run command
func (d Device) RunCommand(name string, args []string) (result DevCommandResponse) {

	request := &DevCommandRequest{
		Name: name,
		Args: args,
	}

	result = DevCommandResponse{}
	data, err := json.Marshal(request)
	if err != nil {
		result.Error = err.Error()
		return
	}

	if d.node == nil {
		result.Error = "node is nil"
		return
	}

	nodeResult, err := d.node.Send(d.dev, data)
	if err != nil {
		result.Error = err.Error()
		return
	}

	if nodeResult.Response == nil || len(nodeResult.Response) == 0 {
		result.Error = "return nil result from response"
		return
	}

	if err = json.Unmarshal(nodeResult.Response, &result); err != nil {
		result.Error = err.Error()
		return
	}

	//debug.Println(nodeResult)

	result.Time = nodeResult.Time

	return
}

func (d Device) SmartBus(command []byte) (result DevSmartBusResponse) {

	request := &DevSmartBusRequest{
		Command: command,
	}

	result = DevSmartBusResponse{}
	data, err := json.Marshal(request)
	if err != nil {
		result.Error = err.Error()
		return
	}

	if d.node == nil {
		result.Error = "node is nil"
		return
	}

	nodeResult, err := d.node.Send(d.dev, data)

	if err = json.Unmarshal(nodeResult.Response, &result); err != nil {
		result.Error = err.Error()
		return
	}

	//debug.Println(nodeResult)
	//debug.Println(result)

	result.Time = nodeResult.Time

	return
}

func (d Device) ModBus(f string, address, count uint16, command []uint16) (result DevModBusResponse) {

	result = DevModBusResponse{}

	if d.node == nil {
		result.Error = "node is nil"
		return
	}

	if d.dev.Type == DevTypeModbusRtu && d.node.stat.Thread == 0 {
		result.Error = "no serial device"
		return
	}

	request := &DevModBusRequest{
		Function: f,
		Address:  address,
		Count:    count,
		Command:  command,
	}

	data, err := json.Marshal(request)
	if err != nil {
		result.Error = err.Error()
		log.Error(err.Error())
		return
	}

	nodeResult, err := d.node.Send(d.dev, data)
	if err != nil {
		result.Error = err.Error()
		//log.Error(err.Error())
		return
	}

	if err = json.Unmarshal(nodeResult.Response, &result); err != nil {
		result.Error = err.Error()
		log.Error(err.Error())
		return
	}

	//debug.Println(nodeResult)
	//debug.Println(result)

	result.Time = nodeResult.Time

	return
}

func (d Device) Zigbee2mqtt(path string, payload []byte) (result DevZigbee2mqttResponse) {

	if d.dev.Type != DevTypeZigbee2mqtt {
		result.Error = "no zigbee settings"
		return
	}

	params := &DevZigbee2mqttConfig{}

	b, _ := d.dev.Properties.MarshalJSON()
	if err := json.Unmarshal(b, params); err != nil {
		result.Error = err.Error()
		return
	}

	device, err := d.adaptors.Zigbee2mqttDevice.GetById(params.Zigbee2mqttDeviceId)
	if err != nil {
		result.Error = err.Error()
		return
	}

	//TODO add cache
	topic, err := d.zigbee2mqtt.GetTopicByDevice(device)
	if err != nil {
		result.Error = err.Error()
		return
	}

	if path != "" {
		if path[:1] != "/" {
			topic = fmt.Sprintf("%s/%s", topic, path)
		} else {
			topic += path
		}
	}

	d.mqtt.Publish(topic, payload, 0, false)

	return
}

func (d Device) Mqtt(path string, payload []byte) (result DevMqttResponse) {

	if d.dev.Type != DevTypeMqtt {
		result.Error = "no mqtt settings"
		return
	}

	b, _ := d.dev.Properties.MarshalJSON()

	params := &DevMqttConfig{}
	_ = json.Unmarshal(b, params)

	topic := params.Address

	if path != "" {
		if path[:1] != "/" {
			topic = fmt.Sprintf("%s/%s", topic, path)
		} else {
			topic += path
		}
	}

	d.mqtt.Publish(topic, payload, 0, false)

	return
}
