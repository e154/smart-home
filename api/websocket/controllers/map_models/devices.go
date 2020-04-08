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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/stream"
)

var (
	log = common.MustGetLogger("map_models")
)

// Devices ...
type Devices struct {
	metric *metrics.MetricManager
}

// NewDevices ...
func NewDevices(metric *metrics.MetricManager) *Devices {
	return &Devices{
		metric: metric,
	}
}

// Broadcast ...
func (d *Devices) Broadcast(cursor metrics.MapElementCursor) (map[string]interface{}, bool) {

	snapshot := d.metric.MapElement.Snapshot()
	element, ok := snapshot.Elements[d.key(cursor.DeviceId, cursor.ElementName)]
	if !ok {
		return nil, false
	}

	return map[string]interface{}{
		"device": map[string]interface{}{
			"id":           element.DeviceId,
			"element_name": element.ElementName,
			"options":      element.StateOptions,
			"status_id":    element.StateId,
		},
	}, true
}

// only on request: 'map.get.devices.states'
//
func (d *Devices) GetDevicesStates(client stream.IStreamClient, message stream.Message) {

	snapshot := d.metric.MapElement.Snapshot()

	devices := make([]MapElement, 0)

	for _, element := range snapshot.Elements {
		devices = append(devices, MapElement{
			Id:          element.DeviceId,
			StatusId:    element.StateId,
			Options:     element.StateOptions,
			ElementName: element.ElementName,
		})
	}

	msg := stream.Message{
		Id:      message.Id,
		Forward: stream.Response,
		Payload: map[string]interface{}{
			"devices": devices,
		},
	}

	client.Write(msg.Pack())
}

func (d *Devices) key(deviceId int64, elementName string) string {
	return fmt.Sprintf("%d_%s", deviceId, elementName)
}

// MapElement ...
type MapElement struct {
	Id          int64       `json:"id"`
	StatusId    int64       `json:"status_id"`
	Options     interface{} `json:"options"`
	ElementName string      `json:"element_name"`
}
