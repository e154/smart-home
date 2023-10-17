// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package dto

import (
	stub "github.com/e154/smart-home/api/stub"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

func Metric(metric *m.Metric) (object *stub.ApiMetric) {
	var options = &stub.ApiMetricOption{}
	for _, item := range metric.Options.Items {
		options.Items = append(options.Items, stub.ApiMetricOptionItem{
			Name:        item.Name,
			Description: item.Description,
			Color:       item.Color,
			Translate:   item.Translate,
			Label:       item.Label,
		})
	}

	var data = make([]stub.ApiMetricOptionData, 0, len(metric.Data))
	for _, item := range metric.Data {
		data = append(data, stub.ApiMetricOptionData{
			Value:    item.Value,
			MetricId: item.MetricId,
			Time:     item.Time,
		})
	}

	object = &stub.ApiMetric{
		Id:          metric.Id,
		Name:        metric.Name,
		Description: metric.Description,
		Options:     options,
		Data:        data,
		Type:        string(metric.Type),
		Ranges:      metric.Ranges,
		CreatedAt:   metric.CreatedAt,
		UpdatedAt:   metric.UpdatedAt,
	}

	return
}

func Metrics(metrics []*m.Metric) (objects []stub.ApiMetric) {
	objects = make([]stub.ApiMetric, 0, len(metrics))
	for _, metric := range metrics {
		objects = append(objects, *Metric(metric))
	}
	return
}

func AddMetric(objects []stub.ApiMetric) (metrics []*m.Metric) {
	metrics = make([]*m.Metric, 0, len(objects))
	for _, obj := range objects {
		if obj.Options == nil {
			continue
		}
		var options m.MetricOptions
		for _, item := range obj.Options.Items {
			options.Items = append(options.Items, m.MetricOptionsItem{
				Name:        item.Name,
				Description: item.Description,
				Color:       item.Color,
				Translate:   item.Translate,
				Label:       item.Label,
			})
		}
		metrics = append(metrics, &m.Metric{
			Id:          obj.Id,
			Name:        obj.Name,
			Description: obj.Description,
			Options:     options,
			Type:        common.MetricType(obj.Type),
			Ranges:      obj.Ranges,
		})
	}
	return
}
