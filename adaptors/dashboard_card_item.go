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
	"github.com/jinzhu/gorm"
)

type IDashboardCardItem interface {
	Add(ver *m.DashboardCardItem) (id int64, err error)
	GetById(mapId int64) (ver *m.DashboardCardItem, err error)
	Update(ver *m.DashboardCardItem) (err error)
	Delete(id int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.DashboardCardItem, total int64, err error)
	fromDb(dbVer *db.DashboardCardItem) (ver *m.DashboardCardItem)
	toDb(ver *m.DashboardCardItem) (dbVer *db.DashboardCardItem)
}

// DashboardCardItem ...
type DashboardCardItem struct {
	IDashboardCardItem
	table *db.DashboardCardItems
	db    *gorm.DB
}

// GetDashboardCardItemAdaptor ...
func GetDashboardCardItemAdaptor(d *gorm.DB) IDashboardCardItem {
	return &DashboardCardItem{
		table: &db.DashboardCardItems{Db: d},
		db:    d,
	}
}

// Add ...
func (n *DashboardCardItem) Add(ver *m.DashboardCardItem) (id int64, err error) {
	dbVer := n.toDb(ver)
	id, err = n.table.Add(dbVer)
	return
}

// GetById ...
func (n *DashboardCardItem) GetById(mapId int64) (ver *m.DashboardCardItem, err error) {

	var dbVer *db.DashboardCardItem
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *DashboardCardItem) Update(ver *m.DashboardCardItem) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

// Delete ...
func (n *DashboardCardItem) Delete(id int64) (err error) {
	err = n.table.Delete(id)
	return
}

// List ...
func (n *DashboardCardItem) List(limit, offset int64, orderBy, sort string) (list []*m.DashboardCardItem, total int64, err error) {

	if sort == "" {
		sort = "name"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.DashboardCardItem
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.DashboardCardItem, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *DashboardCardItem) fromDb(dbVer *db.DashboardCardItem) (ver *m.DashboardCardItem) {
	ver = &m.DashboardCardItem{
		Id:              dbVer.Id,
		Title:           dbVer.Title,
		Type:            dbVer.Type,
		Weight:          dbVer.Weight,
		Enabled:         dbVer.Enabled,
		Payload:         dbVer.Payload, //todo
		EntityId:        dbVer.EntityId,
		DashboardCardId: dbVer.DashboardCardId,
		CreatedAt:       dbVer.CreatedAt,
		UpdatedAt:       dbVer.UpdatedAt,
	}

	return
}

func (n *DashboardCardItem) toDb(ver *m.DashboardCardItem) (dbVer *db.DashboardCardItem) {
	dbVer = &db.DashboardCardItem{
		Id:              ver.Id,
		Title:           ver.Title,
		Type:            ver.Type,
		Weight:          ver.Weight,
		Enabled:         ver.Enabled,
		Payload:         ver.Payload, //todo
		EntityId:        ver.EntityId,
		DashboardCardId: ver.DashboardCardId,
	}

	return
}
