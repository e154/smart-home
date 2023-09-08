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

package db

import (
	"fmt"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// RunHistory ...
type RunHistory struct {
	Db *gorm.DB
}

// RunStory ...
type RunStory struct {
	Id    int64 `gorm:"primary_key"`
	Start time.Time
	End   *time.Time
}

// TableName ...
func (d *RunStory) TableName() string {
	return "run_history"
}

// Add ...
func (n RunHistory) Add(story *RunStory) (id int64, err error) {
	if err = n.Db.Create(&story).Error; err != nil {
		err = errors.Wrap(apperr.ErrRunStoryAdd, err.Error())
		return
	}
	id = story.Id
	return
}

// Update ...
func (n RunHistory) Update(m *RunStory) (err error) {
	q := map[string]interface{}{
		"end": m.End,
	}
	if err = n.Db.Model(&RunStory{Id: m.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrRunStoryUpdate, err.Error())
	}
	return
}

// List ...
func (n *RunHistory) List(limit, offset int, orderBy, sort string) (list []*RunStory, total int64, err error) {

	if err = n.Db.Model(RunStory{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrRunStoryList, err.Error())
		return
	}

	list = make([]*RunStory, 0)
	q := n.Db.Model(&RunStory{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrRunStoryList, err.Error())
	}

	return
}
