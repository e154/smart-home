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
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/stream"
	"reflect"
)

type Devices struct {
	metric *metrics.MetricManager
}

func NewDevices(metric *metrics.MetricManager) *Devices {
	return &Devices{metric: metric}
}

func (d *Devices) Broadcast() (map[string]interface{}, bool) {

	return map[string]interface{}{
		"devices": d.metric.Device.Snapshot(),
	}, true
}

// only on request: 'dashboard.get.devices.states'
//
func (d *Devices) streamGetDevicesStates(client *stream.Client, value interface{}) {

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "states": d.metric.Device.Snapshot()})

	client.Send <- msg
}
