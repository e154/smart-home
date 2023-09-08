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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// IZigbee2mqtt ...
type IZigbee2mqtt interface {
	Add(ver *m.Zigbee2mqtt) (id int64, err error)
	GetById(id int64) (ver *m.Zigbee2mqtt, err error)
	Update(ver *m.Zigbee2mqtt) (err error)
	Delete(id int64) (err error)
	List(limit, offset int64) (list []*m.Zigbee2mqtt, total int64, err error)
	GetByLogin(login string) (ver *m.Zigbee2mqtt, err error)
	fromDb(dbVer *db.Zigbee2mqtt) (ver *m.Zigbee2mqtt)
	toDb(ver *m.Zigbee2mqtt) (dbVer *db.Zigbee2mqtt)
}

// Zigbee2mqtt ...
type Zigbee2mqtt struct {
	IZigbee2mqtt
	table *db.Zigbee2mqtts
	db    *gorm.DB
}

// GetZigbee2mqttAdaptor ...
func GetZigbee2mqttAdaptor(d *gorm.DB) IZigbee2mqtt {
	return &Zigbee2mqtt{
		table: &db.Zigbee2mqtts{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Zigbee2mqtt) Add(ver *m.Zigbee2mqtt) (id int64, err error) {

	id, err = n.table.Add(n.toDb(ver))

	return
}

// GetById ...
func (n *Zigbee2mqtt) GetById(id int64) (ver *m.Zigbee2mqtt, err error) {

	var dbVer *db.Zigbee2mqtt
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *Zigbee2mqtt) Update(ver *m.Zigbee2mqtt) (err error) {
	err = n.table.Update(n.toDb(ver))
	return
}

// Delete ...
func (n *Zigbee2mqtt) Delete(id int64) (err error) {
	err = n.table.Delete(id)
	return
}

// List ...
func (n *Zigbee2mqtt) List(limit, offset int64) (list []*m.Zigbee2mqtt, total int64, err error) {
	var dbList []*db.Zigbee2mqtt
	if dbList, total, err = n.table.List(int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Zigbee2mqtt, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// GetByLogin ...
func (a *Zigbee2mqtt) GetByLogin(login string) (ver *m.Zigbee2mqtt, err error) {

	var dbVer *db.Zigbee2mqtt
	if dbVer, err = a.table.GetByLogin(login); err != nil {
		return
	}

	ver = a.fromDb(dbVer)

	return
}

func (n *Zigbee2mqtt) fromDb(dbVer *db.Zigbee2mqtt) (ver *m.Zigbee2mqtt) {
	ver = &m.Zigbee2mqtt{
		Id:                dbVer.Id,
		Login:             dbVer.Login,
		Name:              dbVer.Name,
		PermitJoin:        dbVer.PermitJoin,
		BaseTopic:         dbVer.BaseTopic,
		CreatedAt:         dbVer.CreatedAt,
		UpdatedAt:         dbVer.UpdatedAt,
		EncryptedPassword: dbVer.EncryptedPassword,
	}

	if len(dbVer.Devices) > 0 {
		zigbee2mqttDeviceAdaptor := GetZigbee2mqttDeviceAdaptor(n.db)
		for _, dbDev := range dbVer.Devices {
			dev := zigbee2mqttDeviceAdaptor.fromDb(dbDev)
			ver.Devices = append(ver.Devices, dev)
		}
	} else {
		ver.Devices = make([]*m.Zigbee2mqttDevice, 0)
	}

	return
}

func (n *Zigbee2mqtt) toDb(ver *m.Zigbee2mqtt) (dbVer *db.Zigbee2mqtt) {
	dbVer = &db.Zigbee2mqtt{
		Id:                ver.Id,
		Login:             ver.Login,
		Name:              ver.Name,
		PermitJoin:        ver.PermitJoin,
		BaseTopic:         ver.BaseTopic,
		EncryptedPassword: ver.EncryptedPassword,
	}
	if ver.Password != nil {
		if *ver.Password == "" {
			dbVer.EncryptedPassword = ""
		} else {
			dbVer.EncryptedPassword, _ = common.HashPassword(*ver.Password)
		}
	}
	return
}
