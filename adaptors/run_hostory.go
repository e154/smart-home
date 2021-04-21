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

// RunStory ...
type RunHistory struct {
	table *db.RunHistory
	db    *gorm.DB
}

// GetRunHistoryAdaptor ...
func GetRunHistoryAdaptor(d *gorm.DB) *RunHistory {
	return &RunHistory{
		table: &db.RunHistory{Db: d},
		db:    d,
	}
}

// Add ...
func (n *RunHistory) Add(story *m.RunStory) (id int64, err error) {

	var dbVer *db.RunStory
	dbVer, err = n.toDb(story)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

// Update ...
func (n *RunHistory) Update(story *m.RunStory) (err error) {

	var dbVer *db.RunStory
	dbVer, err = n.toDb(story)
	err = n.table.Update(dbVer)
	return
}

// List ...
func (n *RunHistory) List(limit, offset int64, orderBy, sort string) (list []*m.RunStory, total int64, err error) {
	var dbList []*db.RunStory
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.RunStory, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

func (n *RunHistory) fromDb(dbVer *db.RunStory) (story *m.RunStory) {
	story = &m.RunStory{
		Id:    dbVer.Id,
		Start: dbVer.Start,
		End:   dbVer.End,
	}

	return
}

func (n *RunHistory) toDb(story *m.RunStory) (dbVer *db.RunStory, err error) {
	dbVer = &db.RunStory{
		Id:    story.Id,
		Start: story.Start,
		End:   story.End,
	}

	return
}
