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
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"
)

// MetricBuckets ...
type MetricBuckets struct {
	Db *gorm.DB
}

// MetricBucket ...
type MetricBucket struct {
	Value    json.RawMessage `gorm:"type:jsonb;not null"`
	Metric   *Metric
	MetricId int64
	Time     time.Time
}

// TableName ...
func (d *MetricBucket) TableName() string {
	return "metric_bucket"
}

// Add ...
func (n MetricBuckets) Add(metric MetricBucket) error {
	return n.Db.Create(&metric).Error
}

// SimpleListByRange ...
func (n *MetricBuckets) SimpleListByRange(from, to time.Time, metricId int64) (list []MetricBucket, err error) {
	list = make([]MetricBucket, 0)
	err = n.Db.Where("time between ? and ? and metric_id = ?", from, to, metricId).Find(&list).Error
	return
}

// SimpleListBySoftRange ...
func (n *MetricBuckets) SimpleListBySoftRange(from, to time.Time, metricId, num int64) (list []MetricBucket, err error) {
	list = make([]MetricBucket, 0)
	err = n.Db.Raw(`SELECT s.*
FROM (
         SELECT t.*,
                ROW_NUMBER() OVER (ORDER BY t.time) as rnk,
                COUNT(*) OVER ()                    as total_cnt
         FROM metric_bucket t
         WHERE t.metric_id = ? and t.time between ? AND ?) s
WHERE MOD(s.rnk, (total_cnt / ?)) = 0`, metricId, from.Format(time.RFC3339), to.Format(time.RFC3339), num).Scan(&list).Error
	return
}

// DeleteOldest ...
func (n *MetricBuckets) DeleteOldest(days int) (err error) {
	err = n.Db.Raw(`delete
    from metric_bucket
    where time < now() - interval '? days'`, days).Error
	return
}
