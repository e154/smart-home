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

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
)

// MetricBuckets ...
type MetricBuckets struct {
	Db        *gorm.DB
	Timescale bool
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
func (n *MetricBuckets) SimpleListWithSoftRange(ctx context.Context, from, to time.Time, metricId int64) (list []*MetricBucket, err error) {
	var num int64 = 1
	t := from.Sub(to).Seconds()
	if t > 3600 {
		num = int64(t / 3600)
	}

	list = make([]*MetricBucket, 0)

	if n.Timescale {
		q := `SELECT time_bucket('%d seconds', time) AS time, value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time between ? and ?
GROUP BY time, value
ORDER BY time ASC
LIMIT 3600`
		q = fmt.Sprintf(q, num)
		if err = n.Db.WithContext(ctx).Raw(q, metricId, from.Format(time.RFC3339), to.Format(time.RFC3339)).Scan(&list).Error; err != nil {
			err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
		}
		return
	}

	q := `SELECT  time, value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time between ? and ?
ORDER BY time ASC;`

	if err = n.Db.WithContext(ctx).Raw(q, metricId, from.Format(time.RFC3339), to.Format(time.RFC3339)).Scan(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
	}

	return
}

// SimpleListByRangeType ...
func (n *MetricBuckets) SimpleListByRangeType(ctx context.Context, metricId int64, metricRange common.MetricRange) (list []*MetricBucket, err error) {

	var interval string
	var r string

	switch metricRange {
	case common.MetricRange6H:
		interval = "6 seconds"
		r = "6 hour"
	case common.MetricRange12H:
		interval = "12 seconds"
		r = "12 hour"
	case common.MetricRange24H:
		interval = "15 seconds"
		r = "24 hour"
	case common.MetricRange7d:
		interval = "2 minutes"
		r = "7 days"
	case common.MetricRange30d, common.MetricRange1m:
		interval = "7 minutes"
		r = "30 days"
	default:
		err = errors.Wrap(apperr.ErrMetricBucketGet, fmt.Sprintf("unknown filter %s", metricRange))
		return
	}

	list = make([]*MetricBucket, 0)

	if n.Timescale {
		var q = `SELECT time_bucket('%s', time) AS time, value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '%s'
GROUP BY time, value
ORDER BY time ASC
LIMIT 3600`
		if err = n.Db.WithContext(ctx).Raw(fmt.Sprintf(q, interval, r), metricId).Scan(&list).Error; err != nil {
			err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
		}
		return
	}

	q := `SELECT  time, value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '%s'
ORDER BY time ASC;`

	if err = n.Db.WithContext(ctx).Raw(fmt.Sprintf(q, r), metricId).Scan(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
	}

	return
}

// Simple24HPreview ...
func (n *MetricBuckets) Simple24HPreview(ctx context.Context, metricId int64) (list []*MetricBucket, err error) {
	list = make([]*MetricBucket, 0)

	if n.Timescale {
		q := `SELECT time_bucket('250 seconds', time) AS time, value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '24 hour'
GROUP BY time, value
ORDER BY time ASC
LIMIT 345`

		if err = n.Db.WithContext(ctx).Raw(q, metricId).Scan(&list).Error; err != nil {
			err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
		}
		return
	}

	q := `SELECT  time, value
FROM metric_bucket c
WHERE c.metric_id = ? and c.time > NOW() - interval '24 hour'
ORDER BY time ASC`

	if err = n.Db.WithContext(ctx).Raw(q, metricId).Scan(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketGet, err.Error())
	}
	return
}

// DeleteOldest ...
func (n *MetricBuckets) DeleteOldest(ctx context.Context, days int) (err error) {
	bucket := &MetricBucket{}
	if err = n.Db.WithContext(ctx).Last(&bucket).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogDelete, err.Error())
		return
	}
	err = n.Db.WithContext(ctx).Delete(&MetricBucket{},
		fmt.Sprintf(`time < CAST('%s' AS DATE) - interval '%d days'`,
			bucket.Time.UTC().Format("2006-01-02 15:04:05"), days)).Error
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

// AddMultiple ...
func (n *MetricBuckets) AddMultiple(ctx context.Context, buckets []*MetricBucket) (err error) {
	if err = n.Db.WithContext(ctx).Create(&buckets).Error; err != nil {
		err = errors.Wrap(apperr.ErrMetricBucketAdd, err.Error())
	}
	return
}
