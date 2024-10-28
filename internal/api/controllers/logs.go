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
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/labstack/echo/v4"
)

// ControllerLogs ...
type ControllerLogs struct {
	*ControllerCommon
}

// NewControllerLogs ...
func NewControllerLogs(common *ControllerCommon) *ControllerLogs {
	return &ControllerLogs{
		ControllerCommon: common,
	}
}

// GetLogList ...
func (c ControllerLogs) LogServiceGetLogList(ctx echo.Context, params stub.LogServiceGetLogListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Log.GetList(ctx.Request().Context(), pagination, params.Query, params.StartDate, params.EndDate)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Log.ToListResult(items), total, pagination))
}
