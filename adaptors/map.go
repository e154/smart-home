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
	"sort"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

// IMap ...
type IMap interface {
	Add(ver *m.Map) (id int64, err error)
	GetById(mapId int64) (ver *m.Map, err error)
	GetFullById(mapId int64) (ver *m.Map, err error)
	Update(ver *m.Map) (err error)
	Delete(mapId int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.Map, total int64, err error)
	Search(query string, limit, offset int) (list []*m.Map, total int64, err error)
	fromDb(dbVer *db.Map) (ver *m.Map)
	toDb(ver *m.Map) (dbVer *db.Map)
}

// Map ...
type Map struct {
	IMap
	table *db.Maps
	db    *gorm.DB
}

// GetMapAdaptor ...
func GetMapAdaptor(d *gorm.DB) IMap {
	return &Map{
		table: &db.Maps{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Map) Add(ver *m.Map) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

// GetById ...
func (n *Map) GetById(mapId int64) (ver *m.Map, err error) {

	var dbVer *db.Map
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// GetFullById ...
func (n *Map) GetFullById(mapId int64) (ver *m.Map, err error) {

	var dbVer *db.Map
	if dbVer, err = n.table.GetFullById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	sort.Sort(m.SortMapLayersByWeight(ver.Layers))

	return
}

// Update ...
func (n *Map) Update(ver *m.Map) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

// Delete ...
func (n *Map) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

// List ...
func (n *Map) List(limit, offset int64, orderBy, sort string) (list []*m.Map, total int64, err error) {
	var dbList []*db.Map
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Map, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// Search ...
func (n *Map) Search(query string, limit, offset int) (list []*m.Map, total int64, err error) {
	var dbList []*db.Map
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Map, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Map) fromDb(dbVer *db.Map) (ver *m.Map) {
	ver = &m.Map{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}
	options, _ := dbVer.Options.MarshalJSON()
	json.Unmarshal(options, &ver.Options)

	// layers
	layerAdaptor := GetMapLayerAdaptor(n.db)
	for _, dbLayer := range dbVer.Layers {
		layer := layerAdaptor.fromDb(dbLayer)
		ver.Layers = append(ver.Layers, layer)
	}

	return
}

func (n *Map) toDb(ver *m.Map) (dbVer *db.Map) {
	dbVer = &db.Map{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
	}
	options, _ := json.Marshal(ver.Options)
	dbVer.Options.UnmarshalJSON(options)
	return
}
