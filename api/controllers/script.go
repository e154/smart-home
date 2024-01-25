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

	"github.com/e154/smart-home/api/dto"
)

// ControllerScript ...
type ControllerScript struct {
	*ControllerCommon
}

// NewControllerScript ...
func NewControllerScript(common *ControllerCommon) *ControllerScript {
	return &ControllerScript{
		ControllerCommon: common,
	}
}

// AddScript ...
func (c ControllerScript) ScriptServiceAddScript(ctx echo.Context, _ stub.ScriptServiceAddScriptParams) error {

	obj := &stub.ApiNewScriptRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	script, err := c.endpoint.Script.Add(ctx.Request().Context(), c.dto.Script.FromNewScriptRequest(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Script.GetStubScript(script)))
}

// GetScriptById ...
func (c ControllerScript) ScriptServiceGetScriptById(ctx echo.Context, id int64) error {

	script, err := c.endpoint.Script.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Script.GetStubScript(script)))
}

// UpdateScriptById ...
func (c ControllerScript) ScriptServiceUpdateScriptById(ctx echo.Context, id int64, _ stub.ScriptServiceUpdateScriptByIdParams) error {

	obj := &stub.ScriptServiceUpdateScriptByIdJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	script, err := c.endpoint.Script.Update(ctx.Request().Context(), c.dto.Script.FromUpdateScriptRequest(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Script.GetStubScript(script)))
}

// GetScriptList ...
func (c ControllerScript) ScriptServiceGetScriptList(ctx echo.Context, params stub.ScriptServiceGetScriptListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Script.GetList(ctx.Request().Context(), pagination, params.Query, params.Ids)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Script.ToListResult(items), total, pagination))
}

// SearchScript ...
func (c ControllerScript) ScriptServiceSearchScript(ctx echo.Context, params stub.ScriptServiceSearchScriptParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Script.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Script.ToSearchResult(items))
}

// DeleteScriptById ...
func (c ControllerScript) ScriptServiceDeleteScriptById(ctx echo.Context, id int64) error {

	if err := c.endpoint.Script.DeleteScriptById(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// ExecScriptById ...
func (c ControllerScript) ScriptServiceExecScriptById(ctx echo.Context, id int64) error {

	result, err := c.endpoint.Script.Execute(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, &stub.ApiExecScriptResult{Result: result}))
}

// ExecSrcScriptById ...
func (c ControllerScript) ScriptServiceExecSrcScriptById(ctx echo.Context, _ stub.ScriptServiceExecSrcScriptByIdParams) error {

	obj := &stub.ApiExecSrcScriptRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	result, err := c.endpoint.Script.ExecuteSource(ctx.Request().Context(), c.dto.Script.FromExecSrcScriptRequest(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, &stub.ApiExecScriptResult{Result: result}))
}

// CopyScriptById ...
func (c ControllerScript) ScriptServiceCopyScriptById(ctx echo.Context, id int64) error {

	script, err := c.endpoint.Script.Copy(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Script.GetStubScript(script)))
}

// GetStatistic ...
func (c ControllerScript) ScriptServiceGetStatistic(ctx echo.Context) error {

	statistic, err := c.endpoint.Script.Statistic(ctx.Request().Context())
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.GetStatistic(statistic)))
}
