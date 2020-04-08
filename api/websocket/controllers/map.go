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
	"bytes"
	"encoding/json"
	mapModels "github.com/e154/smart-home/api/websocket/controllers/map_models"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/stream"
	"sync"
)

// ControllerMap ...
type ControllerMap struct {
	*ControllerCommon
	devices  *mapModels.Devices
	sendLock *sync.Mutex
	buf      *bytes.Buffer
	enc      *json.Encoder
}

// NewControllerMap ...
func NewControllerMap(common *ControllerCommon) *ControllerMap {
	buf := bytes.NewBuffer(nil)
	return &ControllerMap{
		ControllerCommon: common,
		devices:          mapModels.NewDevices(common.metric),
		sendLock:         &sync.Mutex{},
		buf:              buf,
		enc:              json.NewEncoder(buf),
	}
}

// Start ...
func (c *ControllerMap) Start() {
	c.metric.Subscribe("map", c)
	c.stream.Subscribe("map.get.devices", c.devices.GetDevicesStates)
}

// Stop ...
func (c *ControllerMap) Stop() {
	c.metric.UnSubscribe("map")
	c.stream.UnSubscribe("map.get.devices")
}

// method callable from metric service
func (c *ControllerMap) Broadcast(param interface{}) {

	var body map[string]interface{}
	var ok bool

	switch v := param.(type) {
	case string:

	case metrics.MapElementCursor:
		body, ok = c.devices.Broadcast(v)
	default:
		log.Warnf("unknown type %v", v)
	}

	if !ok {
		return
	}

	c.sendMsg(body)
}

func (t *ControllerMap) sendMsg(payload map[string]interface{}) (err error) {

	t.sendLock.Lock()
	defer t.sendLock.Unlock()

	msg := &stream.Message{
		Command: "map.telemetry",
		Type:    stream.Broadcast,
		Forward: stream.Request,
		Payload: payload,
	}

	t.buf.Reset()
	if err = t.enc.Encode(msg); err != nil {
		log.Error(err.Error())
		return
	}

	data := make([]byte, t.buf.Len())
	copy(data, t.buf.Bytes())
	t.stream.Broadcast(data)

	return
}
