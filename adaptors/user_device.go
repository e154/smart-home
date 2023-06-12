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
	"encoding/json"

	"gorm.io/gorm"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

// IUserDevice ...
type IUserDevice interface {
	Add(ver *m.UserDevice) (id int64, err error)
	GetByUserId(userId int64) (list []*m.UserDevice, err error)
	Delete(id int64) (err error)
	fromDb(dbVer *db.UserDevice) (ver *m.UserDevice)
	toDb(ver *m.UserDevice) (dbVer *db.UserDevice)
}

// UserDevice ...
type UserDevice struct {
	IUserDevice
	table *db.UserDevices
	db    *gorm.DB
}

// GetUserDeviceAdaptor ...
func GetUserDeviceAdaptor(d *gorm.DB) IUserDevice {
	return &UserDevice{
		table: &db.UserDevices{Db: d},
		db:    d,
	}
}

// Add ...
func (n *UserDevice) Add(ver *m.UserDevice) (id int64, err error) {

	if id, err = n.table.Add(n.toDb(ver)); err != nil {
		return
	}

	return
}

// GetByUserId ...
func (n *UserDevice) GetByUserId(userId int64) (list []*m.UserDevice, err error) {

	var dbList []*db.UserDevice
	if dbList, err = n.table.GetByUserId(userId); err != nil {
		return
	}

	list = make([]*m.UserDevice, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// Delete ...
func (n *UserDevice) Delete(id int64) (err error) {
	err = n.table.Delete(id)
	return
}

func (n *UserDevice) fromDb(dbVer *db.UserDevice) (ver *m.UserDevice) {
	ver = &m.UserDevice{
		Id:           dbVer.Id,
		UserId:       dbVer.UserId,
		Subscription: &m.Subscription{},
		CreatedAt:    dbVer.CreatedAt,
	}

	// deserialize Subscription
	b, _ := dbVer.PushRegistration.MarshalJSON()
	_ = json.Unmarshal(b, &ver.Subscription)

	return
}

func (n *UserDevice) toDb(ver *m.UserDevice) (dbVer *db.UserDevice) {
	dbVer = &db.UserDevice{
		Id:        ver.Id,
		UserId:    ver.UserId,
		CreatedAt: ver.CreatedAt,
	}

	// serialize Subscription
	b, _ := json.Marshal(ver.Subscription)
	_ = dbVer.PushRegistration.UnmarshalJSON(b)

	return
}
