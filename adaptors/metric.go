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
	"github.com/e154/smart-home/common"
	"time"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/cache"
	"github.com/e154/smart-home/system/orm"
	"gorm.io/gorm"
)

// IMetric ...
type IMetric interface {
	Add(ctx context.Context, ver *m.Metric) (id int64, err error)
	GetByIdWithData(ctx context.Context, id int64, from, to *time.Time, metricRange *common.MetricRange) (metric *m.Metric, err error)
	Update(ctx context.Context, ver *m.Metric) error
	Delete(ctx context.Context, deviceId int64) (err error)
	AddMultiple(ctx context.Context, items []*m.Metric) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.Metric, total int64, err error)
	Search(ctx context.Context, query string, limit, offset int) (list []*m.Metric, total int64, err error)
	fromDb(dbVer *db.Metric) (ver *m.Metric)
	toDb(ver *m.Metric) (dbVer *db.Metric)
}

// Metric ...
type Metric struct {
	IMetric
	table *db.Metrics
	db    *gorm.DB
	c     cache.Cache
	orm   *orm.Orm
}

// GetMetricAdaptor ...
func GetMetricAdaptor(d *gorm.DB, orm *orm.Orm) IMetric {
	c, _ := cache.NewCache("memory", `{"interval":3600}`)
	return &Metric{
		table: &db.Metrics{Db: d},
		db:    d,
		c:     c,
		orm:   orm,
	}
}

// Add ...
func (n *Metric) Add(ctx context.Context, ver *m.Metric) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(ver))
	return
}

// GetByIdWithData ...
func (n *Metric) GetByIdWithData(ctx context.Context, id int64, from, to *time.Time, metricRange *common.MetricRange) (metric *m.Metric, err error) {
	var dbVer *db.Metric
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)

	var optionItems = make([]string, len(metric.Options.Items))
	for i, item := range metric.Options.Items {
		optionItems[i] = item.Name
	}

	metricBucketAdaptor := GetMetricBucketAdaptor(n.db, n.orm)
	if metric.Data, err = metricBucketAdaptor.List(ctx, from, to, id, optionItems, metricRange); err != nil {
		log.Error(err.Error())
		return
	}

	return
}

// Update ...
func (n *Metric) Update(ctx context.Context, ver *m.Metric) error {
	return n.table.Update(ctx, n.toDb(ver))
}

// Delete ...
func (n *Metric) Delete(ctx context.Context, deviceId int64) (err error) {
	err = n.table.Delete(ctx, deviceId)
	return
}

// AddMultiple ...
func (n *Metric) AddMultiple(ctx context.Context, items []*m.Metric) (err error) {

	insertRecords := make([]*db.Metric, 0, len(items))
	for _, ver := range items {
		insertRecords = append(insertRecords, n.toDb(ver))
	}

	err = n.table.AddMultiple(ctx, insertRecords)
	return
}

// List ...
func (n *Metric) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.Metric, total int64, err error) {
	var dbList []*db.Metric
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Metric, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Search ...
func (n *Metric) Search(ctx context.Context, query string, limit, offset int) (list []*m.Metric, total int64, err error) {
	var dbList []*db.Metric
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Metric, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

func (n *Metric) fromDb(dbVer *db.Metric) (ver *m.Metric) {
	ver = &m.Metric{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
		Type:        dbVer.Type,
	}

	// deserialize options
	b, _ := dbVer.Options.MarshalJSON()
	_ = json.Unmarshal(b, &ver.Options)

	if dbVer.Data != nil && len(dbVer.Data) > 0 {
		metricBucketAdaptor := GetMetricBucketAdaptor(n.db, n.orm)
		ver.Data = make([]*m.MetricDataItem, len(dbVer.Data))
		for i, dbVer := range dbVer.Data {
			ver.Data[i] = metricBucketAdaptor.fromDb(dbVer)
		}
	}

	ver.RangesByType()

	return
}

func (n *Metric) toDb(ver *m.Metric) (dbVer *db.Metric) {
	dbVer = &db.Metric{
		Id:          ver.Id,
		Name:        ver.Name,
		Type:        ver.Type,
		Description: ver.Description,
	}

	// serialize options
	b, _ := json.Marshal(ver.Options)
	_ = dbVer.Options.UnmarshalJSON(b)

	return
}
