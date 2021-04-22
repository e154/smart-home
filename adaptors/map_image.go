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

type IMapImage interface {
	Add(ver *m.MapImage) (id int64, err error)
	GetById(mapId int64) (ver *m.MapImage, err error)
	Update(ver *m.MapImage) (err error)
	Sort(ver *m.MapImage) (err error)
	Delete(mapId int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.MapImage, total int64, err error)
	fromDb(dbVer *db.MapImage) (ver *m.MapImage)
	toDb(ver *m.MapImage) (dbVer *db.MapImage)
}

// MapImage ...
type MapImage struct {
	IMapImage
	table *db.MapImages
	db    *gorm.DB
}

// GetMapImageAdaptor ...
func GetMapImageAdaptor(d *gorm.DB) IMapImage {
	return &MapImage{
		table: &db.MapImages{Db: d},
		db:    d,
	}
}

// Add ...
func (n *MapImage) Add(ver *m.MapImage) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

// GetById ...
func (n *MapImage) GetById(mapId int64) (ver *m.MapImage, err error) {

	var dbVer *db.MapImage
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *MapImage) Update(ver *m.MapImage) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

// Sort ...
func (n *MapImage) Sort(ver *m.MapImage) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Sort(dbVer)
	return
}

// Delete ...
func (n *MapImage) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

// List ...
func (n *MapImage) List(limit, offset int64, orderBy, sort string) (list []*m.MapImage, total int64, err error) {
	var dbList []*db.MapImage
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.MapImage, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *MapImage) fromDb(dbVer *db.MapImage) (ver *m.MapImage) {
	ver = &m.MapImage{
		Id:      dbVer.Id,
		ImageId: dbVer.ImageId,
		Style:   dbVer.Style,
	}

	if dbVer.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		ver.Image = imageAdaptor.fromDb(dbVer.Image)
	}

	return
}

func (n *MapImage) toDb(ver *m.MapImage) (dbVer *db.MapImage) {
	dbVer = &db.MapImage{
		Id:      ver.Id,
		ImageId: ver.ImageId,
		Style:   ver.Style,
	}
	return
}
