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

type MapZone struct {
	table *db.MapZones
	db    *gorm.DB
}

func GetMapZoneAdaptor(d *gorm.DB) *MapZone {
	return &MapZone{
		table: &db.MapZones{Db: d},
		db:    d,
	}
}

func (n *MapZone) Add(tag *m.MapZone) (id int64, err error) {

	dbTag := n.toDb(tag)
	id, err = n.table.Add(dbTag)

	return
}


func (n *MapZone) GetByName(zoneName string) (ver *m.MapZone, err error) {

	var dbVer *db.MapZone
	if dbVer, err = n.table.GetByName(zoneName); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *MapZone) Delete(name string) (err error) {

	err = n.table.Delete(name)

	return
}

func (n *MapZone) Search(query string, limit, offset int) (list []*m.MapZone, total int64, err error) {
	var dbList []*db.MapZone
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.MapZone, 0)
	for _, dbTag := range dbList {
		node := n.fromDb(dbTag)
		list = append(list, node)
	}

	return
}

func (n *MapZone) Clean() (err error) {

	err = n.table.Clean()

	return
}

func (n *MapZone) toDb(tag *m.MapZone) *db.MapZone {
	return &db.MapZone{
		Id:   tag.Id,
		Name: tag.Name,
	}
}

func (n *MapZone) fromDb(tag *db.MapZone) *m.MapZone {
	return &m.MapZone{
		Id:   tag.Id,
		Name: tag.Name,
	}
}
