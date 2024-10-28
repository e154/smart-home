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

package controllers

import (
	"github.com/e154/smart-home/internal/api/dto"
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/pkg/common"
	"github.com/labstack/echo/v4"
)

// ControllerMetric ...
type ControllerMetric struct {
	*ControllerCommon
}

// NewControllerMetric ...
func NewControllerMetric(common *ControllerCommon) *ControllerMetric {
	return &ControllerMetric{
		ControllerCommon: common,
	}
}

// MetricServiceGetMetric ...
func (c ControllerMetric) MetricServiceGetMetric(ctx echo.Context, params stub.MetricServiceGetMetricParams) error {

	var metricRange *string
	if params.Range != nil {
		metricRange = common.String(string(*params.Range))
	}
	metric, err := c.endpoint.Metric.GetByIdWithData(ctx.Request().Context(), params.StartDate, params.EndDate, params.Id, metricRange)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.Metric(metric)))
}
