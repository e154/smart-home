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

package db

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common"
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
	Data        []MetricBucket
	Name        string
	Description string
	Options     json.RawMessage `gorm:"type:jsonb;not null"`
	Type        common.MetricType
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

// GetById ...
func (n Metrics) GetById(id int64) (metric Metric, err error) {
	metric = Metric{Id: id}
	err = n.Db.First(&metric).Error
	return
}

// Update ...
func (n Metrics) Update(m Metric) (err error) {
	q := map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"options":     m.Options,
		"type":        m.Type,
	}
	err = n.Db.Model(&Metric{}).Where("id = ?", m.Id).Updates(q).Error
	return
}

// Delete ...
func (n Metrics) Delete(id int64) error {
	return n.Db.Delete(&Metric{}, "id = ?", id).Error
}

// List ...
func (n *Metrics) List(limit, offset int64, orderBy, sort string) (list []Metric, total int64, err error) {

	if err = n.Db.Model(Metric{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]Metric, 0)
	q := n.Db.Model(&Metric{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	return
}

// Search ...q
func (n *Metrics) Search(query string, limit, offset int) (list []Metric, total int64, err error) {

	q := n.Db.Model(&Metric{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")


	list = make([]Metric, 0)
	err = q.Find(&list).Error

	return
}
