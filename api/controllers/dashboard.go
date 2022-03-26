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

	"github.com/e154/smart-home/api/stub/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerDashboard ...
type ControllerDashboard struct {
	*ControllerCommon
}

// NewControllerDashboard ...
func NewControllerDashboard(common *ControllerCommon) ControllerDashboard {
	return ControllerDashboard{
		ControllerCommon: common,
	}
}

// AddDashboard ...
func (c ControllerDashboard) AddDashboard(ctx context.Context, req *api.NewDashboardRequest) (*api.Dashboard, error) {

	board := c.dto.Dashboard.AddDashboard(req)

	board, errs, err := c.endpoint.Dashboard.Add(ctx, board)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Dashboard.ToDashboard(board), nil
}

// UpdateDashboard ...
func (c ControllerDashboard) UpdateDashboard(ctx context.Context, req *api.UpdateDashboardRequest) (*api.Dashboard, error) {

	board := c.dto.Dashboard.UpdateDashboard(req)

	board, errs, err := c.endpoint.Dashboard.Update(ctx, board)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Dashboard.ToDashboard(board), nil
}

// GetDashboardById ...
func (c ControllerDashboard) GetDashboardById(ctx context.Context, req *api.GetDashboardRequest) (*api.Dashboard, error) {

	board, err := c.endpoint.Dashboard.GetById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Dashboard.ToDashboard(board), nil
}

// GetDashboardList ...
func (c ControllerDashboard) GetDashboardList(ctx context.Context, req *api.PaginationRequest) (*api.GetDashboardListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.Dashboard.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Dashboard.ToListResult(items, uint64(total), pagination), nil
}

// DeleteDashboard ...
func (c ControllerDashboard) DeleteDashboard(ctx context.Context, req *api.DeleteDashboardRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Dashboard.Delete(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}
