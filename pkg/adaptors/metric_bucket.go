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
	"time"

	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

// MetricBucketRepo ...
type MetricBucketRepo interface {
	Add(ctx context.Context, ver *m.MetricDataItem) error
	AddMultiple(ctx context.Context, items []*m.MetricDataItem) (err error)
	List(ctx context.Context, from, to *time.Time, metricId int64, optionItems []string, metricRange *common.MetricRange) (list []*m.MetricDataItem, err error)
	DeleteOldest(ctx context.Context, days int) (err error)
	DeleteById(ctx context.Context, id int64) (err error)
	DeleteByMetricId(ctx context.Context, metricId int64) (err error)
}
