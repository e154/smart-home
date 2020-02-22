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

type Zigbee2mqttDevice struct {
	table *db.Zigbee2mqttDevices
	db    *gorm.DB
}

func GetZigbee2mqttDeviceAdaptor(d *gorm.DB) *Zigbee2mqttDevice {
	return &Zigbee2mqttDevice{
		table: &db.Zigbee2mqttDevices{Db: d},
		db:    d,
	}
}

func (n *Zigbee2mqttDevice) Add(ver *m.Zigbee2mqttDevice) (err error) {

	err = n.table.Add(n.toDb(ver))

	return
}

func (n *Zigbee2mqttDevice) GetById(id string) (ver *m.Zigbee2mqttDevice, err error) {

	var dbVer *db.Zigbee2mqttDevice
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Zigbee2mqttDevice) Update(ver *m.Zigbee2mqttDevice) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *Zigbee2mqttDevice) Delete(id string) (err error) {
	err = n.table.Delete(id)
	return
}

func (n *Zigbee2mqttDevice) List(limit, offset int64) (list []*m.Zigbee2mqttDevice, total int64, err error) {
	var dbList []*db.Zigbee2mqttDevice
	if dbList, total, err = n.table.List(limit, offset); err != nil {
		return
	}

	list = make([]*m.Zigbee2mqttDevice, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Zigbee2mqttDevice) fromDb(dbVer *db.Zigbee2mqttDevice) (ver *m.Zigbee2mqttDevice) {
	ver = &m.Zigbee2mqttDevice{
		Id:           dbVer.Id,
		Name:         dbVer.Name,
		Type:         dbVer.Type,
		Model:        dbVer.Model,
		Description:  dbVer.Description,
		Manufacturer: dbVer.Manufacturer,
		Functions:    dbVer.Functions,
		Status:       dbVer.Status,
		CreatedAt:    dbVer.CreatedAt,
		UpdatedAt:    dbVer.UpdatedAt,
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
		CreatedAt:     ver.CreatedAt,
		UpdatedAt:     ver.UpdatedAt,
	}
	return
}
