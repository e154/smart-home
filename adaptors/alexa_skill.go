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

type IAlexaSkill interface {
	Add(app *m.AlexaSkill) (id int64, err error)
	GetById(appId int64) (app *m.AlexaSkill, err error)
	Update(params *m.AlexaSkill) (err error)
	Delete(appId int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.AlexaSkill, total int64, err error)
	ListEnabled(limit, offset int64) (list []*m.AlexaSkill, err error)
	fromDb(dbVer *db.AlexaSkill) (app *m.AlexaSkill)
	toDb(ver *m.AlexaSkill) (dbVer *db.AlexaSkill)
}

// AlexaSkill ...
type AlexaSkill struct {
	IAlexaSkill
	table *db.AlexaSkills
	db    *gorm.DB
}

// GetAlexaSkillAdaptor ...
func GetAlexaSkillAdaptor(d *gorm.DB) IAlexaSkill {
	return &AlexaSkill{
		table: &db.AlexaSkills{Db: d},
		db:    d,
	}
}

// Add ...
func (n *AlexaSkill) Add(app *m.AlexaSkill) (id int64, err error) {
	id, err = n.table.Add(n.toDb(app))
	return
}

// GetById ...
func (n *AlexaSkill) GetById(appId int64) (app *m.AlexaSkill, err error) {

	var dbVer *db.AlexaSkill
	if dbVer, err = n.table.GetById(appId); err != nil {
		return
	}

	app = n.fromDb(dbVer)

	return
}

// Update ...
func (n *AlexaSkill) Update(params *m.AlexaSkill) (err error) {

	var app *db.AlexaSkill
	if app, err = n.table.GetById(params.Id); err != nil {
		return
	}

	tx := n.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	intentAdaptor := GetAlexaIntentAdaptor(tx)

	// обновление, либо удаление alexa intent
	dbParams := n.toDb(params)
	for _, intent := range app.Intents {
		var exist bool
		for _, parIntent := range dbParams.Intents {
			if intent.Name == parIntent.Name {
				exist = true
			}
		}
		if !exist {
			if err = intentAdaptor.Delete(intentAdaptor.fromDb(intent)); err != nil {
				return
			}
		} else {
			if err = intentAdaptor.Update(intentAdaptor.fromDb(intent)); err != nil {
				return
			}
		}
	}

	// добавление alexa intent
	for _, parIntent := range params.Intents {
		var exist bool
		for _, intent := range app.Intents {
			if intent.Name == parIntent.Name {
				exist = true
			}
		}
		if !exist {
			if err = intentAdaptor.Add(parIntent); err != nil {
				return
			}
		}
	}

	table := &db.AlexaSkills{Db: tx}
	err = table.Update(n.toDb(params))

	tx.Commit()

	return
}

// Delete ...
func (n *AlexaSkill) Delete(appId int64) (err error) {
	err = n.table.Delete(appId)
	return
}

// List ...
func (n *AlexaSkill) List(limit, offset int64, orderBy, sort string) (list []*m.AlexaSkill, total int64, err error) {
	var dbList []*db.AlexaSkill
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.AlexaSkill, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

// ListEnabled ...
func (n *AlexaSkill) ListEnabled(limit, offset int64) (list []*m.AlexaSkill, err error) {
	var dbList []*db.AlexaSkill
	if dbList, err = n.table.ListEnabled(limit, offset); err != nil {
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
