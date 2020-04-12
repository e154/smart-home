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
	"github.com/jinzhu/gorm"
	"time"
)

// Metrics ...
type Metrics struct {
	Db *gorm.DB
}

// Metric ...
type Metric struct {
	Id          int64 `gorm:"primary_key"`
	MapDevice   *MapDevice
	MapDeviceId int64
	Buckets     []MetricBucket
	Name        string
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// TableName ...
func (Metric) TableName() string {
	return "metrics"
}

// Add ...
func (n Metrics) Add(metric Metric) (id int64, err error) {
	if err = n.Db.Create(&metric).Error; err != nil {
		return
	}
	id = metric.Id
	return
}

// Update ...
func (n Metrics) Update(m Metric) (err error) {
	q := map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
	}
	err = n.Db.Model(&Metric{}).Updates(q).Error
	return
}

// GetByMapDeviceId ...
func (n Metrics) GetByMapDeviceId(mapDeviceId int64, name string) (list []Metric, err error) {
	list = make([]Metric, 0)
	err = n.Db.Model(&Metric{}).Where("map_device_id = ? and name = ?", mapDeviceId, name).Find(&list).Error
	return
}

// Delete ...
func (n Metrics) Delete(id int64) error {
	return n.Db.Delete(&Metric{Id: id}).Error
}
