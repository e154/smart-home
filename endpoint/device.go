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

package endpoint

import (
	"errors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/validation"
)

// DeviceEndpoint ...
type DeviceEndpoint struct {
	*CommonEndpoint
}

// NewDeviceEndpoint ...
func NewDeviceEndpoint(common *CommonEndpoint) *DeviceEndpoint {
	return &DeviceEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (d *DeviceEndpoint) Add(params *m.Device) (device *m.Device, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = d.adaptors.Device.Add(params); err != nil {
		return
	}

	if device, err = d.adaptors.Device.GetById(id); err != nil {
		return
	}

	var disabled int64
	if device.Status == "disabled" {
		disabled++
	}
	d.metric.Update(metrics.DeviceAdd{TotalNum: 1, DisabledNum: disabled})

	return
}

// GetById ...
func (d *DeviceEndpoint) GetById(deviceId int64) (device *m.Device, err error) {
	device, err = d.adaptors.Device.GetById(deviceId)
	return
}

// Update ...
func (d *DeviceEndpoint) Update(device *m.Device) (result *m.Device, errs []*validation.Error, err error) {

	_, errs = device.Valid()
	if len(errs) > 0 {
		return
	}

	if err = d.adaptors.Device.Update(device); err != nil {
		return
	}

	result, err = d.adaptors.Device.GetById(device.Id)

	return
}

// GetList ...
func (d *DeviceEndpoint) GetList(limit, offset int64, order, sortBy string) (devices []*m.Device, total int64, err error) {

	devices, total, err = d.adaptors.Device.List(limit, offset, order, sortBy)

	return
}

// Delete ...
func (d *DeviceEndpoint) Delete(deviceId int64) (err error) {

	if deviceId == 0 {
		err = errors.New("device id is null")
		return
	}

	var device *m.Device
	if device, err = d.adaptors.Device.GetById(deviceId); err != nil {
		return
	}

	if err = d.adaptors.Device.Delete(device.Id); err != nil {
		return
	}

	var disabled int64
	if device.Status == "disabled" {
		disabled++
	}

	d.metric.Update(metrics.DeviceDelete{TotalNum: 1, DisabledNum: disabled})

	return
}

// Search ...
func (d *DeviceEndpoint) Search(query string, limit, offset int) (devices []*m.Device, total int64, err error) {

	devices, total, err = d.adaptors.Device.Search(query, limit, offset)

	return
}
