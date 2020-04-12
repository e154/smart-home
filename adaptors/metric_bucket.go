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

package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

// MetricBucket ...
type MetricBucket struct {
	table *db.MetricBuckets
	db    *gorm.DB
}

// GetMetricBucketAdaptor ...
func GetMetricBucketAdaptor(d *gorm.DB) *MetricBucket {
	return &MetricBucket{
		table: &db.MetricBuckets{Db: d},
		db:    d,
	}
}

// Add ...
func (n *MetricBucket) Add(ver m.MetricBucket) error {
	return n.table.Add(n.toDb(ver))
}

// AddMultiple ...
func (n *MetricBucket) AddMultiple(items []m.MetricBucket) (err error) {

	insertRecords := make([]interface{}, 0, len(items))
	for _, ver := range items {
		dbVer := n.toDb(ver)
		insertRecords = append(insertRecords, dbVer)
	}

	err = gormbulk.BulkInsert(n.db, insertRecords, len(insertRecords))

	return
}

func (n *MetricBucket) fromDb(dbVer *db.MetricBucket) (ver m.MetricBucket) {
	ver = m.MetricBucket{
		Value:          dbVer.Value,
		MetricMetricId: dbVer.MetricMetricId,
		Time:           dbVer.Time,
	}

	return
}

func (n *MetricBucket) toDb(ver m.MetricBucket) (dbVer db.MetricBucket) {
	dbVer = db.MetricBucket{
		Value:          ver.Value,
		MetricMetricId: ver.MetricMetricId,
		Time:           ver.Time,
	}

	return
}
