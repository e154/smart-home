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

package dashboard_models

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
	"reflect"
)

type Devices struct {
	Total    int64                              `json:"total"`
	Status   map[int64]*m.DashboardDeviceStatus `json:"status"`
	adaptors *adaptors.Adaptors                 `json:"-"`
	core     *core.Core                         `json:"-"`
}

func NewDevices(adaptors *adaptors.Adaptors,
	core *core.Core) *Devices {
	return &Devices{
		Status:   make(map[int64]*m.DashboardDeviceStatus),
		adaptors: adaptors,
		core:     core,
	}
}

func (d *Devices) Update() {

	var err error
	if _, d.Total, err = d.adaptors.Device.List(999, 0, "", ""); err != nil {
		log.Error(err.Error())
	}

	d.Status = d.core.Map.GetDevicesStates()
}

func (d *Devices) Broadcast() (map[string]interface{}, bool) {

	d.Update()

	return map[string]interface{}{
		"devices": d,
	}, true
}

func (d *Devices) BroadcastOne(deviceId int64, elementName string) (map[string]interface{}, bool) {

	d.Status = d.core.Map.GetDevicesStates()
	state, ok := d.Status[deviceId]
	if !ok {
		return nil, false
	}

	return map[string]interface{}{
		"device":  map[string]interface{}{"id": deviceId, "state": state},
		"devices": d,
	}, true
}

// only on request: 'dashboard.get.devices.states'
//
func (d *Devices) streamGetDevicesStates(client *stream.Client, value interface{}) {

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	d.Update()

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "states": d.Status})

	client.Send <- msg
}
