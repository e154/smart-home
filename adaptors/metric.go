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
	"github.com/e154/smart-home/system/cache"
	"github.com/e154/smart-home/system/orm"
	"github.com/jinzhu/gorm"
)

// Metric ...
type Metric struct {
	table *db.Metrics
	db    *gorm.DB
	c     cache.Cache
	orm   *orm.Orm
}

// GetMetricAdaptor ...
func GetMetricAdaptor(d *gorm.DB, orm *orm.Orm) *Metric {
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

// GetByMapDeviceId ...
func (n *Metric) GetByMapDeviceId(mapDeviceId int64, name string) (metric m.Metric, err error) {
	var dbList []db.Metric
	if dbList, err = n.table.GetByMapDeviceId(mapDeviceId, name); err != nil {
		return
	}

	if len(dbList) == 0 {
		err = ErrRecordNotFound
		return
	}

	metric = n.fromDb(dbList[0])

	return
}

// Update ...
func (n *Metric) Update(ver m.Metric) error {
	return n.table.Update(n.toDb(ver))
}

// Delete ...
func (n *Metric) Delete(verId int64) error {
	return n.table.Delete(verId)
}

func (n *Metric) fromDb(dbVer db.Metric) (ver m.Metric) {
	ver = m.Metric{
		Id:           dbVer.Id,
		Name:         dbVer.Name,
		Description:  dbVer.Description,
		CreatedAt:    dbVer.CreatedAt,
		UpdatedAt:    dbVer.UpdatedAt,
		MapDeviceId:  dbVer.MapDeviceId,
		Translations: dbVer.Translations,
	}

	return
}

func (n *Metric) toDb(ver m.Metric) (dbVer db.Metric) {
	dbVer = db.Metric{
		Id:           ver.Id,
		Name:         ver.Name,
		Description:  ver.Description,
		MapDeviceId:  ver.MapDeviceId,
		Translations: ver.Translations,
	}

	return
}
