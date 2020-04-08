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
	"github.com/e154/smart-home/system/validation"
)

// DeviceActionEndpoint ...
type DeviceActionEndpoint struct {
	*CommonEndpoint
}

// NewDeviceActionEndpoint ...
func NewDeviceActionEndpoint(common *CommonEndpoint) *DeviceActionEndpoint {
	return &DeviceActionEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (d *DeviceActionEndpoint) Add(params *m.DeviceAction) (action *m.DeviceAction, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = d.adaptors.DeviceAction.Add(params); err != nil {
		return
	}

	action, err = d.adaptors.DeviceAction.GetById(id)

	return
}

// Update ...
func (d *DeviceActionEndpoint) Update(params *m.DeviceAction) (result *m.DeviceAction, errs []*validation.Error, err error) {

	var action *m.DeviceAction
	if action, err = d.adaptors.DeviceAction.GetById(params.Id); err != nil {
		return
	}

	action.Name = params.Name
	action.Description = params.Description
	action.ScriptId = params.ScriptId

	// validation
	_, errs = action.Valid()
	if len(errs) > 0 {
		return
	}

	if err = d.adaptors.DeviceAction.Update(action); err != nil {
		return
	}

	action, err = d.adaptors.DeviceAction.GetById(params.Id)

	return
}

// GetById ...
func (d *DeviceActionEndpoint) GetById(id int64) (device *m.DeviceAction, err error) {

	device, err = d.adaptors.DeviceAction.GetById(id)

	return
}

// Delete ...
func (d *DeviceActionEndpoint) Delete(id int64) (err error) {

	if id == 0 {
		err = errors.New("action id is null")
		return
	}

	var device *m.DeviceAction
	if device, err = d.adaptors.DeviceAction.GetById(id); err != nil {
		return
	}

	err = d.adaptors.DeviceAction.Delete(device.Id)

	return
}

// GetList ...
func (d *DeviceActionEndpoint) GetList(deviceId int64) (actions []*m.DeviceAction, err error) {

	actions, err = d.adaptors.DeviceAction.GetByDeviceId(deviceId)

	return
}

// Search ...
func (d *DeviceActionEndpoint) Search(query string, limit, offset int) (actions []*m.DeviceAction, total int64, err error) {

	actions, total, err = d.adaptors.DeviceAction.Search(query, limit, offset)

	return
}
