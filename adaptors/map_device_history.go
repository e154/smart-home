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

type MapDeviceHistory struct {
	table *db.MapDeviceHistories
	db    *gorm.DB
}

func GetMapDeviceHistoryAdaptor(d *gorm.DB) *MapDeviceHistory {
	return &MapDeviceHistory{
		table: &db.MapDeviceHistories{Db: d},
		db:    d,
	}
}

func (n *MapDeviceHistory) Add(ver m.MapDeviceHistory) (id int64, err error) {

	id, err = n.table.Add(n.toDb(ver))

	return
}

func (n *MapDeviceHistory) GetAllByDeviceId(id int64, limit, offset int) (list []*m.MapDeviceHistory, total int64, err error) {

	var dbList []*db.MapDeviceHistory
	if dbList, total, err = n.table.GetAllByDeviceId(id, limit, offset); err != nil {
		return
	}

	list = make([]*m.MapDeviceHistory, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

func (n *MapDeviceHistory) fromDb(dbVer *db.MapDeviceHistory) (ver *m.MapDeviceHistory) {
	ver = &m.MapDeviceHistory{
		Id:          dbVer.Id,
		MapDeviceId: dbVer.MapDeviceId,
		Type:        dbVer.Type,
		Description: dbVer.Description,
		CreatedAt:   dbVer.CreatedAt,
	}

	return
}

func (n *MapDeviceHistory) toDb(ver m.MapDeviceHistory) (dbVer db.MapDeviceHistory) {
	dbVer = db.MapDeviceHistory{
		Id:          ver.Id,
		MapDeviceId: ver.MapDeviceId,
		Type:        ver.Type,
		Description: ver.Description,
		CreatedAt:   ver.CreatedAt,
	}

	return
}
