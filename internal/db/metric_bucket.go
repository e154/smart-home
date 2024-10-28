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

	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/pkg/errors"
)

// MetricBuckets ...
type MetricBuckets struct {
	*Common
	Timescale bool
}

// MetricBucket ...
type MetricBucket struct {
	Value    json.RawMessage `gorm:"type:jsonb;not null"`
	Metric   *Metric
	Mins     time.Time `gorm:"->"`
	MetricId int64
	Time     time.Time
}

// TableName ...
func (d *MetricBucket) TableName() string {
	return "metric_bucket"
}

// Add ...
func (n MetricBuckets) Add(ctx context.Context, metric *MetricBucket) (err error) {
	if err = n.DB(ctx).Create(&metric).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketAdd, err.Error())
	}
	return
}

// List ...
func (n *MetricBuckets) List(ctx context.Context, metricId int64, optionItems []string, rFrom, rTo *time.Time, metricRange *pkgCommon.MetricRange) (list []*MetricBucket, err error) {

	var str string
	for i, item := range optionItems {
		str += fmt.Sprintf(" '%s', trunc(avg((value ->> '%[1]s')::numeric), 2)", item)
		if i+1 < len(optionItems) {
			str += ","
		}
	}

	var interval string
	var from time.Time
	var to = time.Now()

	if rTo != nil {
		to = pkgCommon.TimeValue(rTo)
	}

	if metricRange != nil {
		switch *metricRange {
		case pkgCommon.MetricRange6H:
			interval = "6 seconds"
			from = to.Add(-6 * time.Hour)
		case pkgCommon.MetricRange12H:
			interval = "12 seconds"
			from = to.Add(-12 * time.Hour)
		case pkgCommon.MetricRange24H:
			interval = "24 seconds"
			from = to.Add(-24 * time.Hour)
		case pkgCommon.MetricRange7d:
			interval = "150 seconds"
			from = to.Add(-24 * 7 * time.Hour)
		case pkgCommon.MetricRange30d, pkgCommon.MetricRange1m:
			interval = "7 minutes"
			from = to.Add(-24 * 30 * time.Hour)
		default:
			err = errors.Wrap(apperr.ErrMetricBucketGet, fmt.Sprintf("unknown filter %s", metricRange))
			return
		}
	}

	if rFrom != nil {
		from = pkgCommon.TimeValue(rFrom)
	}

	var num int64 = 1
	t := from.Sub(to).Seconds()
	if t > 3600 {
		num = int64(t / 3600)
	}

	if metricRange == nil {
		interval = fmt.Sprintf("%d seconds", num)
	}

	list = make([]*MetricBucket, 0)

	// c.time between ? and ?
	if n.Timescale {
		var q = `SELECT time_bucket('%s', time) AS mins, json_build_object(%s) as value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > ? and c.time < ? 
GROUP BY mins
ORDER BY mins ASC
LIMIT 3600`
		if err = n.DB(ctx).Raw(fmt.Sprintf(q, interval, str), metricId, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339)).Scan(&list).Error; err != nil {
			err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
		}

		for _, metric := range list {
			metric.Time = metric.Mins
		}

		return
	}

	q := `SELECT TIMESTAMP WITH TIME ZONE 'epoch' +
       INTERVAL '1 second' * round(extract('epoch' from time) / %[1]d) * %[1]d as time, json_build_object(%s) as value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > ? and c.time < ? 
GROUP BY round(extract('epoch' from c.time) / %[1]d)
order by time asc
LIMIT 3600`
	if err = n.DB(ctx).Raw(fmt.Sprintf(q, num, str), metricId, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339)).Scan(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
	}

	return
}

// DeleteOldest ...
func (n *MetricBuckets) DeleteOldest(ctx context.Context, days int) (err error) {
	bucket := &MetricBucket{}
	if err = n.DB(ctx).Last(&bucket).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogDelete, err.Error())
		return
	}
	err = n.DB(ctx).Delete(&MetricBucket{},
		fmt.Sprintf(`time < CAST('%s' AS DATE) - interval '%d days'`,
			bucket.Time.UTC().Format("2006-01-02 15:04:05"), days)).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketDelete, err.Error())
	}
	return
}

// DeleteByMetricId ...
func (n MetricBuckets) DeleteByMetricId(ctx context.Context, metricId int64) (err error) {
	if err = n.DB(ctx).Delete(&MetricBucket{}, "metric_id = ?", metricId).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketDelete, err.Error())
	}
	return
}

// DeleteById ...
func (n MetricBuckets) DeleteById(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Delete(&MetricBucket{}, "id = ?", id).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketDelete, err.Error())
	}
	return
}

// AddMultiple ...
func (n *MetricBuckets) AddMultiple(ctx context.Context, buckets []*MetricBucket) (err error) {
	if err = n.DB(ctx).Create(&buckets).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketAdd, err.Error())
	}
	return
}
