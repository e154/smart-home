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
	"github.com/e154/smart-home/api/dto"
	"github.com/e154/smart-home/api/stub"
	"github.com/labstack/echo/v4"
)

// ControllerDashboardTab ...
type ControllerDashboardTab struct {
	*ControllerCommon
}

// NewControllerDashboardTab ...
func NewControllerDashboardTab(common *ControllerCommon) *ControllerDashboardTab {
	return &ControllerDashboardTab{
		ControllerCommon: common,
	}
}

// AddDashboardTab ...
func (c ControllerDashboardTab) DashboardTabServiceAddDashboardTab(ctx echo.Context, _ stub.DashboardTabServiceAddDashboardTabParams) error {

	obj := &stub.ApiNewDashboardTabRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	board, err := c.endpoint.DashboardTab.Add(ctx.Request().Context(), c.dto.DashboardTab.AddDashboardTab(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.DashboardTab.ToDashboardTab(board)))
}

// UpdateDashboardTab ...
func (c ControllerDashboardTab) DashboardTabServiceUpdateDashboardTab(ctx echo.Context, id int64, _ stub.DashboardTabServiceUpdateDashboardTabParams) error {

	obj := &stub.DashboardTabServiceUpdateDashboardTabJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	tab, err := c.endpoint.DashboardTab.Update(ctx.Request().Context(), c.dto.DashboardTab.UpdateDashboardTab(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.DashboardTab.ToDashboardTab(tab)))
}

// GetDashboardTabById ...
func (c ControllerDashboardTab) DashboardTabServiceGetDashboardTabById(ctx echo.Context, id int64) error {

	tab, err := c.endpoint.DashboardTab.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.DashboardTab.ToDashboardTab(tab)))
}

// GetDashboardTabList ...
func (c ControllerDashboardTab) DashboardTabServiceGetDashboardTabList(ctx echo.Context, params stub.DashboardTabServiceGetDashboardTabListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.DashboardTab.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.DashboardTab.ToListResult(items), total, pagination))
}

// DeleteDashboardTab ...
func (c ControllerDashboardTab) DashboardTabServiceDeleteDashboardTab(ctx echo.Context, id int64) error {

	if err := c.endpoint.DashboardTab.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// DashboardTabServiceImportDashboardTab ...
func (c ControllerDashboardTab) DashboardTabServiceImportDashboardTab(ctx echo.Context, _ stub.DashboardTabServiceImportDashboardTabParams) error {

	obj := &stub.ApiDashboardTab{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	board, err := c.endpoint.DashboardTab.Import(ctx.Request().Context(), dto.ImportDashboardTab(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.DashboardTab.ToDashboardTab(board)))
}
