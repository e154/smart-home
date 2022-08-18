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

package controllers

import (
	"context"
	"time"

	"github.com/e154/smart-home/api/dto"
	"github.com/e154/smart-home/common"

	"github.com/e154/smart-home/api/stub/api"
)

// ControllerMetric ...
type ControllerMetric struct {
	*ControllerCommon
}

// NewControllerMetric ...
func NewControllerMetric(common *ControllerCommon) ControllerMetric {
	return ControllerMetric{
		ControllerCommon: common,
	}
}

// GetMetric ...
func (c ControllerMetric) GetMetric(ctx context.Context, req *api.GetMetricRequest) (*api.Metric, error) {

	var from, to *time.Time
	if req.From != nil {
		from = common.Time(req.From.AsTime())
	}
	if req.To != nil {
		to = common.Time(req.To.AsTime())
	}
	metric, err := c.endpoint.Metric.GetByIdWithData(ctx, from, to, req.Id, req.Range)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return dto.Metric(metric), nil
}
