// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/internal/system/orm"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.MetricRepo = (*Metric)(nil)

// Metric ...
type Metric struct {
	table *db.Metrics
	db    *gorm.DB
	orm   *orm.Orm
}

// GetMetricAdaptor ...
func GetMetricAdaptor(d *gorm.DB, orm *orm.Orm) *Metric {
	return &Metric{
		table: &db.Metrics{&db.Common{Db: d}},
		db:    d,
		orm:   orm,
	}
}

// Add ...
func (n *Metric) Add(ctx context.Context, ver *models.Metric) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(ver))
	return
}

// GetByIdWithData ...
func (n *Metric) GetByIdWithData(ctx context.Context, id int64, from, to *time.Time, metricRange *common.MetricRange) (metric *models.Metric, err error) {
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
func (n *Metric) Update(ctx context.Context, ver *models.Metric) error {
	return n.table.Update(ctx, n.toDb(ver))
}

// Delete ...
func (n *Metric) Delete(ctx context.Context, deviceId int64) (err error) {
	err = n.table.Delete(ctx, deviceId)
	return
}

// AddMultiple ...
func (n *Metric) AddMultiple(ctx context.Context, items []*models.Metric) (err error) {

	insertRecords := make([]*db.Metric, 0, len(items))
	for _, ver := range items {
		insertRecords = append(insertRecords, n.toDb(ver))
	}

	err = n.table.AddMultiple(ctx, insertRecords)
	return
}

// List ...
func (n *Metric) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*models.Metric, total int64, err error) {
	var dbList []*db.Metric
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*models.Metric, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Search ...
func (n *Metric) Search(ctx context.Context, query string, limit, offset int) (list []*models.Metric, total int64, err error) {
	var dbList []*db.Metric
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
		return
	}

	list = make([]*models.Metric, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

func (n *Metric) fromDb(dbVer *db.Metric) (ver *models.Metric) {
	ver = &models.Metric{
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
		ver.Data = make([]*models.MetricDataItem, len(dbVer.Data))
		for i, dbVer := range dbVer.Data {
			ver.Data[i] = metricBucketAdaptor.fromDb(dbVer)
		}
	}

	ver.RangesByType()

	return
}

func (n *Metric) toDb(ver *models.Metric) (dbVer *db.Metric) {
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
