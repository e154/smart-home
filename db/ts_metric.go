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

// TsMetrics ...
type TsMetrics struct {
	Db *gorm.DB
}

// TsMetric ...
type TsMetric struct {
	Id          int64 `gorm:"primary_key"`
	MapDevice   *MapDevice
	MapDeviceId int64
	MetricType  common.MetricType
	Buckets     []*TsBucket
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// TableName ...
func (d *TsMetric) TableName() string {
	return "ts_metric"
}

// Add ...
func (n TsMetrics) Add(node *TsMetric) error {
	return n.Db.Create(&node).Error
}
