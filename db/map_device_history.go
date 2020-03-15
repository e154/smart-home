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
	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"time"
)

type MapDeviceHistories struct {
	Db *gorm.DB
}

type MapDeviceHistory struct {
	Id          int64 `gorm:"primary_key"`
	MapDevice   *MapDevice
	MapDeviceId int64
	Type        common.LogLevel
	Description string
	CreatedAt   time.Time
}

func (d *MapDeviceHistory) TableName() string {
	return "map_device_history"
}

func (m MapDeviceHistories) Add(story MapDeviceHistory) (id int64, err error) {
	if err = m.Db.Create(&story).Error; err != nil {
		return
	}
	id = story.Id
	return
}

func (m MapDeviceHistories) GetAllByDeviceId(id int64, limit, offset int) (list []*MapDeviceHistory, total int64, err error) {

	if err = m.Db.Model(MapDeviceHistory{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapDeviceHistory, 0)
	err = m.Db.Model(&MapDeviceHistory{}).
		Limit(limit).
		Offset(offset).
		Where("map_device_id = ?", id).
		Find(&list).Error

	return
}
