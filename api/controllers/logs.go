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
	"context"

	"github.com/e154/smart-home/api/stub/api"
)

// ControllerLogs ...
type ControllerLogs struct {
	*ControllerCommon
}

// NewControllerLogs ...
func NewControllerLogs(common *ControllerCommon) ControllerLogs {
	return ControllerLogs{
		ControllerCommon: common,
	}
}

// GetLogList ...
func (c ControllerLogs) GetLogList(ctx context.Context, req *api.LogPaginationRequest) (*api.GetLogListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.Log.GetList(ctx, pagination, req.Query, req.StartDate, req.EndDate)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Log.ToListResult(items, uint64(total), pagination), nil
}
