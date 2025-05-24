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

var _ adaptors.AlexaIntentRepo = (*AlexaIntent)(nil)

// AlexaIntent ...
type AlexaIntent struct {
	table *db.AlexaIntents
	db    *gorm.DB
}

// GetAlexaIntentAdaptor ...
func GetAlexaIntentAdaptor(d *gorm.DB) *AlexaIntent {
	return &AlexaIntent{
		table: &db.AlexaIntents{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *AlexaIntent) Add(ctx context.Context, ver *m.AlexaIntent) error {
	return n.table.Add(ctx, n.toDb(ver))
}

// GetByName ...
func (n *AlexaIntent) GetByName(ctx context.Context, name string) (ver *m.AlexaIntent, err error) {

	var dbVer *db.AlexaIntent
	if dbVer, err = n.table.GetByName(ctx, name); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *AlexaIntent) Update(ctx context.Context, ver *m.AlexaIntent) (err error) {
	err = n.table.Update(ctx, n.toDb(ver))
	return
}

// Delete ...
func (n *AlexaIntent) Delete(ctx context.Context, ver *m.AlexaIntent) (err error) {
	err = n.table.Delete(ctx, n.toDb(ver))
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
