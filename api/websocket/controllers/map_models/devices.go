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

package map_models

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/websocket/controllers/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
	"github.com/op/go-logging"
	"sync"
)

var (
	log = logging.MustGetLogger("map_models")
)

type Devices struct {
	sync.Mutex
	Total       int64                          `json:"total"`
	DeviceStats map[string]*models.DeviceState `json:"device_stats"`
	adaptors    *adaptors.Adaptors
	core        *core.Core
}

func NewDevices(adaptors *adaptors.Adaptors,
	core *core.Core) *Devices {
	return &Devices{
		DeviceStats: make(map[string]*models.DeviceState),
		adaptors:    adaptors,
		core:        core,
	}
}

func (d *Devices) Update() {

	if d.core == nil {
		return
	}

	mapElements := d.core.Map.GetAllElements()

	d.Lock()
	d.Total = int64(len(mapElements))

	// clear map
	for key, _ := range d.DeviceStats {
		delete(d.DeviceStats, key)
	}

	var status *models.DeviceStateStatus
	for _, mapElement := range mapElements {

		if mapElement.State != nil {
			status = &models.DeviceStateStatus{
				Id:          mapElement.State.Id,
				DeviceId:    mapElement.State.DeviceId,
				SystemName:  mapElement.State.SystemName,
				Description: mapElement.State.Description,
			}
		}

		deviceId := mapElement.Device.Id
		key := d.key(deviceId, mapElement.ElementName)
		d.DeviceStats[key] = &models.DeviceState{
			Id:          deviceId,
			Status:      status,
			Options:     mapElement.Options,
			ElementName: mapElement.ElementName,
		}
	}
	d.Unlock()
}

func (d *Devices) Broadcast() (map[string]interface{}, bool) {

	d.Update()

	return map[string]interface{}{
		"devices": d,
	}, true
}

func (d *Devices) BroadcastOne(deviceId int64, elementName string) (map[string]interface{}, bool) {

	d.Update()

	d.Lock()
	key := d.key(deviceId, elementName)
	state, ok := d.DeviceStats[key]
	if !ok {
		d.Unlock()
		return nil, false
	}
	d.Unlock()

	return map[string]interface{}{
		"device": state,
	}, true
}

// only on request: 'map.get.devices.states'
//
func (d *Devices) GetDevicesStates(client stream.IStreamClient, message stream.Message) {

	d.Update()

	msg := &stream.Message{
		Id:      message.Id,
		Forward: stream.Response,
		Payload: map[string]interface{}{
			"states": d,
		},
	}

	client.Write(msg.Pack())
}

func (d *Devices) key(deviceId int64, elementName string) string {
	return fmt.Sprintf("%d_%s", deviceId, elementName)
}
