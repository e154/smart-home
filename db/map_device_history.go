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
	"database/sql"
	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"time"
)

// MapDeviceHistories ...
type MapDeviceHistories struct {
	Db *gorm.DB
}

// MapDeviceHistory ...
type MapDeviceHistory struct {
	Id           int64 `gorm:"primary_key"`
	MapDevice    *MapDevice
	MapDeviceId  int64
	MapElement   *MapElement
	MapElementId int64
	LogLevel     common.LogLevel
	Type         common.MapDeviceHistoryType
	Description  string
	CreatedAt    time.Time
}

// TableName ...
func (d *MapDeviceHistory) TableName() string {
	return "map_device_history"
}

// Add ...
func (m MapDeviceHistories) Add(story MapDeviceHistory) (id int64, err error) {
	if err = m.Db.Create(&story).Error; err != nil {
		return
	}
	id = story.Id
	return
}

// ListByDeviceId ...
func (m MapDeviceHistories) ListByDeviceId(id int64, limit, offset int) (list []*MapDeviceHistory, total int64, err error) {

	if err = m.Db.Model(MapDeviceHistory{}).
		Where("map_device_id = ?", id).
		Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapDeviceHistory, 0)
	err = m.Db.Model(&MapDeviceHistory{}).
		Limit(limit).
		Offset(offset).
		Where("map_device_id = ?", id).
		Order("id desc").
		Preload("MapElement").
		Find(&list).Error

	return
}

// ListByElementId ...
func (m MapDeviceHistories) ListByElementId(id int64, limit, offset int) (list []*MapDeviceHistory, total int64, err error) {

	if err = m.Db.Model(MapDeviceHistory{}).
		Where("id = ?", id).
		Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapDeviceHistory, 0)
	err = m.Db.Model(&MapDeviceHistory{}).
		Limit(limit).
		Offset(offset).
		Where("id = ?", id).
		Order("id desc").
		Preload("MapElement").
		Find(&list).Error

	return
}

// List ...
func (m MapDeviceHistories) List(limit, offset int) (list []*MapDeviceHistory, err error) {

	list = make([]*MapDeviceHistory, 0)
	err = m.Db.Model(&MapDeviceHistory{}).
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Preload("MapElement").
		Find(&list).Error

	return
}

// ListByMapId ...
func (m MapDeviceHistories) ListByMapId(mapId int64, limit, offset int, orderBy, sort string) (list []*MapDeviceHistory, total int64, err error) {

	err = m.Db.Raw(`select count(mdh.*)
from map_elements me
         left join map_device_history mdh on me.id = mdh.map_element_id
where me.map_id = ? and mdh notnull`, mapId).Count(&total).Error

	if err != nil {
		return
	}

	if orderBy == "" {
		orderBy = "id"
	}

	orderBy = "mdh." + orderBy

	//TODO sort fix
	if sort == "" {
		sort = "DESC"
	}

	var rows *sql.Rows
	rows, err = m.Db.Raw(`select mdh.id, mdh.map_device_id, mdh.map_element_id, 
mdh.type, mdh.description, mdh.created_at, me.id, me.name, me.description
from map_device_history mdh
left join map_elements me on mdh.map_element_id = me.id
where mdh.map_element_id in (select me.id
                             from map_elements me
                             where me.map_id = ?)
order by mdh.id DESC limit ? offset ?`, mapId, limit, offset).Rows()

	if err != nil {
		return
	}

	list = make([]*MapDeviceHistory, 0)

	for rows.Next() {
		item := &MapDeviceHistory{
			MapElement: &MapElement{},
		}
		err = rows.Scan(&item.Id, &item.MapDeviceId, &item.MapElementId, &item.Type, &item.Description,
			&item.CreatedAt, &item.MapElement.Id, &item.MapElement.Name, &item.MapElement.Description)
		list = append(list, item)
	}

	return
}
