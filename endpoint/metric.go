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

package endpoint

import (
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"time"
)

// MetricEndpoint ...
type MetricEndpoint struct {
	*CommonEndpoint
}

// NewMetricEndpoint ...
func NewMetricEndpoint(common *CommonEndpoint) *MetricEndpoint {
	return &MetricEndpoint{
		CommonEndpoint: common,
	}
}

// GetList ...
// Add ...
func (n *MetricEndpoint) Add(params m.Metric) (result m.Metric, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.Metric.Add(params); err != nil {
		return
	}

	if result, err = n.adaptors.Metric.GetById(id); err != nil {
		return
	}

	return
}

// GetById ...
func (n *MetricEndpoint) GetById(id int64, from, to *time.Time, metricRange *string) (result m.Metric, err error) {
	result, err = n.adaptors.Metric.GetByIdWithData(id, from, to, metricRange)
	return
}

// Update ...
func (n *MetricEndpoint) Update(params m.Metric) (result m.Metric, errs []*validation.Error, err error) {

	var metric m.Metric
	if metric, err = n.adaptors.Metric.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&metric, &params, common.JsonEngine)

	// validation
	_, errs = metric.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Metric.Update(metric); err != nil {
		return
	}

	if metric, err = n.adaptors.Metric.GetById(metric.Id); err != nil {
		return
	}

	return
}

// GetList ...
func (n *MetricEndpoint) GetList(limit, offset int64, order, sortBy string) (result []m.Metric, total int64, err error) {
	result, total, err = n.adaptors.Metric.List(limit, offset, order, sortBy)
	return
}

// Delete ...
func (n *MetricEndpoint) Delete(id int64) (err error) {

	if id == 0 {
		err = errors.New("metric id is null")
		return
	}

	var metric m.Metric
	if metric, err = n.adaptors.Metric.GetById(id); err != nil {
		return
	}

	err = n.adaptors.Metric.Delete(metric.Id)

	return
}

// Search ...
func (n *MetricEndpoint) Search(query string, limit, offset int) (result []m.Metric, total int64, err error) {
	result, total, err = n.adaptors.Metric.Search(query, limit, offset)
	return
}

// AddData ...
func (n *MetricEndpoint) AddData(metricId int64, value []byte) (result []m.Metric, total int64, err error) {
	data := m.MetricDataItem{
		Value:    value,
		MetricId: metricId,
		Time:     time.Now(),
	}
	err = n.adaptors.MetricBucket.Add(data)
	return
}
