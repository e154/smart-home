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

package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Metrics ...
type Metrics struct {
	Db *gorm.DB
}

// Metric ...
type Metric struct {
	Id          int64 `gorm:"primary_key"`
	Data        []*MetricBucket
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
func (n Metrics) Add(ctx context.Context, metric *Metric) (id int64, err error) {
	if err = n.Db.WithContext(ctx).Create(&metric).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricAdd, err.Error())
		return
	}
	id = metric.Id
	return
}

// GetById ...
func (n Metrics) GetById(ctx context.Context, id int64) (metric *Metric, err error) {
	metric = &Metric{Id: id}
	if err = n.Db.WithContext(ctx).First(&metric).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrMetricNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrMetricGet, err.Error())
	}
	return
}

// Update ...
func (n Metrics) Update(ctx context.Context, m *Metric) (err error) {
	q := map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"options":     m.Options,
		"type":        m.Type,
	}
	if err = n.Db.WithContext(ctx).Model(&Metric{}).Where("id = ?", m.Id).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricUpdate, err.Error())
	}
	return
}

// Delete ...
func (n Metrics) Delete(ctx context.Context, id int64) (err error) {
	if err = n.Db.WithContext(ctx).Delete(&Metric{}, "id = ?", id).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricDelete, err.Error())
	}
	return
}

// List ...
func (n *Metrics) List(ctx context.Context, limit, offset int, orderBy, sort string) (list []*Metric, total int64, err error) {

	if err = n.Db.WithContext(ctx).Model(Metric{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricList, err.Error())
		return
	}

	list = make([]*Metric, 0)
	q := n.Db.WithContext(ctx).Model(&Metric{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricList, err.Error())
	}
	return
}

// Search ...q
func (n *Metrics) Search(ctx context.Context, query string, limit, offset int) (list []*Metric, total int64, err error) {

	q := n.Db.WithContext(ctx).Model(&Metric{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Metric, 0)
	err = q.Find(&list).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrMetricSearch, err.Error())
	}
	return
}

// AddMultiple ...
func (n *Metrics) AddMultiple(ctx context.Context, metrics []*Metric) (err error) {
	if err = n.Db.WithContext(ctx).Create(&metrics).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricAdd, err.Error())
	}
	return
}
