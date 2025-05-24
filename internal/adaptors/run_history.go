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
	"time"

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	m "github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.RunHistoryRepo = (*RunHistory)(nil)

// RunHistory ...
type RunHistory struct {
	table *db.RunHistory
	db    *gorm.DB
}

// GetRunHistoryAdaptor ...
func GetRunHistoryAdaptor(d *gorm.DB) *RunHistory {
	return &RunHistory{
		table: &db.RunHistory{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *RunHistory) Add(ctx context.Context, story *m.RunStory) (id int64, err error) {

	var dbVer *db.RunStory
	if dbVer, err = n.toDb(story); err != nil {
		return
	}
	id, err = n.table.Add(ctx, dbVer)
	return
}

// Update ...
func (n *RunHistory) Update(ctx context.Context, story *m.RunStory) (err error) {

	var dbVer *db.RunStory
	if dbVer, err = n.toDb(story); err != nil {
		return
	}
	err = n.table.Update(ctx, dbVer)
	return
}

// List ...
func (n *RunHistory) List(ctx context.Context, limit, offset int64, orderBy, sort string, from *time.Time) (list []*m.RunStory, total int64, err error) {
	var dbList []*db.RunStory
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, from); err != nil {
		return
	}

	list = make([]*m.RunStory, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// DeleteOldest ...
func (n *RunHistory) DeleteOldest(ctx context.Context, days int) (err error) {
	err = n.table.DeleteOldest(ctx, days)
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
