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
	"sort"

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	m "github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.ScriptRepo = (*Script)(nil)

// Script ...
type Script struct {
	table *db.Scripts
	db    *gorm.DB
}

// GetScriptAdaptor ...
func GetScriptAdaptor(d *gorm.DB) *Script {
	return &Script{
		table: &db.Scripts{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *Script) Add(ctx context.Context, script *m.Script) (id int64, err error) {
	dbScript := n.toDb(script)
	id, err = n.table.Add(ctx, dbScript)
	return
}

// GetById ...
func (n *Script) GetById(ctx context.Context, scriptId int64) (script *m.Script, err error) {

	var dbScript *db.Script
	if dbScript, err = n.table.GetById(ctx, scriptId); err != nil {
		return
	}

	script, _ = n.fromDb(dbScript)

	return
}

// GetByName ...
func (n *Script) GetByName(ctx context.Context, name string) (script *m.Script, err error) {

	var dbScript *db.Script
	if dbScript, err = n.table.GetByName(ctx, name); err != nil {
		return
	}

	script, _ = n.fromDb(dbScript)

	return
}

// Update ...
func (n *Script) Update(ctx context.Context, script *m.Script) (err error) {
	dbScript := n.toDb(script)
	err = n.table.Update(ctx, dbScript)
	return
}

// Delete ...
func (n *Script) Delete(ctx context.Context, scriptId int64) (err error) {
	err = n.table.Delete(ctx, scriptId)
	return
}

// List ...
func (n *Script) List(ctx context.Context, limit, offset int64, orderBy, sort string, query *string, ids *[]uint64) (list []*m.Script, total int64, err error) {

	if sort == "" {
		sort = "id"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.Script
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, query, ids); err != nil {
		return
	}

	list = make([]*m.Script, 0)
	for _, dbScript := range dbList {
		script, _ := n.fromDb(dbScript)
		list = append(list, script)
	}

	return
}

// Search ...
func (n *Script) Search(ctx context.Context, query string, limit, offset int64) (list []*m.Script, total int64, err error) {
	var dbList []*db.Script
	if dbList, total, err = n.table.Search(ctx, query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Script, 0)
	for _, dbScript := range dbList {
		dev, _ := n.fromDb(dbScript)
		list = append(list, dev)
	}

	return
}

func (n *Script) Statistic(ctx context.Context) (statistic *m.ScriptsStatistic, err error) {
	var dbVer *db.ScriptsStatistic
	if dbVer, err = n.table.Statistic(ctx); err != nil {
		return
	}
	statistic = &m.ScriptsStatistic{
		Total:        dbVer.Total,
		Used:         dbVer.Used,
		Unused:       dbVer.Unused,
		CoffeeScript: dbVer.CoffeeScript,
		TypeScript:   dbVer.TypeScript,
		JavaScript:   dbVer.JavaScript,
	}
	return
}

func (n *Script) fromDb(dbVer *db.Script) (ver *m.Script, err error) {
	ver = &m.Script{
		Id:          dbVer.Id,
		Lang:        dbVer.Lang,
		Name:        dbVer.Name,
		Source:      dbVer.Source,
		Description: dbVer.Description,
		Compiled:    dbVer.Compiled,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
		Info: &m.ScriptInfo{
			AlexaIntents:         dbVer.AlexaIntents,
			EntityActions:        dbVer.EntityActions,
			EntityScripts:        dbVer.EntityScripts,
			AutomationTriggers:   dbVer.AutomationTriggers,
			AutomationConditions: dbVer.AutomationConditions,
			AutomationActions:    dbVer.AutomationActions,
		},
	}
	if dbVer.Versions != nil {
		ver.Versions = make([]*m.ScriptVersion, 0, len(dbVer.Versions))
		for _, version := range dbVer.Versions {
			ver.Versions = append(ver.Versions, &m.ScriptVersion{
				Id:        version.Id,
				Lang:      version.Lang,
				Source:    version.Source,
				CreatedAt: version.CreatedAt,
			})
		}
		sort.Sort(ver.Versions)
	}
	return
}

func (n *Script) toDb(script *m.Script) (dbVer *db.Script) {
	dbVer = &db.Script{
		Id:          script.Id,
		Lang:        script.Lang,
		Name:        script.Name,
		Source:      script.Source,
		Description: script.Description,
		Compiled:    script.Compiled,
		CreatedAt:   script.CreatedAt,
		UpdatedAt:   script.UpdatedAt,
	}
	return
}
