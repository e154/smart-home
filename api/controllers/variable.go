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

// ControllerVariable ...
type ControllerVariable struct {
	*ControllerCommon
}

// NewControllerVariable ...
func NewControllerVariable(common *ControllerCommon) *ControllerVariable {
	return &ControllerVariable{
		ControllerCommon: common,
	}
}

// AddVariable ...
func (c ControllerVariable) VariableServiceAddVariable(ctx echo.Context, _ stub.VariableServiceAddVariableParams) error {

	obj := &stub.ApiNewVariableRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	variable := c.dto.Variable.AddVariable(obj)

	if err := c.endpoint.Variable.Add(ctx.Request().Context(), variable); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Variable.ToVariable(variable)))
}

// UpdateVariable ...
func (c ControllerVariable) VariableServiceUpdateVariable(ctx echo.Context, name string, _ stub.VariableServiceUpdateVariableParams) error {

	obj := &stub.VariableServiceUpdateVariableJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	variable := c.dto.Variable.UpdateVariable(obj, name)

	if err := c.endpoint.Variable.Update(ctx.Request().Context(), variable); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Variable.ToVariable(variable)))
}

// GetVariableByName ...
func (c ControllerVariable) VariableServiceGetVariableByName(ctx echo.Context, name string) error {

	variable, err := c.endpoint.Variable.GetById(ctx.Request().Context(), name)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Variable.ToVariable(variable)))
}

// GetVariableList ...
func (c ControllerVariable) VariableServiceGetVariableList(ctx echo.Context, params stub.VariableServiceGetVariableListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Variable.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Variable.ToListResult(items), total, pagination))
}

// DeleteVariable ...
func (c ControllerVariable) VariableServiceDeleteVariable(ctx echo.Context, name string) error {

	if err := c.endpoint.Variable.Delete(ctx.Request().Context(), name); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// SearchVariable ...
func (c ControllerVariable) VariableServiceSearchVariable(ctx echo.Context, params stub.VariableServiceSearchVariableParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Variable.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Variable.ToSearchResult(items))
}
