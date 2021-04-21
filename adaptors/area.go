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

// Area ...
type Area struct {
	table *db.Areas
	db    *gorm.DB
}

// GetAreaAdaptor ...
func GetAreaAdaptor(d *gorm.DB) *Area {
	return &Area{
		table: &db.Areas{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Area) Add(ver *m.Area) (id int64, err error) {

	if id, err = n.table.Add(n.toDb(ver)); err != nil {
		return
	}

	return
}

// GetById ...
func (n *Area) GetById(verId int64) (ver *m.Area, err error) {

	var dbVer *db.Area
	if dbVer, err = n.table.GetById(verId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *Area) Update(ver *m.Area) (err error) {
	err = n.table.Update(n.toDb(ver))
	return
}

// DeleteByName ...
func (n *Area) DeleteByName(name string) (err error) {
	err = n.table.DeleteByName(name)
	return
}

// List ...
func (n *Area) List(limit, offset int64, orderBy, sort string) (list []*m.Area, total int64, err error) {
	var dbList []*db.Area
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Area, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// Search ...
func (n *Area) Search(query string, limit, offset int) (list []*m.Area, total int64, err error) {
	var dbList []*db.Area
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Area, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// GetByName ...
func (a *Area) GetByName(name string) (ver *m.Area, err error) {

	var dbVer *db.Area
	if dbVer, err = a.table.GetByName(name); err != nil {
		return
	}

	ver = a.fromDb(dbVer)

	return
}

func (n *Area) fromDb(dbVer *db.Area) (ver *m.Area) {
	ver = &m.Area{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
	}

	return
}

func (n *Area) toDb(ver *m.Area) (dbVer *db.Area) {
	dbVer = &db.Area{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
	}

	return
}
