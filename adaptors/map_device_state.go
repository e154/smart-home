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
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

// MapDeviceState ...
type MapDeviceState struct {
	table *db.MapDeviceStates
	db    *gorm.DB
}

// GetMapDeviceStateAdaptor ...
func GetMapDeviceStateAdaptor(d *gorm.DB) *MapDeviceState {
	return &MapDeviceState{
		table: &db.MapDeviceStates{Db: d},
		db:    d,
	}
}

// Add ...
func (n *MapDeviceState) Add(ver *m.MapDeviceState) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

// DeleteByDeviceId ...
func (n *MapDeviceState) DeleteByDeviceId(deviceId int64) (err error) {
	err = n.table.DeleteByDeviceId(deviceId)
	return
}

// AddMultiple ...
func (n *MapDeviceState) AddMultiple(items []*m.MapDeviceState) (err error) {

	insertRecords := make([]interface{}, 0, len(items))

	for _, ver := range items {
		if ver.Image == nil {
			continue
		}
		insertRecords = append(insertRecords, n.toDb(ver))
	}

	err = gormbulk.BulkInsert(n.db, insertRecords, len(insertRecords))

	return
}

func (n *MapDeviceState) fromDb(dbVer *db.MapDeviceState) (ver *m.MapDeviceState) {
	ver = &m.MapDeviceState{
		Id:            dbVer.Id,
		DeviceStateId: dbVer.DeviceStateId,
		MapDeviceId:   dbVer.MapDeviceId,
		ImageId:       dbVer.ImageId,
		Style:         dbVer.Style,
		CreatedAt:     dbVer.CreatedAt,
		UpdatedAt:     dbVer.UpdatedAt,
	}

	// image
	if dbVer.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		ver.Image = imageAdaptor.fromDb(dbVer.Image)
	}

	// state
	if dbVer.DeviceState != nil {
		stateAdaptor := GetDeviceStateAdaptor(n.db)
		ver.DeviceState = stateAdaptor.fromDb(dbVer.DeviceState)
	}

	return
}

func (n *MapDeviceState) toDb(ver *m.MapDeviceState) (dbVer *db.MapDeviceState) {
	dbVer = &db.MapDeviceState{
		Id:            ver.Id,
		DeviceStateId: ver.DeviceStateId,
		MapDeviceId:   ver.MapDeviceId,
		ImageId:       ver.ImageId,
		Style:         ver.Style,
	}
	if ver.DeviceState != nil && ver.DeviceState.Id != 0 {
		dbVer.DeviceStateId = ver.DeviceState.Id
	}
	if ver.Image != nil && ver.Image.Id != 0 {
		dbVer.ImageId = ver.Image.Id
	}
	return
}
