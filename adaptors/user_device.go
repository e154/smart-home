// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

// IUserDevice ...
type IUserDevice interface {
	Add(ctx context.Context, ver *m.UserDevice) (id int64, err error)
	GetByUserId(ctx context.Context, userId int64) (list []*m.UserDevice, err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.UserDevice, total int64, err error)
	Delete(ctx context.Context, id int64) (err error)
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
func (n *UserDevice) Add(ctx context.Context, ver *m.UserDevice) (id int64, err error) {

	if id, err = n.table.Add(ctx, n.toDb(ver)); err != nil {
		return
	}

	return
}

// GetByUserId ...
func (n *UserDevice) GetByUserId(ctx context.Context, userId int64) (list []*m.UserDevice, err error) {

	var dbList []*db.UserDevice
	if dbList, err = n.table.GetByUserId(ctx, userId); err != nil {
		return
	}

	list = make([]*m.UserDevice, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// List ...
func (n *UserDevice) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.UserDevice, total int64, err error) {

	if sort == "" {
		sort = "id"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.UserDevice
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.UserDevice, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

// Delete ...
func (n *UserDevice) Delete(ctx context.Context, id int64) (err error) {
	err = n.table.Delete(ctx, id)
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
