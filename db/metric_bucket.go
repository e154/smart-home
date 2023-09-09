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
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// MetricBuckets ...
type MetricBuckets struct {
	Db *gorm.DB
}

// MetricBucket ...
type MetricBucket struct {
	Value    json.RawMessage `gorm:"type:jsonb;not null"`
	Metric   *Metric
	MetricId int64
	Time     time.Time
}

// TableName ...
func (d *MetricBucket) TableName() string {
	return "metric_bucket"
}

// Add ...
func (n MetricBuckets) Add(ctx context.Context, metric *MetricBucket) (err error) {
	if err = n.Db.WithContext(ctx).Create(&metric).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketAdd, err.Error())
	}
	return
}

// SimpleListWithSoftRange ...
func (n *MetricBuckets) SimpleListWithSoftRange(ctx context.Context, from, to time.Time, metricId int64, optionItems []string) (list []*MetricBucket, err error) {
	var num int64 = 1
	t := from.Sub(to).Seconds()
	if t > 3600 {
		num = int64(t / 3600)
	}

	list = make([]*MetricBucket, 0)
	var str string
	for i, item := range optionItems {
		str += fmt.Sprintf(" '%s', trunc(avg((value ->> '%[1]s')::numeric), 2)", item)
		if i+1 < len(optionItems) {
			str += ","
		}
	}

	q := `SELECT TIMESTAMP WITH TIME ZONE 'epoch' +
       INTERVAL '1 second' * round(extract('epoch' from time) / %[1]d) * %[1]d as time,  json_build_object(
` + str + `
    ) as value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time between ? and ?
GROUP BY round(extract('epoch' from c.time) / %[1]d)
order by time asc
LIMIT 3600`
	q = fmt.Sprintf(q, num)
	if err = n.Db.WithContext(ctx).Raw(q, metricId, from.Format(time.RFC3339), to.Format(time.RFC3339)).Scan(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
	}
	return
}

// SimpleListByRangeType ...
func (n *MetricBuckets) SimpleListByRangeType(ctx context.Context, metricId int64, metricRange common.MetricRange, optionItems []string) (list []*MetricBucket, err error) {
	var str string
	for i, item := range optionItems {
		str += fmt.Sprintf(" '%s', trunc(avg((value ->> '%[1]s')::numeric), 2)", item)
		if i+1 < len(optionItems) {
			str += ","
		}
	}

	var q string
	switch metricRange {
	case common.MetricRange6H:
		q = `SELECT TIMESTAMP WITH TIME ZONE 'epoch' +
       INTERVAL '1 second' * round(extract('epoch' from time) / 6) * 6 as time,  json_build_object(
` + str + `
    ) as value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '6 hour'
GROUP BY round(extract('epoch' from c.time) / 6)
order by time asc
LIMIT 3600`
	case common.MetricRange12H:
		q = `SELECT TIMESTAMP WITH TIME ZONE 'epoch' +
       INTERVAL '1 second' * round(extract('epoch' from time) / 12) * 12 as time,  json_build_object(
` + str + `
    ) as value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '12 hour'
GROUP BY round(extract('epoch' from c.time) / 12)
order by time asc
LIMIT 3600`
	case common.MetricRange24H:
		q = `SELECT TIMESTAMP WITH TIME ZONE 'epoch' +
       INTERVAL '1 second' * round(extract('epoch' from time) / 24) * 24 as time,  json_build_object(
` + str + `
    ) as value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '24 hour'
GROUP BY round(extract('epoch' from c.time) / 24)
order by time asc
LIMIT 3600`
	case common.MetricRange7d:
		q = `SELECT TIMESTAMP WITH TIME ZONE 'epoch' +
       INTERVAL '1 second' * round(extract('epoch' from time) / 168) * 168 as time,  json_build_object(
` + str + `
    ) as value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '7 days'
GROUP BY round(extract('epoch' from c.time) / 168)
order by time asc
LIMIT 3600`
	case common.MetricRange30d, common.MetricRange1m:
		q = `SELECT TIMESTAMP WITH TIME ZONE 'epoch' +
       INTERVAL '1 second' * round(extract('epoch' from time) / 720) * 720 as time,  json_build_object(
` + str + `
    ) as value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '1 month'
GROUP BY round(extract('epoch' from c.time) / 720)
order by time asc
LIMIT 3600`
	default:
		q = `SELECT b.*
from metric_bucket b
where b.metric_id = ? and b.time > NOW() - interval '1 hour'
limit 3600`
	}
	list = make([]*MetricBucket, 0)
	if err = n.Db.WithContext(ctx).Raw(q, metricId).Scan(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
	}

	return
}

// Simple24HPreview ...
func (n *MetricBuckets) Simple24HPreview(ctx context.Context, metricId int64, optionItems []string) (list []*MetricBucket, err error) {
	list = make([]*MetricBucket, 0)
	var str string
	for i, item := range optionItems {
		str += fmt.Sprintf(" '%s', trunc(avg((value ->> '%[1]s')::numeric), 2)", item)
		if i+1 < len(optionItems) {
			str += ","
		}
	}
	q := `SELECT TIMESTAMP WITH TIME ZONE 'epoch' +
       INTERVAL '1 second' * round(extract('epoch' from time) / 250) * 250 as time,  json_build_object(
` + str + `
    ) as value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '1 day'
GROUP BY round(extract('epoch' from c.time) / 250)
order by time asc
LIMIT 345`

	if err = n.Db.WithContext(ctx).Raw(q, metricId).Scan(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
	}
	return
}

// DeleteOldest ...
func (n *MetricBuckets) DeleteOldest(ctx context.Context, days int) (err error) {
	err = n.Db.WithContext(ctx).Delete(&MetricBucket{}, fmt.Sprintf(`time < now() - interval '%d days'`, days)).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketDelete, err.Error())
	}
	return
}

// DeleteByMetricId ...
func (n MetricBuckets) DeleteByMetricId(ctx context.Context, metricId int64) (err error) {
	if err = n.Db.WithContext(ctx).Delete(&MetricBucket{}, "metric_id = ?", metricId).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketDelete, err.Error())
	}
	return
}

// DeleteById ...
func (n MetricBuckets) DeleteById(ctx context.Context, id int64) (err error) {
	if err = n.Db.WithContext(ctx).Delete(&MetricBucket{}, "id = ?", id).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketDelete, err.Error())
	}
	return
}

// CreateHypertable ...
func (n MetricBuckets) CreateHypertable(ctx context.Context, ) (err error) {
	if err = n.Db.WithContext(ctx).Raw(`SELECT create_hypertable('metric_bucket', 'time', migrate_data=>'TRUE')`).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketAdd, err.Error())
	}
	return
}

// AddMultiple ...
func (n *MetricBuckets) AddMultiple(ctx context.Context, buckets []*MetricBucket) (err error) {
	if err = n.Db.WithContext(ctx).Create(&buckets).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketAdd, err.Error())
	}
	return
}
