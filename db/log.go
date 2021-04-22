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

package db

import (
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"time"
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
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n Logs) GetById(mapId int64) (v *Log, err error) {
	v = &Log{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

// Delete ...
func (n Logs) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&Log{Id: mapId}).Error
	return
}

// List ...
func (n *Logs) List(limit, offset int64, orderBy, sort string, queryObj *LogQuery) (list []*Log, total int64, err error) {

	if err = n.Db.Model(Log{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Log, 0)
	q := n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy))

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

	err = q.
		Find(&list).
		Error

	if err != nil {
		log.Error(err.Error())
	}

	return
}

// Search ...
func (n *Logs) Search(query string, limit, offset int) (list []*Log, total int64, err error) {

	q := n.Db.Model(&Log{}).
		Where("body LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("body ASC")


	list = make([]*Log, 0)
	err = q.Find(&list).Error

	return
}
