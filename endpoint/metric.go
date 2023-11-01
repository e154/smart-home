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

package endpoint

import (
	"context"
	"fmt"
	"github.com/e154/smart-home/common"
	"time"

	m "github.com/e154/smart-home/models"
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

// GetByIdWithData ...
func (l *MetricEndpoint) GetByIdWithData(ctx context.Context, from, to *time.Time, metricId int64, _metricRange *string) (metric *m.Metric, err error) {

	key := fmt.Sprintf("metric_%d_%v_%v_%v", metricId, from, to, _metricRange)
	if l.cache.IsExist(key) {
		v := l.cache.Get(key)
		metric = v.(*m.Metric)
		return
	}

	var metricRange *common.MetricRange
	if _metricRange != nil {
		metricRange = common.MetricRange(*_metricRange).Ptr()
	}
	if metric, err = l.adaptors.Metric.GetByIdWithData(ctx, metricId, from, to, metricRange); err != nil {
		return
	}

	l.cache.Put(key, metric, 10*time.Second)

	return
}
