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

	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Logs ...
type Logs struct {
	Db *gorm.DB
}

// Log ...
type Log struct {
	Id        int64 `gorm:"primary_key"`
	Body      string
	Level     common.LogLevel
	Owner     string
	CreatedAt time.Time
}

// LogQuery ...
type LogQuery struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Levels    []string   `json:"levels"`
}

// TableName ...
func (m *Log) TableName() string {
	return "logs"
}

// Add ...
func (n Logs) Add(v *Log) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogAdd, err.Error())
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n Logs) GetById(id int64) (v *Log, err error) {
	v = &Log{Id: id}
	if err = n.Db.First(&v).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrLogNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrLogGet, err.Error())
	}
	return
}

// Delete ...
func (n Logs) Delete(mapId int64) (err error) {
	if err = n.Db.Delete(&Log{Id: mapId}).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogDelete, err.Error())
	}
	return
}

// List ...
func (n *Logs) List(limit, offset int64, orderBy, sort string, queryObj *LogQuery) (list []*Log, total int64, err error) {

	q := n.Db.Model(Log{})

	if queryObj != nil {
		if queryObj.StartDate != nil {
			q = q.Where("created_at >= ?", &queryObj.StartDate)
		}
		if queryObj.EndDate != nil {
			q = q.Where("created_at <= ?", &queryObj.EndDate)
		}
		if len(queryObj.Levels) > 0 {
			q = q.Where("level in (?)", queryObj.Levels)
		}
	}

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogList, err.Error())
		return
	}

	list = make([]*Log, 0)
	err = q.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrLogList, err.Error())
	}

	return
}

// Search ...
func (n *Logs) Search(query string, limit, offset int) (list []*Log, total int64, err error) {

	q := n.Db.Model(&Log{}).
		Where("body LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogList, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("body ASC")

	list = make([]*Log, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogList, err.Error())
	}

	return
}
