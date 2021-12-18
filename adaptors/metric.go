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
	"encoding/json"
	"time"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/cache"
	"github.com/e154/smart-home/system/orm"
	"github.com/jinzhu/gorm"
)

// IMetric ...
type IMetric interface {
	Add(ver m.Metric) (id int64, err error)
	GetById(id int64) (metric m.Metric, err error)
	GetByIdWithData(id int64, from, to *time.Time, metricRange *string) (metric m.Metric, err error)
	Update(ver m.Metric) error
	Delete(deviceId int64) (err error)
	AddMultiple(items []m.Metric) (err error)
	List(limit, offset int64, orderBy, sort string) (list []m.Metric, total int64, err error)
	Search(query string, limit, offset int) (list []m.Metric, total int64, err error)
	fromDb(dbVer db.Metric) (ver m.Metric)
	toDb(ver m.Metric) (dbVer db.Metric)
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
func (n *Metric) Add(ver m.Metric) (id int64, err error) {
	id, err = n.table.Add(n.toDb(ver))
	return
}

// GetById ...
func (n *Metric) GetById(id int64) (metric m.Metric, err error) {
	var dbVer db.Metric
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)

	var optionItems = make([]string, len(metric.Options.Items))
	for i, item := range metric.Options.Items {
		optionItems[i] = item.Name
	}

	metricBucketAdaptor := GetMetricBucketAdaptor(n.db, nil)
	if metric.Data, err = metricBucketAdaptor.Simple24HPreview(metric.Id, optionItems); err != nil {
		log.Error(err.Error())
		return
	}

	return
}

// GetByIdWithData ...
func (n *Metric) GetByIdWithData(id int64, from, to *time.Time, metricRange *string) (metric m.Metric, err error) {
	var dbVer db.Metric
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)

	var optionItems = make([]string, len(metric.Options.Items))
	for i, item := range metric.Options.Items {
		optionItems[i] = item.Name
	}

	metricBucketAdaptor := GetMetricBucketAdaptor(n.db, nil)
	if metric.Data, err = metricBucketAdaptor.SimpleListWithSoftRange(from, to, id, metricRange, optionItems); err != nil {
		log.Error(err.Error())
		return
	}

	return
}

// Update ...
func (n *Metric) Update(ver m.Metric) error {
	return n.table.Update(n.toDb(ver))
}

// Delete ...
func (n *Metric) Delete(deviceId int64) (err error) {
	err = n.table.Delete(deviceId)
	return
}

// AddMultiple ...
func (n *Metric) AddMultiple(items []m.Metric) (err error) {

	//TODO not work
	//insertRecords := make([]interface{}, 0, len(items))
	//for _, ver := range items {
	//	insertRecords = append(insertRecords, n.toDb(ver))
	//}
	//
	//err = gormbulk.BulkInsert(n.db, insertRecords, len(insertRecords))

	for _, ver := range items {
		if _, err = n.table.Add(n.toDb(ver)); err != nil {
			return
		}
	}

	return
}

// List ...
func (n *Metric) List(limit, offset int64, orderBy, sort string) (list []m.Metric, total int64, err error) {
	var dbList []db.Metric
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]m.Metric, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Search ...
func (n *Metric) Search(query string, limit, offset int) (list []m.Metric, total int64, err error) {
	var dbList []db.Metric
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]m.Metric, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

func (n *Metric) fromDb(dbVer db.Metric) (ver m.Metric) {
	ver = m.Metric{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
		Type:        dbVer.Type,
	}

	// deserialize options
	b, _ := dbVer.Options.MarshalJSON()
	json.Unmarshal(b, &ver.Options)

	if dbVer.Data != nil && len(dbVer.Data) > 0 {
		metricBucketAdaptor := GetMetricBucketAdaptor(n.db, nil)
		ver.Data = make([]m.MetricDataItem, len(dbVer.Data))
		for i, dbVer := range dbVer.Data {
			ver.Data[i] = metricBucketAdaptor.fromDb(dbVer)
		}
	}

	ver.RangesByType()

	return
}

func (n *Metric) toDb(ver m.Metric) (dbVer db.Metric) {
	dbVer = db.Metric{
		Id:          ver.Id,
		Name:        ver.Name,
		Type:        ver.Type,
		Description: ver.Description,
	}

	// serialize options
	b, _ := json.Marshal(ver.Options)
	dbVer.Options.UnmarshalJSON(b)

	return
}
