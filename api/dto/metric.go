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
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Metric(metric *m.Metric) (object *api.Metric) {
	var options = &api.MetricOption{}
	for _, item := range metric.Options.Items {
		options.Items = append(options.Items, &api.MetricOptionItem{
			Name:        item.Name,
			Description: item.Description,
			Color:       item.Color,
			Translate:   item.Translate,
			Label:       item.Label,
		})
	}

	var data = make([]*api.MetricOptionData, 0, len(metric.Data))
	for _, item := range metric.Data {
		data = append(data, &api.MetricOptionData{
			Value:    item.Value,
			MetricId: item.MetricId,
			Time:     timestamppb.New(item.Time),
		})
	}

	object = &api.Metric{
		Id:          metric.Id,
		Name:        metric.Name,
		Description: metric.Description,
		Options:     options,
		Data:        data,
		Type:        string(metric.Type),
		Ranges:      metric.Ranges,
		CreatedAt:   timestamppb.New(metric.CreatedAt),
		UpdatedAt:   timestamppb.New(metric.UpdatedAt),
	}

	return
}

func Metrics(metrics []*m.Metric) (objects []*api.Metric) {
	objects = make([]*api.Metric, 0, len(metrics))
	for _, metric := range metrics {
		objects = append(objects, Metric(metric))
	}
	return
}

func AddMetric(objects []*api.Metric) (metrics []*m.Metric) {
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
