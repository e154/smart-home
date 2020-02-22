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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type Zigbee2mqtt struct {
	table *db.Zigbee2mqtts
	db    *gorm.DB
}

func GetZigbee2mqttAdaptor(d *gorm.DB) *Zigbee2mqtt {
	return &Zigbee2mqtt{
		table: &db.Zigbee2mqtts{Db: d},
		db:    d,
	}
}

func (n *Zigbee2mqtt) Add(ver *m.Zigbee2mqtt) (id int64, err error) {

	id, err = n.table.Add(n.toDb(ver))

	return
}

func (n *Zigbee2mqtt) GetById(id int64) (ver *m.Zigbee2mqtt, err error) {

	var dbVer *db.Zigbee2mqtt
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Zigbee2mqtt) Update(ver *m.Zigbee2mqtt) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *Zigbee2mqtt) Delete(id int64) (err error) {
	err = n.table.Delete(id)
	return
}

func (n *Zigbee2mqtt) List(limit, offset int64) (list []*m.Zigbee2mqtt, total int64, err error) {
	var dbList []*db.Zigbee2mqtt
	if dbList, total, err = n.table.List(limit, offset); err != nil {
		return
	}

	list = make([]*m.Zigbee2mqtt, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Zigbee2mqtt) fromDb(dbVer *db.Zigbee2mqtt) (ver *m.Zigbee2mqtt) {
	ver = &m.Zigbee2mqtt{
		Id:        dbVer.Id,
		Login:     dbVer.Login,
		Name:      dbVer.Name,
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
	}

	zigbee2mqttDeviceAdaptor := GetZigbee2mqttDeviceAdaptor(n.db)
	for _, dbDev := range dbVer.Devices {
		dev := zigbee2mqttDeviceAdaptor.fromDb(dbDev)
		ver.Devices = append(ver.Devices, dev)
	}

	return
}

func (n *Zigbee2mqtt) toDb(ver *m.Zigbee2mqtt) (dbVer *db.Zigbee2mqtt) {
	dbVer = &db.Zigbee2mqtt{
		Id:        ver.Id,
		Login:     ver.Login,
		Name:      ver.Name,
		CreatedAt: ver.CreatedAt,
		UpdatedAt: ver.UpdatedAt,
	}
	if ver.Password != "" {
		dbVer.EncryptedPassword, _ = common.HashPassword(ver.Password)
	}
	return
}
