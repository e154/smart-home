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

// DeviceAction ...
type DeviceAction struct {
	table *db.DeviceActions
	db    *gorm.DB
}

// GetDeviceActionAdaptor ...
func GetDeviceActionAdaptor(d *gorm.DB) *DeviceAction {
	return &DeviceAction{
		table: &db.DeviceActions{Db: d},
		db:    d,
	}
}

// Add ...
func (n *DeviceAction) Add(device *m.DeviceAction) (id int64, err error) {

	dbDeviceAction := n.toDb(device)
	if id, err = n.table.Add(dbDeviceAction); err != nil {
		return
	}

	return
}

// GetById ...
func (n *DeviceAction) GetById(actionId int64) (device *m.DeviceAction, err error) {

	var dbDeviceAction *db.DeviceAction
	if dbDeviceAction, err = n.table.GetById(actionId); err != nil {
		return
	}

	device = n.fromDb(dbDeviceAction)

	return
}

// GetByDeviceId ...
func (n *DeviceAction) GetByDeviceId(deviceId int64) (actions []*m.DeviceAction, err error) {

	var dbDeviceActions []*db.DeviceAction
	if dbDeviceActions, err = n.table.GetByDeviceId(deviceId); err != nil {
		return
	}

	actions = make([]*m.DeviceAction, 0)
	for _, dbActino := range dbDeviceActions {
		action := n.fromDb(dbActino)
		actions = append(actions, action)
	}

	return
}

// Update ...
func (n *DeviceAction) Update(device *m.DeviceAction) (err error) {
	dbDeviceAction := n.toDb(device)
	err = n.table.Update(dbDeviceAction)
	return
}

// Delete ...
func (n *DeviceAction) Delete(deviceId int64) (err error) {
	err = n.table.Delete(deviceId)
	return
}

// List ...
func (n *DeviceAction) List(limit, offset int64, orderBy, sort string) (list []*m.DeviceAction, total int64, err error) {
	var dbList []*db.DeviceAction
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.DeviceAction, 0)
	for _, dbDeviceAction := range dbList {
		device := n.fromDb(dbDeviceAction)
		list = append(list, device)
	}

	return
}

// Search ...
func (n *DeviceAction) Search(query string, limit, offset int) (list []*m.DeviceAction, total int64, err error) {
	var dbList []*db.DeviceAction
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.DeviceAction, 0)
	for _, dbDeviceAction := range dbList {
		ver := n.fromDb(dbDeviceAction)
		list = append(list, ver)
	}

	return
}

func (n *DeviceAction) fromDb(dbDeviceAction *db.DeviceAction) (device *m.DeviceAction) {

	device = &m.DeviceAction{
		Id:          dbDeviceAction.Id,
		Name:        dbDeviceAction.Name,
		Description: dbDeviceAction.Description,
		DeviceId:    dbDeviceAction.DeviceId,
		ScriptId:    dbDeviceAction.ScriptId,
		CreatedAt:   dbDeviceAction.CreatedAt,
		UpdatedAt:   dbDeviceAction.UpdatedAt,
	}

	// device
	if dbDeviceAction.Device != nil {
		deviceAdaptor := GetDeviceAdaptor(n.db)
		device.Device = deviceAdaptor.fromDb(dbDeviceAction.Device)
	}

	// script
	if dbDeviceAction.Script != nil {
		scriptADaptor := GetScriptAdaptor(n.db)
		device.Script, _ = scriptADaptor.fromDb(dbDeviceAction.Script)
	}

	return
}

func (n *DeviceAction) toDb(device *m.DeviceAction) (dbDeviceAction *db.DeviceAction) {
	dbDeviceAction = &db.DeviceAction{
		Id:          device.Id,
		Name:        device.Name,
		Description: device.Description,
		DeviceId:    device.DeviceId,
		ScriptId:    device.ScriptId,
	}
	return
}
