// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.DashboardRepo = (*Dashboard)(nil)

// Dashboard ...
type Dashboard struct {
	table *db.Dashboards
	db    *gorm.DB
}

// GetDashboardAdaptor ...
func GetDashboardAdaptor(d *gorm.DB) *Dashboard {
	return &Dashboard{
		table: &db.Dashboards{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *Dashboard) Add(ctx context.Context, ver *models.Dashboard) (id int64, err error) {
	dbVer := n.toDb(ver)
	id, err = n.table.Add(ctx, dbVer)
	return
}

// GetById ...
func (n *Dashboard) GetById(ctx context.Context, mapId int64) (ver *models.Dashboard, err error) {

	var dbVer *db.Dashboard
	if dbVer, err = n.table.GetById(ctx, mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Search ...
func (n *Dashboard) Search(ctx context.Context, query string, limit, offset int64) (list []*models.Dashboard, total int64, err error) {
	var dbList []*db.Dashboard
	if dbList, total, err = n.table.Search(ctx, query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*models.Dashboard, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

// Update ...
func (n *Dashboard) Update(ctx context.Context, ver *models.Dashboard) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(ctx, dbVer)
	return
}

// Delete ...
func (n *Dashboard) Delete(ctx context.Context, id int64) (err error) {
	err = n.table.Delete(ctx, id)
	return
}

// List ...
func (n *Dashboard) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*models.Dashboard, total int64, err error) {

	if sort == "" {
		sort = "name"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.Dashboard
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*models.Dashboard, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Dashboard) fromDb(dbVer *db.Dashboard) (ver *models.Dashboard) {
	ver = &models.Dashboard{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Enabled:     dbVer.Enabled,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
		AreaId:      dbVer.AreaId,
	}

	// Area
	if dbVer.Area != nil {
		areaAdaptor := GetAreaAdaptor(n.db)
		ver.Area = areaAdaptor.fromDb(dbVer.Area)
		ver.AreaId = common.Int64(dbVer.Area.Id)
	}

	// tabs
	tabsAdaptor := GetDashboardTabAdaptor(n.db)
	ver.Tabs = make([]*models.DashboardTab, 0, len(ver.Tabs))
	for _, tab := range dbVer.Tabs {
		ver.Tabs = append(ver.Tabs, tabsAdaptor.fromDb(tab))
	}

	return
}

func (n *Dashboard) toDb(ver *models.Dashboard) (dbVer *db.Dashboard) {
	dbVer = &db.Dashboard{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Enabled:     ver.Enabled,
		AreaId:      ver.AreaId,
	}

	// area
	if ver.Area != nil && ver.Area.Id != 0 {
		dbVer.AreaId = &ver.Area.Id
	}

	return
}
