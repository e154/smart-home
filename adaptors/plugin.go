// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// IPlugin ...
type IPlugin interface {
	Add(ctx context.Context, plugin *m.Plugin) error
	CreateOrUpdate(ctx context.Context, ver *m.Plugin) error
	Update(ctx context.Context, plugin *m.Plugin) error
	Delete(ctx context.Context, pluginId string) error
	List(ctx context.Context, limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*m.Plugin, total int64, err error)
	Search(ctx context.Context, query string, limit, offset int64) (list []*m.Plugin, total int64, err error)
	GetByName(ctx context.Context, name string) (ver *m.Plugin, err error)
	fromDb(dbVer *db.Plugin) (plugin *m.Plugin)
	toDb(plugin *m.Plugin) (dbVer *db.Plugin)
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
func (p *Plugin) Add(ctx context.Context, plugin *m.Plugin) (err error) {
	err = p.table.Add(ctx, p.toDb(plugin))
	return
}

// CreateOrUpdate ...
func (p *Plugin) CreateOrUpdate(ctx context.Context, plugin *m.Plugin) (err error) {
	err = p.table.CreateOrUpdate(ctx, p.toDb(plugin))
	return
}

// Update ...
func (p *Plugin) Update(ctx context.Context, plugin *m.Plugin) (err error) {
	err = p.table.Update(ctx, p.toDb(plugin))
	return
}

// Delete ...
func (p *Plugin) Delete(ctx context.Context, name string) (err error) {
	err = p.table.Delete(ctx, name)
	return
}

// List ...
func (p *Plugin) List(ctx context.Context, limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*m.Plugin, total int64, err error) {
	var dbList []*db.Plugin
	if dbList, total, err = p.table.List(ctx, int(limit), int(offset), orderBy, sort, onlyEnabled); err != nil {
		return
	}

	list = make([]*m.Plugin, len(dbList))
	for i, dbVer := range dbList {
		list[i] = p.fromDb(dbVer)
	}
	return
}

// Search ...
func (p *Plugin) Search(ctx context.Context, query string, limit, offset int64) (list []*m.Plugin, total int64, err error) {
	var dbList []*db.Plugin
	if dbList, total, err = p.table.Search(ctx, query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Plugin, len(dbList))
	for i, dbVer := range dbList {
		list[i] = p.fromDb(dbVer)
	}

	return
}

// GetByName ...
func (p *Plugin) GetByName(ctx context.Context, name string) (ver *m.Plugin, err error) {

	var dbVer *db.Plugin
	if dbVer, err = p.table.GetByName(ctx, name); err != nil {
		return
	}

	ver = p.fromDb(dbVer)

	return
}

func (p *Plugin) fromDb(dbVer *db.Plugin) (ver *m.Plugin) {
	ver = &m.Plugin{
		Name:    dbVer.Name,
		Version: dbVer.Version,
		Enabled: dbVer.Enabled,
		System:  dbVer.System,
		Actor:   dbVer.Actor,
	}

	// deserialize settings
	b, _ := dbVer.Settings.MarshalJSON()
	settings := make(m.AttributeValue, 0)
	_ = json.Unmarshal(b, &settings)
	ver.Settings = settings

	return
}

func (p *Plugin) toDb(ver *m.Plugin) (dbVer *db.Plugin) {
	dbVer = &db.Plugin{
		Name:    ver.Name,
		Version: ver.Version,
		Enabled: ver.Enabled,
		System:  ver.System,
		Actor:   ver.Actor,
	}

	// serialize settings
	b, _ := json.Marshal(ver.Settings)
	_ = dbVer.Settings.UnmarshalJSON(b)

	return
}
