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
	m "github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.AlexaSkillRepo = (*AlexaSkill)(nil)

// AlexaSkill ...
type AlexaSkill struct {
	table *db.AlexaSkills
	db    *gorm.DB
}

// GetAlexaSkillAdaptor ...
func GetAlexaSkillAdaptor(d *gorm.DB) *AlexaSkill {
	return &AlexaSkill{
		table: &db.AlexaSkills{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *AlexaSkill) Add(ctx context.Context, app *m.AlexaSkill) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(app))
	return
}

// GetById ...
func (n *AlexaSkill) GetById(ctx context.Context, appId int64) (app *m.AlexaSkill, err error) {

	var dbVer *db.AlexaSkill
	if dbVer, err = n.table.GetById(ctx, appId); err != nil {
		return
	}

	app = n.fromDb(dbVer)

	return
}

// Update ...
func (n *AlexaSkill) Update(ctx context.Context, params *m.AlexaSkill) (err error) {
	err = n.table.Update(ctx, n.toDb(params))
	return
}

// Delete ...
func (n *AlexaSkill) Delete(ctx context.Context, appId int64) (err error) {
	err = n.table.Delete(ctx, appId)
	return
}

// List ...
func (n *AlexaSkill) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.AlexaSkill, total int64, err error) {
	var dbList []*db.AlexaSkill
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.AlexaSkill, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

// ListEnabled ...
func (n *AlexaSkill) ListEnabled(ctx context.Context, limit, offset int64) (list []*m.AlexaSkill, err error) {
	var dbList []*db.AlexaSkill
	if dbList, err = n.table.ListEnabled(ctx, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.AlexaSkill, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

func (n *AlexaSkill) fromDb(dbVer *db.AlexaSkill) (app *m.AlexaSkill) {

	app = &m.AlexaSkill{
		Id:          dbVer.Id,
		SkillId:     dbVer.SkillId,
		Description: dbVer.Description,
		Status:      dbVer.Status,
		ScriptId:    dbVer.ScriptId,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}

	intentAdaptor := GetAlexaIntentAdaptor(n.db)
	for _, dbVer := range dbVer.Intents {
		app.Intents = append(app.Intents, intentAdaptor.fromDb(dbVer))
	}

	scriptAdaptor := GetScriptAdaptor(n.db)
	if dbVer.ScriptId != nil {
		app.Script, _ = scriptAdaptor.fromDb(dbVer.Script)
	}

	return
}

func (n *AlexaSkill) toDb(ver *m.AlexaSkill) (dbVer *db.AlexaSkill) {

	dbVer = &db.AlexaSkill{
		Id:          ver.Id,
		SkillId:     ver.SkillId,
		Description: ver.Description,
		ScriptId:    ver.ScriptId,
		Status:      ver.Status,
	}

	intentAdaptor := GetAlexaIntentAdaptor(n.db)
	for _, ver := range ver.Intents {
		dbVer.Intents = append(dbVer.Intents, intentAdaptor.toDb(ver))
	}

	return
}
