// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package metrics

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"sync"
)

// Device ...
type Device struct {
	Total    int64 `json:"total"`
	Disabled int64 `json:"disabled"`
}

// DeviceManager ...
type DeviceManager struct {
	publisher  IPublisher
	total      metrics.Counter
	disabled   metrics.Counter
	updateLock sync.Mutex
}

// NewDeviceManager ...
func NewDeviceManager(publisher IPublisher) *DeviceManager {
	return &DeviceManager{
		total:     metrics.NewCounter(),
		disabled:  metrics.NewCounter(),
		publisher: publisher,
	}
}

func (d *DeviceManager) update(t interface{}) {
	switch v := t.(type) {
	case DeviceAdd:
		d.total.Inc(v.TotalNum)
		d.disabled.Inc(v.DisabledNum)
	case DeviceDelete:
		d.total.Dec(v.TotalNum)
		d.disabled.Dec(v.DisabledNum)
	default:
		return
	}

	d.broadcast()
}

func (d *DeviceManager) broadcast() {
	go d.publisher.Broadcast("device")
}

// Snapshot ...
func (d *DeviceManager) Snapshot() Device {
	return Device{
		Total:    d.total.Count(),
		Disabled: d.disabled.Count(),
	}
}

func (d *DeviceManager) key(deviceId int64, elementName string) string {
	return fmt.Sprintf("%d_%s", deviceId, elementName)
}

// DeviceAdd ...
type DeviceAdd struct {
	TotalNum    int64
	DisabledNum int64
}

// DeviceDelete ...
type DeviceDelete struct {
	TotalNum    int64
	DisabledNum int64
}
