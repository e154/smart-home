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

package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

// DeviceState ...
type DeviceState struct {
	table *db.DeviceStates
	db    *gorm.DB
}

// GetDeviceStateAdaptor ...
func GetDeviceStateAdaptor(d *gorm.DB) *DeviceState {
	return &DeviceState{
		table: &db.DeviceStates{Db: d},
		db:    d,
	}
}

// Add ...
func (n *DeviceState) Add(device *m.DeviceState) (id int64, err error) {

	dbDeviceState := n.toDb(device)
	if id, err = n.table.Add(dbDeviceState); err != nil {
		return
	}

	return
}

// GetById ...
func (n *DeviceState) GetById(deviceId int64) (device *m.DeviceState, err error) {

	var dbDeviceState *db.DeviceState
	if dbDeviceState, err = n.table.GetById(deviceId); err != nil {
		return
	}

	device = n.fromDb(dbDeviceState)

	return
}

// GetByDeviceId ...
func (n *DeviceState) GetByDeviceId(deviceId int64) (states []*m.DeviceState, err error) {

	var dbDeviceStates []*db.DeviceState
	if dbDeviceStates, err = n.table.GetByDeviceId(deviceId); err != nil {
		return
	}

	states = make([]*m.DeviceState, 0)
	for _, dbActino := range dbDeviceStates {
		state := n.fromDb(dbActino)
		states = append(states, state)
	}

	return
}

// Update ...
func (n *DeviceState) Update(device *m.DeviceState) (err error) {
	dbDeviceState := n.toDb(device)
	err = n.table.Update(dbDeviceState)
	return
}

// Delete ...
func (n *DeviceState) Delete(deviceId int64) (err error) {
	err = n.table.Delete(deviceId)
	return
}

// List ...
func (n *DeviceState) List(limit, offset int64, orderBy, sort string) (list []*m.DeviceState, total int64, err error) {
	var dbList []*db.DeviceState
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.DeviceState, 0)
	for _, dbDeviceState := range dbList {
		device := n.fromDb(dbDeviceState)
		list = append(list, device)
	}

	return
}

func (n *DeviceState) fromDb(dbVer *db.DeviceState) (ver *m.DeviceState) {
	ver = &m.DeviceState{
		Id:          dbVer.Id,
		Description: dbVer.Description,
		SystemName:  dbVer.SystemName,
		DeviceId:    dbVer.DeviceId,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}

	if dbVer.Device != nil {
		ver.DeviceId = dbVer.Device.Id
		deviceAdaptor := GetDeviceAdaptor(n.db)
		ver.Device = deviceAdaptor.fromDb(dbVer.Device)
	}
	return
}

func (n *DeviceState) toDb(device *m.DeviceState) (dbDeviceState *db.DeviceState) {
	dbDeviceState = &db.DeviceState{
		Id:          device.Id,
		Description: device.Description,
		DeviceId:    device.DeviceId,
		SystemName:  device.SystemName,
	}
	return
}
