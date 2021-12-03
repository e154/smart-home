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
	"sort"
)

// IMapLayer ...
type IMapLayer interface {
	Add(ver *m.MapLayer) (id int64, err error)
	GetById(mapId int64) (ver *m.MapLayer, err error)
	Update(ver *m.MapLayer) (err error)
	Sort(ver *m.MapLayer) (err error)
	Delete(mapId int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.MapLayer, total int64, err error)
	fromDb(dbVer *db.MapLayer) (ver *m.MapLayer)
	toDb(ver *m.MapLayer) (dbVer *db.MapLayer)
}

// MapLayer ...
type MapLayer struct {
	IMapLayer
	table *db.MapLayers
	db    *gorm.DB
}

// GetMapLayerAdaptor ...
func GetMapLayerAdaptor(d *gorm.DB) IMapLayer {
	return &MapLayer{
		table: &db.MapLayers{Db: d},
		db:    d,
	}
}

// Add ...
func (n *MapLayer) Add(ver *m.MapLayer) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

// GetById ...
func (n *MapLayer) GetById(mapId int64) (ver *m.MapLayer, err error) {

	var dbVer *db.MapLayer
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *MapLayer) Update(ver *m.MapLayer) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

// Sort ...
func (n *MapLayer) Sort(ver *m.MapLayer) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Sort(dbVer)
	return
}

// Delete ...
func (n *MapLayer) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

// List ...
func (n *MapLayer) List(limit, offset int64, orderBy, sort string) (list []*m.MapLayer, total int64, err error) {
	var dbList []*db.MapLayer
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.MapLayer, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *MapLayer) fromDb(dbVer *db.MapLayer) (ver *m.MapLayer) {
	ver = &m.MapLayer{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		MapId:       dbVer.MapId,
		Status:      dbVer.Status,
		Weight:      dbVer.Weight,
		Description: dbVer.Description,
		Elements:    make([]*m.MapElement, 0),
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}

	// elements
	mapElementAdaptor := GetMapElementAdaptor(n.db)
	for _, dbElement := range dbVer.Elements {
		element := mapElementAdaptor.fromDb(dbElement)
		ver.Elements = append(ver.Elements, element)
	}

	// map
	if dbVer.Map != nil {
		mapAdaptor := GetMapAdaptor(n.db)
		ver.Map = mapAdaptor.fromDb(dbVer.Map)
	}

	sort.Sort(m.SortMapElementByWeight(ver.Elements))

	return
}

func (n *MapLayer) toDb(ver *m.MapLayer) (dbVer *db.MapLayer) {
	dbVer = &db.MapLayer{
		Id:          ver.Id,
		Name:        ver.Name,
		MapId:       ver.MapId,
		Status:      ver.Status,
		Weight:      ver.Weight,
		Description: ver.Description,
	}
	return
}
