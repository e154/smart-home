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

package zigbee2mqtt

import (
	"github.com/e154/smart-home/models"
	"sync"
)

// Device ...
type Device struct {
	friendlyName string
	modelLock    *sync.Mutex
	model        models.Zigbee2mqttDevice
}

// NewDevice ...
func NewDevice(friendlyName string, model *models.Zigbee2mqttDevice) *Device {
	return &Device{
		friendlyName: friendlyName,
		modelLock:    &sync.Mutex{},
		model:        *model,
	}
}

// AddFunc ...
func (d *Device) AddFunc(name string) {
	d.modelLock.Lock()
	defer d.modelLock.Unlock()
	for _, f := range d.model.Functions {
		if name == f {
			return
		}
	}
	d.model.Functions = append(d.model.Functions, name)
}

// DeviceType ...
func (d *Device) DeviceType(devType string) {
	d.modelLock.Lock()
	defer d.modelLock.Unlock()

	if d.model.Type != "" && devType == "sensor" {
		return
	}
	d.model.Type = devType
}

// GetModel ...
func (d *Device) GetModel() models.Zigbee2mqttDevice {
	d.modelLock.Lock()
	defer d.modelLock.Unlock()
	return d.model
}

// SetStatus ...
func (d *Device) SetStatus(status string) {
	d.modelLock.Lock()
	d.model.Status = status
	d.modelLock.Unlock()
}

// Status ...
func (d *Device) Status() string {
	d.modelLock.Lock()
	defer d.modelLock.Unlock()
	return d.model.Status
}

// SetName ...
func (d *Device) SetName(name string) {
	d.modelLock.Lock()
	d.model.Name = name
	d.modelLock.Unlock()
}

// SetModel ...
func (d *Device) SetModel(model string) {
	d.modelLock.Lock()
	d.model.Model = model
	d.modelLock.Unlock()
}

// SetDescription ...
func (d *Device) SetDescription(desc string) {
	d.modelLock.Lock()
	d.model.Description = desc
	d.modelLock.Unlock()
}

// SetVendor ...
func (d *Device) SetVendor(vendor string) {
	d.modelLock.Lock()
	d.model.Manufacturer = vendor
	d.modelLock.Unlock()
}

// GetImage ...
func (d *Device) GetImage() string {
	d.modelLock.Lock()
	defer d.modelLock.Unlock()
	d.model.GetImageUrl()
	return d.model.ImageUrl
}
