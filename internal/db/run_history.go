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

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/e154/smart-home/pkg/apperr"
	"github.com/pkg/errors"
)

// RunHistory ...
type RunHistory struct {
	*Common
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
func (n RunHistory) Add(ctx context.Context, story *RunStory) (id int64, err error) {
	if err = n.DB(ctx).Create(&story).Error; err != nil {
		err = errors.Wrap(apperr.ErrRunStoryAdd, err.Error())
		return
	}
	id = story.Id
	return
}

// Update ...
func (n RunHistory) Update(ctx context.Context, m *RunStory) (err error) {
	q := map[string]interface{}{
		"end": m.End,
	}
	if err = n.DB(ctx).Model(&RunStory{Id: m.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrRunStoryUpdate, err.Error())
	}
	return
}

// List ...
func (n *RunHistory) List(ctx context.Context, limit, offset int, orderBy, sort string, from *time.Time) (list []*RunStory, total int64, err error) {

	list = make([]*RunStory, 0)
	q := n.DB(ctx).Model(&RunStory{})

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if from != nil {
		q = q.Where("start > ?", from.UTC().Format(time.RFC3339))
	}

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrRunStoryList, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset)

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrRunStoryList, err.Error())
	}

	return
}

// DeleteOldest ...
func (n *RunHistory) DeleteOldest(ctx context.Context, days int) (err error) {
	story := &RunStory{}
	if err = n.DB(ctx).Last(&story).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogDelete, err.Error())
		return
	}
	err = n.DB(ctx).Delete(&RunStory{},
		fmt.Sprintf(`start < CAST('%s' AS DATE) - interval '%d days'`,
			story.Start.UTC().Format("2006-01-02 15:04:05"), days)).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrEntityStorageDelete, err.Error())
	}
	return
}
