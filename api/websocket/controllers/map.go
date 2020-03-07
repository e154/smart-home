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

package controllers

import (
	mapModels "github.com/e154/smart-home/api/websocket/controllers/map_models"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/stream"
)

type ControllerMap struct {
	*ControllerCommon
	devices *mapModels.Devices
}

func NewControllerMap(common *ControllerCommon) *ControllerMap {
	return &ControllerMap{
		ControllerCommon: common,
		devices:          mapModels.NewDevices(common.metric),
	}
}

func (c *ControllerMap) Start() {
	c.metric.Subscribe("map", c)
	c.stream.Subscribe("map.get.devices", c.devices.GetDevicesStates)
	c.gate.Subscribe("map.get.devices", c.devices.GetDevicesStates)
}

func (c *ControllerMap) Stop() {
	c.metric.UnSubscribe("map")
	c.stream.UnSubscribe("map.get.devices")
	c.gate.UnSubscribe("map.get.devices")
}

// method callable from metric service
func (c *ControllerMap) Broadcast(param interface{}) {

	var body map[string]interface{}
	var ok bool

	switch v := param.(type) {
	case metrics.MapElementCursor:
		body, ok = c.devices.UpdateMapElement(v)
	}

	if ok {
		c.sendMsg(body)
	}
}

func (t *ControllerMap) sendMsg(payload map[string]interface{}) {

	msg := &stream.Message{
		Command: "map.telemetry",
		Type:    stream.Broadcast,
		Forward: stream.Request,
		Payload: payload,
	}

	t.stream.Broadcast(msg.Pack())
}
