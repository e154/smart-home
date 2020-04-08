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

// MapDevice ...
type MapDevice struct {
	table *db.MapDevices
	db    *gorm.DB
}

// GetMapDeviceAdaptor ...
func GetMapDeviceAdaptor(d *gorm.DB) *MapDevice {
	return &MapDevice{
		table: &db.MapDevices{Db: d},
		db:    d,
	}
}

// Add ...
func (n *MapDevice) Add(ver *m.MapDevice) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

// GetById ...
func (n *MapDevice) GetById(mapId int64) (ver *m.MapDevice, err error) {

	var dbVer *db.MapDevice
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Delete ...
func (n *MapDevice) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

func (n *MapDevice) fromDb(dbVer *db.MapDevice) (ver *m.MapDevice) {
	ver = &m.MapDevice{
		Id:        dbVer.Id,
		DeviceId:  dbVer.DeviceId,
		ImageId:   dbVer.ImageId,
		Actions:   make([]*m.MapDeviceAction, 0),
		States:    make([]*m.MapDeviceState, 0),
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
	}

	// actions
	mapDeviceActionAdaptor := GetMapDeviceActionAdaptor(n.db)
	for _, dbAction := range dbVer.Actions {
		action := mapDeviceActionAdaptor.fromDb(dbAction)
		ver.Actions = append(ver.Actions, action)
	}

	// states
	mapDeviceStateAdaptor := GetMapDeviceStateAdaptor(n.db)
	for _, dbState := range dbVer.States {
		state := mapDeviceStateAdaptor.fromDb(dbState)
		ver.States = append(ver.States, state)
	}

	// device
	if dbVer.Device != nil {
		deviceAdaptor := GetDeviceAdaptor(n.db)
		ver.Device = deviceAdaptor.fromDb(dbVer.Device)
	}

	// image
	if dbVer.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		ver.Image = imageAdaptor.fromDb(dbVer.Image)
	}

	return
}

func (n *MapDevice) toDb(ver *m.MapDevice) (dbVer *db.MapDevice) {
	dbVer = &db.MapDevice{
		Id:       ver.Id,
		DeviceId: ver.DeviceId,
	}
	if ver.ImageId != 0 {
		dbVer.ImageId = ver.ImageId
	} else if ver.Image != nil && ver.Image.Id != 0 {
		dbVer.ImageId = ver.Image.Id
	}
	return
}
