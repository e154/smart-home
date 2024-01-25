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
	"github.com/labstack/echo/v4"

	"github.com/e154/smart-home/api/stub"
)

// ControllerAction ...
type ControllerAction struct {
	*ControllerCommon
}

// NewControllerAction ...
func NewControllerAction(common *ControllerCommon) *ControllerAction {
	return &ControllerAction{
		ControllerCommon: common,
	}
}

// AddAction ...
func (c ControllerAction) ActionServiceAddAction(ctx echo.Context, _ stub.ActionServiceAddActionParams) error {

	obj := &stub.ApiNewActionRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	action, err := c.endpoint.Action.Add(ctx.Request().Context(), c.dto.Action.Add(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Action.ToAction(action)))
}

// UpdateAction ...
func (c ControllerAction) ActionServiceUpdateAction(ctx echo.Context, id int64, _ stub.ActionServiceUpdateActionParams) error {

	obj := &stub.ActionServiceUpdateActionJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	action, err := c.endpoint.Action.Update(ctx.Request().Context(), c.dto.Action.Update(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Action.ToAction(action)))
}

// GetActionById ...
func (c ControllerAction) ActionServiceGetActionById(ctx echo.Context, id int64) error {

	action, err := c.endpoint.Action.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Action.ToAction(action)))
}

// GetActionList ...
func (c ControllerAction) ActionServiceGetActionList(ctx echo.Context, params stub.ActionServiceGetActionListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Action.GetList(ctx.Request().Context(), pagination, params.Ids)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Action.ToListResult(items), total, pagination))
}

// DeleteAction ...
func (c ControllerAction) ActionServiceDeleteAction(ctx echo.Context, id int64) error {

	if err := c.endpoint.Action.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// SearchAction ...
func (c ControllerAction) ActionServiceSearchAction(ctx echo.Context, params stub.ActionServiceSearchActionParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Action.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Action.ToSearchResult(items))
}
