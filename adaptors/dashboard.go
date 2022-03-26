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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type IDashboard interface {
	Add(ver *m.Dashboard) (id int64, err error)
	GetById(mapId int64) (ver *m.Dashboard, err error)
	Update(ver *m.Dashboard) (err error)
	Delete(id int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.Dashboard, total int64, err error)
}

// Dashboard ...
type Dashboard struct {
	IDashboard
	table *db.Dashboards
	db    *gorm.DB
}

// GetDashboardAdaptor ...
func GetDashboardAdaptor(d *gorm.DB) IDashboard {
	return &Dashboard{
		table: &db.Dashboards{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Dashboard) Add(ver *m.Dashboard) (id int64, err error) {
	dbVer := n.toDb(ver)
	id, err = n.table.Add(dbVer)
	return
}

// GetById ...
func (n *Dashboard) GetById(mapId int64) (ver *m.Dashboard, err error) {

	var dbVer *db.Dashboard
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *Dashboard) Update(ver *m.Dashboard) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

// Delete ...
func (n *Dashboard) Delete(id int64) (err error) {
	err = n.table.Delete(id)
	return
}

// List ...
func (n *Dashboard) List(limit, offset int64, orderBy, sort string) (list []*m.Dashboard, total int64, err error) {

	if sort == "" {
		sort = "name"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.Dashboard
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Dashboard, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Dashboard) fromDb(dbVer *db.Dashboard) (ver *m.Dashboard) {
	ver = &m.Dashboard{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Enabled:     dbVer.Enabled,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}

	// Area
	if dbVer.Area != nil {
		areaAdaptor := GetAreaAdaptor(n.db)
		ver.Area = areaAdaptor.fromDb(dbVer.Area)
		ver.AreaId = common.Int64(dbVer.Area.Id)
	}

	// tabs
	tabsAdaptor := GetDashboardTabAdaptor(n.db)
	ver.Tabs = make([]*m.DashboardTab, 0, len(ver.Tabs))
	for _, tab := range dbVer.Tabs {
		ver.Tabs = append(ver.Tabs, tabsAdaptor.fromDb(tab))
	}

	return
}

func (n *Dashboard) toDb(ver *m.Dashboard) (dbVer *db.Dashboard) {
	dbVer = &db.Dashboard{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Enabled:     ver.Enabled,
	}

	// area
	if ver.Area != nil && ver.Area.Id != 0 {
		dbVer.AreaId = &ver.Area.Id
	}

	return
}
