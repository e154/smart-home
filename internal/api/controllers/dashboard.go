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
	"github.com/labstack/echo/v4"
)

// ControllerDashboard ...
type ControllerDashboard struct {
	*ControllerCommon
}

// NewControllerDashboard ...
func NewControllerDashboard(common *ControllerCommon) *ControllerDashboard {
	return &ControllerDashboard{
		ControllerCommon: common,
	}
}

// AddDashboard ...
func (c ControllerDashboard) DashboardServiceAddDashboard(ctx echo.Context, _ stub.DashboardServiceAddDashboardParams) error {

	obj := &stub.ApiNewDashboardRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	board, err := c.endpoint.Dashboard.Add(ctx.Request().Context(), c.dto.Dashboard.AddDashboard(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Dashboard.ToDashboard(board)))
}

// UpdateDashboard ...
func (c ControllerDashboard) DashboardServiceUpdateDashboard(ctx echo.Context, id int64, _ stub.DashboardServiceUpdateDashboardParams) error {

	obj := &stub.DashboardServiceUpdateDashboardJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	board, err := c.endpoint.Dashboard.Update(ctx.Request().Context(), c.dto.Dashboard.UpdateDashboard(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Dashboard.ToDashboard(board)))
}

// GetDashboardById ...
func (c ControllerDashboard) DashboardServiceGetDashboardById(ctx echo.Context, id int64) error {

	board, err := c.endpoint.Dashboard.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Dashboard.ToDashboard(board)))
}

// GetDashboardList ...
func (c ControllerDashboard) DashboardServiceGetDashboardList(ctx echo.Context, params stub.DashboardServiceGetDashboardListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Dashboard.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Dashboard.ToListResult(items), total, pagination))
}

// DeleteDashboard ...
func (c ControllerDashboard) DashboardServiceDeleteDashboard(ctx echo.Context, id int64) error {

	if err := c.endpoint.Dashboard.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// ImportDashboard ...
func (c ControllerDashboard) DashboardServiceImportDashboard(ctx echo.Context, _ stub.DashboardServiceImportDashboardParams) error {

	obj := &stub.ApiDashboard{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	board, err := c.endpoint.Dashboard.Import(ctx.Request().Context(), dto.ImportDashboard(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Dashboard.ToDashboard(board)))
}

// SearchDashboard ...
func (c ControllerDashboard) DashboardServiceSearchDashboard(ctx echo.Context, params stub.DashboardServiceSearchDashboardParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Dashboard.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Dashboard.ToSearchResult(items))
}
