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
	"github.com/labstack/echo/v4"

	"github.com/e154/smart-home/api/stub"
)

// ControllerArea ...
type ControllerArea struct {
	*ControllerCommon
}

// NewControllerArea ...
func NewControllerArea(common *ControllerCommon) *ControllerArea {
	return &ControllerArea{
		ControllerCommon: common,
	}
}

// AddArea ...
func (c ControllerArea) AreaServiceAddArea(ctx echo.Context, _ stub.AreaServiceAddAreaParams) error {

	obj := &stub.ApiNewAreaRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	area, err := c.endpoint.Area.Add(ctx.Request().Context(), c.dto.Area.AddArea(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, dto.ToArea(area)))
}

// UpdateArea ...
func (c ControllerArea) AreaServiceUpdateArea(ctx echo.Context, id int64, _ stub.AreaServiceUpdateAreaParams) error {

	obj := &stub.AreaServiceUpdateAreaJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	area, err := c.endpoint.Area.Update(ctx.Request().Context(), c.dto.Area.UpdateArea(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.ToArea(area)))
}

// GetAreaById ...
func (c ControllerArea) AreaServiceGetAreaById(ctx echo.Context, id int64) error {

	area, err := c.endpoint.Area.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.ToArea(area)))
}

// GetAreaList ...
func (c ControllerArea) AreaServiceGetAreaList(ctx echo.Context, params stub.AreaServiceGetAreaListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Area.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Area.ToListResult(items), total, pagination))
}

// DeleteArea ...
func (c ControllerArea) AreaServiceDeleteArea(ctx echo.Context, id int64) error {

	if err := c.endpoint.Area.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// SearchArea ...
func (c ControllerArea) AreaServiceSearchArea(ctx echo.Context, params stub.AreaServiceSearchAreaParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Area.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Area.ToSearchResult(items))
}
