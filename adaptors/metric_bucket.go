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

package adaptors

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/orm"
)

// IMetricBucket ...
type IMetricBucket interface {
	Add(ctx context.Context, ver *m.MetricDataItem) error
	AddMultiple(ctx context.Context, items []*m.MetricDataItem) (err error)
	List(ctx context.Context, from, to *time.Time, metricId int64, optionItems []string, metricRange *common.MetricRange) (list []*m.MetricDataItem, err error)
	DeleteOldest(ctx context.Context, days int) (err error)
	DeleteById(ctx context.Context, id int64) (err error)
	DeleteByMetricId(ctx context.Context, metricId int64) (err error)
	fromDb(dbVer *db.MetricBucket) (ver *m.MetricDataItem)
	toDb(ver *m.MetricDataItem) (dbVer *db.MetricBucket)
}

// MetricBucket ...
type MetricBucket struct {
	IMetricBucket
	table *db.MetricBuckets
	db    *gorm.DB
}

// GetMetricBucketAdaptor ...
func GetMetricBucketAdaptor(d *gorm.DB, orm *orm.Orm) IMetricBucket {
	table := &db.MetricBuckets{Db: d}
	if orm != nil {
		table.Timescale = orm.CheckAvailableExtension("timescaledb") && orm.CheckInstalledExtension("timescaledb")
	}
	return &MetricBucket{
		table: table,
		db:    d,
	}
}

// Add ...
func (n *MetricBucket) Add(ctx context.Context, ver *m.MetricDataItem) error {
	return n.table.Add(ctx, n.toDb(ver))
}

// AddMultiple ...
func (n *MetricBucket) AddMultiple(ctx context.Context, items []*m.MetricDataItem) (err error) {

	insertRecords := make([]*db.MetricBucket, 0, len(items))
	for _, ver := range items {
		dbVer := n.toDb(ver)
		insertRecords = append(insertRecords, dbVer)
	}

	err = n.table.AddMultiple(ctx, insertRecords)

	return
}

// List ...
func (n *MetricBucket) List(ctx context.Context, from, to *time.Time, metricId int64, optionItems []string, metricRange *common.MetricRange) (list []*m.MetricDataItem, err error) {

	var dbList []*db.MetricBucket
	if dbList, err = n.table.List(ctx, metricId, optionItems, from, to, metricRange); err != nil {
		return
	}

	list = make([]*m.MetricDataItem, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// DeleteOldest ...
func (n *MetricBucket) DeleteOldest(ctx context.Context, days int) (err error) {
	err = n.table.DeleteOldest(ctx, days)
	return
}

// DeleteById ...
func (n *MetricBucket) DeleteById(ctx context.Context, id int64) (err error) {
	err = n.table.DeleteById(ctx, id)
	return
}

// DeleteByMetricId ...
func (n *MetricBucket) DeleteByMetricId(ctx context.Context, metricId int64) (err error) {
	err = n.table.DeleteByMetricId(ctx, metricId)
	return
}

func (n *MetricBucket) fromDb(dbVer *db.MetricBucket) (ver *m.MetricDataItem) {
	ver = &m.MetricDataItem{
		MetricId: dbVer.MetricId,
		Time:     dbVer.Time,
	}

	// deserialize value
	b, _ := dbVer.Value.MarshalJSON()
	value := make(map[string]interface{})
	_ = json.Unmarshal(b, &value)
	ver.Value = value

	return
}

func (n *MetricBucket) toDb(ver *m.MetricDataItem) (dbVer *db.MetricBucket) {
	dbVer = &db.MetricBucket{
		MetricId: ver.MetricId,
		Time:     ver.Time,
	}

	// serialize value
	b, _ := json.Marshal(ver.Value)
	_ = dbVer.Value.UnmarshalJSON(b)

	return
}
