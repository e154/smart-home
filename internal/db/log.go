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
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Logs ...
type Logs struct {
	*Common
}

// Log ...
type Log struct {
	Id        int64 `gorm:"primary_key"`
	Body      string
	Level     pkgCommon.LogLevel
	Owner     string
	CreatedAt time.Time `gorm:"<-:create"`
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
func (n Logs) Add(ctx context.Context, v *Log) (id int64, err error) {
	if err = n.DB(ctx).Create(&v).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogAdd, err.Error())
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n Logs) GetById(ctx context.Context, id int64) (v *Log, err error) {
	v = &Log{Id: id}
	if err = n.DB(ctx).First(&v).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrLogNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrLogGet, err.Error())
	}
	return
}

// Delete ...
func (n Logs) Delete(ctx context.Context, mapId int64) (err error) {
	if err = n.DB(ctx).Delete(&Log{Id: mapId}).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogDelete, err.Error())
	}
	return
}

// List ...
func (n *Logs) List(ctx context.Context, limit, offset int, orderBy, sort string, queryObj *LogQuery) (list []*Log, total int64, err error) {

	q := n.DB(ctx).Model(Log{})

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
func (n *Logs) Search(ctx context.Context, query string, limit, offset int) (list []*Log, total int64, err error) {

	q := n.DB(ctx).Model(&Log{}).
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

// DeleteOldest ...
func (n *Logs) DeleteOldest(ctx context.Context, days int) (err error) {

	log := &Log{}
	if err = n.DB(ctx).Last(&log).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogDelete, err.Error())
		return
	}
	err = n.DB(ctx).Delete(&Log{},
		fmt.Sprintf(`created_at < CAST('%s' AS DATE) - interval '%d days'`,
			log.CreatedAt.UTC().Format("2006-01-02 15:04:05"), days)).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrLogDelete, err.Error())
	}
	return
}

// AddMultiple ...
func (n *Logs) AddMultiple(ctx context.Context, logs []*Log) (err error) {
	if err = n.DB(ctx).Create(&logs).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogAdd, err.Error())
	}
	return
}
