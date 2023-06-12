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

package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// IZigbee2mqttDevice ...
type IZigbee2mqttDevice interface {
	Add(ver *m.Zigbee2mqttDevice) (err error)
	GetById(id string) (ver *m.Zigbee2mqttDevice, err error)
	Update(ver *m.Zigbee2mqttDevice) (err error)
	Delete(id string) (err error)
	List(limit, offset int64) (list []*m.Zigbee2mqttDevice, total int64, err error)
	ListByBridgeId(bridgeId, limit, offset int64) (list []*m.Zigbee2mqttDevice, total int64, err error)
	Search(query string, limit, offset int64) (list []*m.Zigbee2mqttDevice, total int64, err error)
	fromDb(dbVer *db.Zigbee2mqttDevice) (ver *m.Zigbee2mqttDevice)
	toDb(ver *m.Zigbee2mqttDevice) (dbVer *db.Zigbee2mqttDevice)
}

// Zigbee2mqttDevice ...
type Zigbee2mqttDevice struct {
	IZigbee2mqttDevice
	table *db.Zigbee2mqttDevices
	db    *gorm.DB
}

// GetZigbee2mqttDeviceAdaptor ...
func GetZigbee2mqttDeviceAdaptor(d *gorm.DB) IZigbee2mqttDevice {
	return &Zigbee2mqttDevice{
		table: &db.Zigbee2mqttDevices{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Zigbee2mqttDevice) Add(ver *m.Zigbee2mqttDevice) (err error) {

	err = n.table.Add(n.toDb(ver))

	return
}

// GetById ...
func (n *Zigbee2mqttDevice) GetById(id string) (ver *m.Zigbee2mqttDevice, err error) {

	var dbVer *db.Zigbee2mqttDevice
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *Zigbee2mqttDevice) Update(ver *m.Zigbee2mqttDevice) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

// Delete ...
func (n *Zigbee2mqttDevice) Delete(id string) (err error) {
	err = n.table.Delete(id)
	return
}

// List ...
func (n *Zigbee2mqttDevice) List(limit, offset int64) (list []*m.Zigbee2mqttDevice, total int64, err error) {
	var dbList []*db.Zigbee2mqttDevice
	if dbList, total, err = n.table.List(int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Zigbee2mqttDevice, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// ListByBridgeId ...
func (n *Zigbee2mqttDevice) ListByBridgeId(bridgeId, limit, offset int64) (list []*m.Zigbee2mqttDevice, total int64, err error) {
	var dbList []*db.Zigbee2mqttDevice
	if dbList, total, err = n.table.ListByBridgeId(bridgeId, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Zigbee2mqttDevice, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// Search ...
func (n *Zigbee2mqttDevice) Search(query string, limit, offset int64) (list []*m.Zigbee2mqttDevice, total int64, err error) {
	var dbList []*db.Zigbee2mqttDevice
	if dbList, total, err = n.table.Search(query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Zigbee2mqttDevice, 0)
	for _, dbNode := range dbList {
		node := n.fromDb(dbNode)
		list = append(list, node)
	}

	return
}

func (n *Zigbee2mqttDevice) fromDb(dbVer *db.Zigbee2mqttDevice) (ver *m.Zigbee2mqttDevice) {
	ver = &m.Zigbee2mqttDevice{
		Id:            dbVer.Id,
		Name:          dbVer.Name,
		Type:          dbVer.Type,
		Zigbee2mqttId: dbVer.Zigbee2mqttId,
		Model:         dbVer.Model,
		Description:   dbVer.Description,
		Manufacturer:  dbVer.Manufacturer,
		Functions:     dbVer.Functions,
		Status:        dbVer.Status,
		Payload:       dbVer.Payload,
		CreatedAt:     dbVer.CreatedAt,
		UpdatedAt:     dbVer.UpdatedAt,
	}
	ver.GetImageUrl()
	return
}

func (n *Zigbee2mqttDevice) toDb(ver *m.Zigbee2mqttDevice) (dbVer *db.Zigbee2mqttDevice) {
	dbVer = &db.Zigbee2mqttDevice{
		Id:            ver.Id,
		Zigbee2mqttId: ver.Zigbee2mqttId,
		Name:          ver.Name,
		Type:          ver.Type,
		Model:         ver.Model,
		Description:   ver.Description,
		Manufacturer:  ver.Manufacturer,
		Functions:     ver.Functions,
		Status:        ver.Status,
		Payload:       ver.Payload,
		CreatedAt:     ver.CreatedAt,
		UpdatedAt:     ver.UpdatedAt,
	}
	return
}
