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

// ControllerDashboardCardItem ...
type ControllerDashboardCardItem struct {
	*ControllerCommon
}

// NewControllerDashboardCardItem ...
func NewControllerDashboardCardItem(common *ControllerCommon) ControllerDashboardCardItem {
	return ControllerDashboardCardItem{
		ControllerCommon: common,
	}
}

// AddDashboardCardItem ...
func (c ControllerDashboardCardItem) AddDashboardCardItem(ctx context.Context, req *api.NewDashboardCardItemRequest) (*api.DashboardCardItem, error) {

	board := c.dto.DashboardCardItem.AddDashboardCardItem(req)

	board, errs, err := c.endpoint.DashboardCardItem.Add(ctx, board)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.DashboardCardItem.ToDashboardCardItem(board), nil
}

// UpdateDashboardCardItem ...
func (c ControllerDashboardCardItem) UpdateDashboardCardItem(ctx context.Context, req *api.UpdateDashboardCardItemRequest) (*api.DashboardCardItem, error) {

	board := c.dto.DashboardCardItem.UpdateDashboardCardItem(req)

	board, errs, err := c.endpoint.DashboardCardItem.Update(ctx, board)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.DashboardCardItem.ToDashboardCardItem(board), nil
}

// GetDashboardCardItemById ...
func (c ControllerDashboardCardItem) GetDashboardCardItemById(ctx context.Context, req *api.GetDashboardCardItemRequest) (*api.DashboardCardItem, error) {

	board, err := c.endpoint.DashboardCardItem.GetById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.DashboardCardItem.ToDashboardCardItem(board), nil
}

// GetDashboardCardItemList ...
func (c ControllerDashboardCardItem) GetDashboardCardItemList(ctx context.Context, req *api.PaginationRequest) (*api.GetDashboardCardItemListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.DashboardCardItem.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.DashboardCardItem.ToListResult(items, uint64(total), pagination), nil
}

// DeleteDashboardCardItem ...
func (c ControllerDashboardCardItem) DeleteDashboardCardItem(ctx context.Context, req *api.DeleteDashboardCardItemRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.DashboardCardItem.Delete(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}
