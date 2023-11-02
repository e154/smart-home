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
	"github.com/e154/smart-home/api/stub"
	"github.com/labstack/echo/v4"
)

// ControllerDashboardCardItem ...
type ControllerDashboardCardItem struct {
	*ControllerCommon
}

// NewControllerDashboardCardItem ...
func NewControllerDashboardCardItem(common *ControllerCommon) *ControllerDashboardCardItem {
	return &ControllerDashboardCardItem{
		ControllerCommon: common,
	}
}

// AddDashboardCardItem ...
func (c ControllerDashboardCardItem) DashboardCardItemServiceAddDashboardCardItem(ctx echo.Context, _ stub.DashboardCardItemServiceAddDashboardCardItemParams) error {

	obj := &stub.ApiNewDashboardCardItemRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	cardItem, err := c.endpoint.DashboardCardItem.Add(ctx.Request().Context(), c.dto.DashboardCardItem.AddDashboardCardItem(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.DashboardCardItem.ToDashboardCardItem(cardItem)))
}

// UpdateDashboardCardItem ...
func (c ControllerDashboardCardItem) DashboardCardItemServiceUpdateDashboardCardItem(ctx echo.Context, id int64, _ stub.DashboardCardItemServiceUpdateDashboardCardItemParams) error {

	obj := &stub.DashboardCardItemServiceUpdateDashboardCardItemJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	cardItem, err := c.endpoint.DashboardCardItem.Update(ctx.Request().Context(), c.dto.DashboardCardItem.UpdateDashboardCardItem(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.DashboardCardItem.ToDashboardCardItem(cardItem)))
}

// GetDashboardCardItemById ...
func (c ControllerDashboardCardItem) DashboardCardItemServiceGetDashboardCardItemById(ctx echo.Context, id int64) error {

	cardItem, err := c.endpoint.DashboardCardItem.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.DashboardCardItem.ToDashboardCardItem(cardItem)))
}

// GetDashboardCardItemList ...
func (c ControllerDashboardCardItem) DashboardCardItemServiceGetDashboardCardItemList(ctx echo.Context, params stub.DashboardCardItemServiceGetDashboardCardItemListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.DashboardCardItem.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.DashboardCardItem.ToListResult(items), total, pagination))
}

// DeleteDashboardCardItem ...
func (c ControllerDashboardCardItem) DashboardCardItemServiceDeleteDashboardCardItem(ctx echo.Context, id int64) error {

	if err := c.endpoint.DashboardCardItem.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}
