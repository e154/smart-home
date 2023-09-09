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
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

// IArea ...
type IArea interface {
	Add(ctx context.Context, ver *m.Area) (id int64, err error)
	GetById(ctx context.Context, verId int64) (ver *m.Area, err error)
	GetByName(ctx context.Context, name string) (ver *m.Area, err error)
	Update(ctx context.Context, ver *m.Area) (err error)
	DeleteByName(ctx context.Context, name string) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.Area, total int64, err error)
	ListByPoint(ctx context.Context, point m.Point, limit, offset int64) (list []*m.Area, err error)
	Search(ctx context.Context, query string, limit, offset int64) (list []*m.Area, total int64, err error)
	GetDistance(ctx context.Context, point m.Point, areaId int64) (distance float64, err error)
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
func (a *Area) Add(ctx context.Context, ver *m.Area) (id int64, err error) {

	if id, err = a.table.Add(ctx, a.toDb(ver)); err != nil {
		return
	}

	return
}

// GetById ...
func (a *Area) GetById(ctx context.Context, verId int64) (ver *m.Area, err error) {

	var dbVer *db.Area
	if dbVer, err = a.table.GetById(ctx, verId); err != nil {
		return
	}

	ver = a.fromDb(dbVer)

	return
}

// Update ...
func (a *Area) Update(ctx context.Context, ver *m.Area) (err error) {
	err = a.table.Update(ctx, a.toDb(ver))
	return
}

// DeleteByName ...
func (a *Area) DeleteByName(ctx context.Context, name string) (err error) {
	err = a.table.DeleteByName(ctx, name)
	return
}

// List ...
func (a *Area) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.Area, total int64, err error) {
	var dbList []*db.Area
	if dbList, total, err = a.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Area, len(dbList))
	for i, dbVer := range dbList {
		list[i] = a.fromDb(dbVer)
	}
	return
}

// ListByPoint ...
func (a *Area) ListByPoint(ctx context.Context, point m.Point, limit, offset int64) (list []*m.Area, err error) {

	var dbList []*db.Area
	if dbList, err = a.table.ListByPoint(ctx, db.Point{Lon: point.Lon, Lat: point.Lat}, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Area, len(dbList))
	for i, dbVer := range dbList {
		list[i] = a.fromDb(dbVer)
	}
	return
}

// Search ...
func (a *Area) Search(ctx context.Context, query string, limit, offset int64) (list []*m.Area, total int64, err error) {
	var dbList []*db.Area
	if dbList, total, err = a.table.Search(ctx, query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Area, len(dbList))
	for i, dbVer := range dbList {
		list[i] = a.fromDb(dbVer)
	}

	return
}

// GetByName ...
func (a *Area) GetByName(ctx context.Context, name string) (ver *m.Area, err error) {

	var dbVer *db.Area
	if dbVer, err = a.table.GetByName(ctx, name); err != nil {
		return
	}

	ver = a.fromDb(dbVer)

	return
}

// GetDistance ...
func (a *Area) GetDistance(ctx context.Context, point m.Point, areaId int64) (distance float64, err error) {
	distance, err = a.table.GetDistance(ctx, db.Point{Lon: point.Lon, Lat: point.Lat}, areaId)
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
	if len(dbVer.Payload) > 0 {
		payload := &m.AreaPayload{}
		_ = json.Unmarshal(dbVer.Payload, payload)
		ver.Zoom = payload.Zoom
		ver.Center = payload.Center
		ver.Resolution = payload.Resolution
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
	if ver.Polygon != nil && len(ver.Polygon) > 2 {
		dbVer.Polygon = &db.Polygon{}
		for _, point := range ver.Polygon {
			dbVer.Polygon.Points = append(dbVer.Polygon.Points, db.Point{
				Lon: float64(point.Lon),
				Lat: float64(point.Lat),
			})
		}
	}
	b, _ := json.Marshal(m.AreaPayload{
		Zoom:       ver.Zoom,
		Center:     ver.Center,
		Resolution: ver.Resolution,
	})
	_ = dbVer.Payload.UnmarshalJSON(b)
	return
}
