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

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

type IDashboardTab interface {
	Add(ctx context.Context, ver *m.DashboardTab) (id int64, err error)
	GetById(ctx context.Context, mapId int64) (ver *m.DashboardTab, err error)
	Update(ctx context.Context, ver *m.DashboardTab) (err error)
	Delete(ctx context.Context, id int64) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.DashboardTab, total int64, err error)
	fromDb(dbVer *db.DashboardTab) (ver *m.DashboardTab)
	toDb(ver *m.DashboardTab) (dbVer *db.DashboardTab)
}

// DashboardTab ...
type DashboardTab struct {
	IDashboardTab
	table *db.DashboardTabs
	db    *gorm.DB
}

// GetDashboardTabAdaptor ...
func GetDashboardTabAdaptor(d *gorm.DB) IDashboardTab {
	return &DashboardTab{
		table: &db.DashboardTabs{Db: d},
		db:    d,
	}
}

// Add ...
func (n *DashboardTab) Add(ctx context.Context, ver *m.DashboardTab) (id int64, err error) {
	dbVer := n.toDb(ver)
	id, err = n.table.Add(ctx, dbVer)
	return
}

// GetById ...
func (n *DashboardTab) GetById(ctx context.Context, mapId int64) (ver *m.DashboardTab, err error) {

	var dbVer *db.DashboardTab
	if dbVer, err = n.table.GetById(ctx, mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *DashboardTab) Update(ctx context.Context, ver *m.DashboardTab) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(ctx, dbVer)
	return
}

// Delete ...
func (n *DashboardTab) Delete(ctx context.Context, id int64) (err error) {
	err = n.table.Delete(ctx, id)
	return
}

// List ...
func (n *DashboardTab) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.DashboardTab, total int64, err error) {

	if sort == "" {
		sort = "name"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.DashboardTab
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.DashboardTab, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *DashboardTab) fromDb(dbVer *db.DashboardTab) (ver *m.DashboardTab) {
	ver = &m.DashboardTab{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		ColumnWidth: dbVer.ColumnWidth,
		Gap:         dbVer.Gap,
		Background:  dbVer.Background,
		Icon:        dbVer.Icon,
		Enabled:     dbVer.Enabled,
		Weight:      dbVer.Weight,
		DashboardId: dbVer.DashboardId,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}

	// Cards
	cardsAdaptor := GetDashboardCardAdaptor(n.db)
	ver.Cards = make([]*m.DashboardCard, 0, len(ver.Cards))
	for _, tab := range dbVer.Cards {
		ver.Cards = append(ver.Cards, cardsAdaptor.fromDb(tab))
	}

	return
}

func (n *DashboardTab) toDb(ver *m.DashboardTab) (dbVer *db.DashboardTab) {
	dbVer = &db.DashboardTab{
		Id:          ver.Id,
		Name:        ver.Name,
		Icon:        ver.Icon,
		ColumnWidth: ver.ColumnWidth,
		Gap:         ver.Gap,
		Background:  ver.Background,
		Enabled:     ver.Enabled,
		Weight:      ver.Weight,
		DashboardId: ver.DashboardId,
	}

	return
}
