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

// MapText ...
type MapText struct {
	table *db.MapTexts
	db    *gorm.DB
}

// GetMapTextAdaptor ...
func GetMapTextAdaptor(d *gorm.DB) *MapText {
	return &MapText{
		table: &db.MapTexts{Db: d},
		db:    d,
	}
}

// Add ...
func (n *MapText) Add(ver *m.MapText) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

// GetById ...
func (n *MapText) GetById(mapId int64) (ver *m.MapText, err error) {

	var dbVer *db.MapText
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *MapText) Update(ver *m.MapText) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

// Sort ...
func (n *MapText) Sort(ver *m.MapText) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Sort(dbVer)
	return
}

// Delete ...
func (n *MapText) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

// List ...
func (n *MapText) List(limit, offset int64, orderBy, sort string) (list []*m.MapText, total int64, err error) {
	var dbList []*db.MapText
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.MapText, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *MapText) fromDb(dbVer *db.MapText) (ver *m.MapText) {
	ver = &m.MapText{
		Id:    dbVer.Id,
		Text:  dbVer.Text,
		Style: dbVer.Style,
	}

	return
}

func (n *MapText) toDb(ver *m.MapText) (dbVer *db.MapText) {
	dbVer = &db.MapText{
		Id:    ver.Id,
		Text:  ver.Text,
		Style: ver.Style,
	}
	return
}
