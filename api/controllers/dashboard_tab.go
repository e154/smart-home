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

// ControllerDashboardTab ...
type ControllerDashboardTab struct {
	*ControllerCommon
}

// NewControllerDashboardTab ...
func NewControllerDashboardTab(common *ControllerCommon) ControllerDashboardTab {
	return ControllerDashboardTab{
		ControllerCommon: common,
	}
}

// AddDashboardTab ...
func (c ControllerDashboardTab) AddDashboardTab(ctx context.Context, req *api.NewDashboardTabRequest) (*api.DashboardTab, error) {

	board := c.dto.DashboardTab.AddDashboardTab(req)

	board, errs, err := c.endpoint.DashboardTab.Add(ctx, board)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.DashboardTab.ToDashboardTab(board), nil
}

// UpdateDashboardTab ...
func (c ControllerDashboardTab) UpdateDashboardTab(ctx context.Context, req *api.UpdateDashboardTabRequest) (*api.DashboardTab, error) {

	board := c.dto.DashboardTab.UpdateDashboardTab(req)

	board, errs, err := c.endpoint.DashboardTab.Update(ctx, board)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.DashboardTab.ToDashboardTab(board), nil
}

// GetDashboardTabById ...
func (c ControllerDashboardTab) GetDashboardTabById(ctx context.Context, req *api.GetDashboardTabRequest) (*api.DashboardTab, error) {

	board, err := c.endpoint.DashboardTab.GetById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.DashboardTab.ToDashboardTab(board), nil
}

// GetDashboardTabList ...
func (c ControllerDashboardTab) GetDashboardTabList(ctx context.Context, req *api.PaginationRequest) (*api.GetDashboardTabListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.DashboardTab.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.DashboardTab.ToListResult(items, uint64(total), pagination), nil
}

// DeleteDashboardTab ...
func (c ControllerDashboardTab) DeleteDashboardTab(ctx context.Context, req *api.DeleteDashboardTabRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.DashboardTab.Delete(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}
