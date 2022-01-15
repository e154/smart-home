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

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

// IPlugin ...
type IPlugin interface {
	Add(plugin m.Plugin) error
	CreateOrUpdate(ver m.Plugin) error
	Update(plugin m.Plugin) error
	Delete(pluginId string) error
	List(limit, offset int64, orderBy, sort string) (list []m.Plugin, total int64, err error)
	Search(query string, limit, offset int64) (list []m.Plugin, total int64, err error)
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
func (p *Plugin) Add(plugin m.Plugin) (err error) {
	err = p.table.Add(p.toDb(plugin))
	return
}

// CreateOrUpdate ...
func (p *Plugin) CreateOrUpdate(plugin m.Plugin) (err error) {
	err = p.table.CreateOrUpdate(p.toDb(plugin))
	return
}

// Update ...
func (p *Plugin) Update(plugin m.Plugin) (err error) {
	err = p.table.Update(p.toDb(plugin))
	return
}

// Delete ...
func (p *Plugin) Delete(name string) (err error) {
	err = p.table.Delete(name)
	return
}

// List ...
func (p *Plugin) List(limit, offset int64, orderBy, sort string) (list []m.Plugin, total int64, err error) {
	var dbList []db.Plugin
	if dbList, total, err = p.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]m.Plugin, len(dbList))
	for i, dbVer := range dbList {
		list[i] = p.fromDb(dbVer)
	}
	return
}

// Search ...
func (p *Plugin) Search(query string, limit, offset int64) (list []m.Plugin, total int64, err error) {
	var dbList []db.Plugin
	if dbList, total, err = p.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]m.Plugin, len(dbList))
	for i, dbVer := range dbList {
		list[i] = p.fromDb(dbVer)
	}

	return
}

// GetByName ...
func (p *Plugin) GetByName(name string) (ver m.Plugin, err error) {

	var dbVer db.Plugin
	if dbVer, err = p.table.GetByName(name); err != nil {
		return
	}

	ver = p.fromDb(dbVer)

	return
}

func (p *Plugin) fromDb(dbVer db.Plugin) (ver m.Plugin) {
	ver = m.Plugin{
		Name:    dbVer.Name,
		Version: dbVer.Version,
		Enabled: dbVer.Enabled,
		System:  dbVer.System,
		Actor:   dbVer.Actor,
	}

	// deserialize settings
	b, _ := dbVer.Settings.MarshalJSON()
	settings := m.EntitySettings{}
	json.Unmarshal(b, &settings)
	ver.Settings = settings.Settings

	return
}

func (p *Plugin) toDb(ver m.Plugin) (dbVer db.Plugin) {
	dbVer = db.Plugin{
		Name:    ver.Name,
		Version: ver.Version,
		Enabled: ver.Enabled,
		System:  ver.System,
		Actor:   ver.Actor,
	}

	// serialize settings
	b, _ := json.Marshal(m.EntitySettings{
		Settings: ver.Settings,
	})
	dbVer.Settings.UnmarshalJSON(b)

	return
}
