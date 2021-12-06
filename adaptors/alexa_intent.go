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
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

// IAlexaIntent ...
type IAlexaIntent interface {
	Add(ver *m.AlexaIntent) (err error)
	GetByName(name string) (ver *m.AlexaIntent, err error)
	Update(ver *m.AlexaIntent) (err error)
	Delete(ver *m.AlexaIntent) (err error)
	fromDb(dbVer *db.AlexaIntent) (ver *m.AlexaIntent)
	toDb(ver *m.AlexaIntent) (dbVer *db.AlexaIntent)
}

// AlexaIntent ...
type AlexaIntent struct {
	IAlexaIntent
	table *db.AlexaIntents
	db    *gorm.DB
}

// GetAlexaIntentAdaptor ...
func GetAlexaIntentAdaptor(d *gorm.DB) IAlexaIntent {
	return &AlexaIntent{
		table: &db.AlexaIntents{Db: d},
		db:    d,
	}
}

// Add ...
func (n *AlexaIntent) Add(ver *m.AlexaIntent) error {
	return n.table.Add(n.toDb(ver))
}

// GetByName ...
func (n *AlexaIntent) GetByName(name string) (ver *m.AlexaIntent, err error) {

	var dbVer *db.AlexaIntent
	if dbVer, err = n.table.GetByName(name); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *AlexaIntent) Update(ver *m.AlexaIntent) (err error) {
	err = n.table.Update(n.toDb(ver))
	return
}

// Delete ...
func (n *AlexaIntent) Delete(ver *m.AlexaIntent) (err error) {
	err = n.table.Delete(n.toDb(ver))
	return
}

func (n *AlexaIntent) fromDb(dbVer *db.AlexaIntent) (ver *m.AlexaIntent) {
	ver = &m.AlexaIntent{
		Name:         dbVer.Name,
		AlexaSkillId: dbVer.AlexaSkillId,
		ScriptId:     dbVer.ScriptId,
		Description:  dbVer.Description,
		CreatedAt:    dbVer.CreatedAt,
		UpdatedAt:    dbVer.UpdatedAt,
	}

	if dbVer.Script != nil {
		scriptAdaptor := GetScriptAdaptor(n.db)
		ver.Script, _ = scriptAdaptor.fromDb(dbVer.Script)
	}

	return
}

func (n *AlexaIntent) toDb(ver *m.AlexaIntent) (dbVer *db.AlexaIntent) {

	dbVer = &db.AlexaIntent{
		Name:         ver.Name,
		AlexaSkillId: ver.AlexaSkillId,
		ScriptId:     ver.ScriptId,
		Description:  ver.Description,
	}

	return
}
