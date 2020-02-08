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
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Device struct {
	table *db.Devices
	db    *gorm.DB
}

func GetDeviceAdaptor(d *gorm.DB) *Device {
	return &Device{
		table: &db.Devices{Db: d},
		db:    d,
	}
}

func (n *Device) Add(device *m.Device) (id int64, err error) {

	dbDevice := n.toDb(device)
	if id, err = n.table.Add(dbDevice); err != nil {
		return
	}

	return
}

func (n *Device) GetAllEnabled() (list []*m.Device, err error) {

	var dbList []*db.Device
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.Device, 0)
	for _, dbDevice := range dbList {
		device := n.fromDb(dbDevice)
		list = append(list, device)
	}

	return
}

func (n *Device) GetById(deviceId int64) (device *m.Device, err error) {

	var dbDevice *db.Device
	if dbDevice, err = n.table.GetById(deviceId); err != nil {
		return
	}

	device = n.fromDb(dbDevice)

	return
}

func (n *Device) Update(device *m.Device) (err error) {
	dbDevice := n.toDb(device)
	err = n.table.Update(dbDevice)
	return
}

func (n *Device) Delete(deviceId int64) (err error) {
	err = n.table.Delete(deviceId)
	return
}

func (n *Device) List(limit, offset int64, orderBy, sort string) (list []*m.Device, total int64, err error) {
	var dbList []*db.Device
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Device, 0)
	for _, dbDevice := range dbList {
		device := n.fromDb(dbDevice)
		list = append(list, device)
	}

	return
}

func (n *Device) Search(query string, limit, offset int) (list []*m.Device, total int64, err error) {
	var dbList []*db.Device
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Device, 0)
	for _, dbDevice := range dbList {
		dev := n.fromDb(dbDevice)
		list = append(list, dev)
	}

	return
}

func (n *Device) fromDb(dbDevice *db.Device) (device *m.Device) {
	device = &m.Device{
		Id:          dbDevice.Id,
		Name:        dbDevice.Name,
		Status:      dbDevice.Status,
		Description: dbDevice.Description,
		Type:        dbDevice.Type,
		Properties:  dbDevice.Properties,
		IsGroup:     dbDevice.IsGroup,
		Actions:     make([]*m.DeviceAction, 0),
		States:      make([]*m.DeviceState, 0),
		Devices:     make([]*m.Device, 0),
		CreatedAt:   dbDevice.CreatedAt,
		UpdatedAt:   dbDevice.UpdatedAt,
	}

	// parent device
	if dbDevice.Device != nil && dbDevice.DeviceId.Valid {
		device.Device = n.fromDb(dbDevice.Device)
		device.DeviceId = &dbDevice.DeviceId.Int64
	}

	// actions
	deviceActionAdaptor := GetDeviceActionAdaptor(n.db)
	for _, dbAction := range dbDevice.Actions {
		action := deviceActionAdaptor.fromDb(dbAction)
		device.Actions = append(device.Actions, action)
	}

	// states
	deviceStatesAdaptor := GetDeviceStateAdaptor(n.db)
	for _, dbState := range dbDevice.States {
		state := deviceStatesAdaptor.fromDb(dbState)
		device.States = append(device.States, state)
	}

	// devices
	for _, dbDevice := range dbDevice.Devices {
		dev := n.fromDb(dbDevice)
		device.Devices = append(device.Devices, dev)
	}

	// node
	if dbDevice.Node != nil {
		nodeAdaptor := GetNodeAdaptor(n.db)
		device.Node = nodeAdaptor.fromDb(dbDevice.Node)
	}

	return
}

func (n *Device) toDb(device *m.Device) (dbDevice *db.Device) {
	dbDevice = &db.Device{
		Id:          device.Id,
		Name:        device.Name,
		Status:      device.Status,
		Description: device.Description,
		Properties:  device.Properties,
		Type:        device.Type,
		IsGroup:     device.IsGroup,
	}

	// device
	if device.Device != nil && device.Device.Id != 0 {
		dbDevice.DeviceId.Scan(device.Device.Id)
	}

	// node
	if device.Node != nil && device.Node.Id != 0 {
		dbDevice.NodeId.Scan(device.Node.Id)
	}

	return
}
