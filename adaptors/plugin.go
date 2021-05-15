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

type IPlugin interface {
	Add(plugin m.Plugin) error
	CreateOrUpdate(ver m.Plugin) error
	Update(plugin m.Plugin) error
	Delete(pluginId string) error
	List(limit, offset int64, orderBy, sort string) (list []m.Plugin, total int64, err error)
	Search(query string, limit, offset int) (list []m.Plugin, total int64, err error)
	GetByName(name string) (ver m.Plugin, err error)
	fromDb(dbVer db.Plugin) (plugin m.Plugin)
	toDb(plugin m.Plugin) (dbVer db.Plugin)
}

// Plugin ...
type Plugin struct {
	IPlugin
	table *db.Plugins
	db    *gorm.DB
}

// GetPluginAdaptor ...
func GetPluginAdaptor(d *gorm.DB) IPlugin {
	return &Plugin{
		table: &db.Plugins{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Plugin) Add(plugin m.Plugin) (err error) {
	err = n.table.Add(n.toDb(plugin))
	return
}

// CreateOrUpdate ...
func (n *Plugin) CreateOrUpdate(ver m.Plugin) (err error) {
	err = n.table.CreateOrUpdate(n.toDb(ver))
	return
}

// Update ...
func (n *Plugin) Update(plugin m.Plugin) (err error) {
	err = n.table.Update(n.toDb(plugin))
	return
}

// Delete ...
func (n *Plugin) Delete(name string) (err error) {
	err = n.table.Delete(name)
	return
}

// List ...
func (n *Plugin) List(limit, offset int64, orderBy, sort string) (list []m.Plugin, total int64, err error) {
	var dbList []db.Plugin
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]m.Plugin, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// Search ...
func (n *Plugin) Search(query string, limit, offset int) (list []m.Plugin, total int64, err error) {
	var dbList []db.Plugin
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]m.Plugin, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// GetByName ...
func (a *Plugin) GetByName(name string) (ver m.Plugin, err error) {

	var dbVer db.Plugin
	if dbVer, err = a.table.GetByName(name); err != nil {
		return
	}

	ver = a.fromDb(dbVer)

	return
}

func (n *Plugin) fromDb(dbVer db.Plugin) (plugin m.Plugin) {
	plugin = m.Plugin{
		Name:    dbVer.Name,
		Version: dbVer.Version,
		Enabled: dbVer.Enabled,
		System:  dbVer.System,
	}

	return
}

func (n *Plugin) toDb(ver m.Plugin) (dbVer db.Plugin) {
	dbVer = db.Plugin{
		Name:    ver.Name,
		Version: ver.Version,
		Enabled: ver.Enabled,
		System:  ver.System,
	}

	return
}
