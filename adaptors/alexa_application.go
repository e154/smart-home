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

type AlexaApplication struct {
	table *db.AlexaApplications
	db    *gorm.DB
}

func GetAlexaApplicationAdaptor(d *gorm.DB) *AlexaApplication {
	return &AlexaApplication{
		table: &db.AlexaApplications{Db: d},
		db:    d,
	}
}

func (n *AlexaApplication) Add(app *m.AlexaApplication) (id int64, err error) {

	var dbVer *db.AlexaApplication
	dbVer, err = n.toDb(app)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *AlexaApplication) GetById(appId int64) (app *m.AlexaApplication, err error) {

	var dbVer *db.AlexaApplication
	if dbVer, err = n.table.GetById(appId); err != nil {
		return
	}

	app = n.fromDb(dbVer)

	return
}

func (n *AlexaApplication) Update(app *m.AlexaApplication) (err error) {

	var dbVer *db.AlexaApplication
	dbVer, err = n.toDb(app)
	err = n.table.Update(dbVer)
	return
}

func (n *AlexaApplication) Delete(appId int64) (err error) {
	err = n.table.Delete(appId)
	return
}

func (n *AlexaApplication) List(limit, offset int64, orderBy, sort string) (list []*m.AlexaApplication, total int64, err error) {
	var dbList []*db.AlexaApplication
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.AlexaApplication, 0)
	for _, dbVer := range dbList {
		app := n.fromDb(dbVer)
		list = append(list, app)
	}

	return
}

func (n *AlexaApplication) fromDb(dbVer *db.AlexaApplication) (app *m.AlexaApplication) {
	app = &m.AlexaApplication{
		Id:            dbVer.Id,
		ApplicationId: dbVer.ApplicationId,
		Description:   dbVer.Description,
		Intents:       nil,
		CreatedAt:     dbVer.CreatedAt,
		UpdatedAt:     dbVer.UpdatedAt,
	}

	return
}

func (n *AlexaApplication) toDb(ver *m.AlexaApplication) (dbVer *db.AlexaApplication, err error) {

	dbVer = &db.AlexaApplication{
		Id:            ver.Id,
		ApplicationId: ver.ApplicationId,
		Description:   ver.Description,
	}

	return
}
