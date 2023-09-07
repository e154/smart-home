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
	"gorm.io/gorm"
)

// IArea ...
type IArea interface {
	Add(ver *m.Area) (id int64, err error)
	GetById(verId int64) (ver *m.Area, err error)
	GetByName(name string) (ver *m.Area, err error)
	Update(ver *m.Area) (err error)
	DeleteByName(name string) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.Area, total int64, err error)
	Search(query string, limit, offset int64) (list []*m.Area, total int64, err error)
	fromDb(dbVer *db.Area) (ver *m.Area)
	toDb(ver *m.Area) (dbVer *db.Area)
}

// Area ...
type Area struct {
	IArea
	table *db.Areas
	db    *gorm.DB
}

// GetAreaAdaptor ...
func GetAreaAdaptor(d *gorm.DB) IArea {
	return &Area{
		table: &db.Areas{Db: d},
		db:    d,
	}
}

// Add ...
func (a *Area) Add(ver *m.Area) (id int64, err error) {

	if id, err = a.table.Add(a.toDb(ver)); err != nil {
		return
	}

	return
}

// GetById ...
func (a *Area) GetById(verId int64) (ver *m.Area, err error) {

	var dbVer *db.Area
	if dbVer, err = a.table.GetById(verId); err != nil {
		return
	}

	ver = a.fromDb(dbVer)

	return
}

// Update ...
func (a *Area) Update(ver *m.Area) (err error) {
	err = a.table.Update(a.toDb(ver))
	return
}

// DeleteByName ...
func (a *Area) DeleteByName(name string) (err error) {
	err = a.table.DeleteByName(name)
	return
}

// List ...
func (a *Area) List(limit, offset int64, orderBy, sort string) (list []*m.Area, total int64, err error) {
	var dbList []*db.Area
	if dbList, total, err = a.table.List(int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Area, len(dbList))
	for i, dbVer := range dbList {
		list[i] = a.fromDb(dbVer)
	}
	return
}

// Search ...
func (a *Area) Search(query string, limit, offset int64) (list []*m.Area, total int64, err error) {
	var dbList []*db.Area
	if dbList, total, err = a.table.Search(query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Area, len(dbList))
	for i, dbVer := range dbList {
		list[i] = a.fromDb(dbVer)
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

func (a *Area) fromDb(dbVer *db.Area) (ver *m.Area) {
	ver = &m.Area{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}
	if dbVer.Polygon != nil {
		for _, point := range dbVer.Polygon.Points {
			ver.Polygon = append(ver.Polygon, m.Point{
				Lon: point.Lon,
				Lat: point.Lat,
			})
		}
	}
	return
}

func (a *Area) toDb(ver *m.Area) (dbVer *db.Area) {
	dbVer = &db.Area{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}
	if ver.Polygon != nil {
		dbVer.Polygon = &db.Polygon{}
		for _, point := range ver.Polygon {
			dbVer.Polygon.Points = append(dbVer.Polygon.Points, db.Point{
				Lon: point.Lon,
				Lat: point.Lat,
			})
		}
	}
	return
}
