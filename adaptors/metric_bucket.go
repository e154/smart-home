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

package adaptors

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/orm"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
	"time"
)

type IMetricBucket interface {
	Add(ver m.MetricDataItem) error
	AddMultiple(items []m.MetricDataItem) (err error)
	SimpleListWithSoftRange(_from, _to *time.Time, metricId int64, _metricRange *string, optionItems []string) (list []m.MetricDataItem, err error)
	Simple24HPreview(metricId int64, optionItems []string) (list []m.MetricDataItem, err error)
	DeleteOldest(days int) (err error)
	DeleteById(id int64) (err error)
	DeleteByMetricId(metricId int64) (err error)
	CreateHypertable() (err error)
	fromDb(dbVer db.MetricBucket) (ver m.MetricDataItem)
	toDb(ver m.MetricDataItem) (dbVer db.MetricBucket)
}

// MetricDataItem ...
type MetricBucket struct {
	IMetricBucket
	table *db.MetricBuckets
	db    *gorm.DB
	orm   *orm.Orm
}

// GetMetricBucketAdaptor ...
func GetMetricBucketAdaptor(d *gorm.DB, orm *orm.Orm) IMetricBucket {
	return &MetricBucket{
		table: &db.MetricBuckets{Db: d},
		db:    d,
		orm:   orm,
	}
}

// Add ...
func (n *MetricBucket) Add(ver m.MetricDataItem) error {
	return n.table.Add(n.toDb(ver))
}

// AddMultiple ...
func (n *MetricBucket) AddMultiple(items []m.MetricDataItem) (err error) {

	insertRecords := make([]interface{}, 0, len(items))
	for _, ver := range items {
		dbVer := n.toDb(ver)
		insertRecords = append(insertRecords, dbVer)
	}

	err = gormbulk.BulkInsert(n.db, insertRecords, len(insertRecords))

	return
}

// SimpleListWithSoftRange ...
func (n *MetricBucket) SimpleListWithSoftRange(_from, _to *time.Time, metricId int64, _metricRange *string, optionItems []string) (list []m.MetricDataItem, err error) {

	var dbList []db.MetricBucket

	if _metricRange != nil {
		if dbList, err = n.table.SimpleListByRangeType(metricId, common.StringValue(_metricRange), optionItems); err != nil {
			return
		}
	}

	if _from != nil && _to != nil {
		if dbList, err = n.table.SimpleListWithSoftRange(common.TimeValue(_from), common.TimeValue(_to), metricId, optionItems); err != nil {
			return
		}
	}

	list = make([]m.MetricDataItem, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// Simple24HPreview ...
func (n *MetricBucket) Simple24HPreview(metricId int64, optionItems []string) (list []m.MetricDataItem, err error) {
	var dbList []db.MetricBucket
	if dbList, err = n.table.Simple24HPreview(metricId, optionItems); err != nil {
		return
	}

	list = make([]m.MetricDataItem, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// DeleteOldest ...
func (n *MetricBucket) DeleteOldest(days int) (err error) {
	err = n.table.DeleteOldest(days)
	return
}

// DeleteById ...
func (n *MetricBucket) DeleteById(id int64) (err error) {
	err = n.table.DeleteById(id)
	return
}

// DeleteByMetricId ...
func (n *MetricBucket) DeleteByMetricId(metricId int64) (err error) {
	err = n.table.DeleteByMetricId(metricId)
	return
}

// CreateHypertable ...
func (n *MetricBucket) CreateHypertable() (err error) {
	err = n.table.CreateHypertable()
	return
}

func (n *MetricBucket) fromDb(dbVer db.MetricBucket) (ver m.MetricDataItem) {
	ver = m.MetricDataItem{
		Value:    dbVer.Value,
		MetricId: dbVer.MetricId,
		Time:     dbVer.Time,
	}

	return
}

func (n *MetricBucket) toDb(ver m.MetricDataItem) (dbVer db.MetricBucket) {
	dbVer = db.MetricBucket{
		Value:    ver.Value,
		MetricId: ver.MetricId,
		Time:     ver.Time,
	}

	return
}
