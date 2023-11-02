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

// ControllerCondition ...
type ControllerCondition struct {
	*ControllerCommon
}

// NewControllerCondition ...
func NewControllerCondition(common *ControllerCommon) *ControllerCondition {
	return &ControllerCondition{
		ControllerCommon: common,
	}
}

// AddCondition ...
func (c ControllerCondition) ConditionServiceAddCondition(ctx echo.Context, _ stub.ConditionServiceAddConditionParams) error {

	obj := &stub.ApiNewConditionRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	condition, err := c.endpoint.Condition.Add(ctx.Request().Context(), c.dto.Condition.AddCondition(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, dto.ToCondition(condition)))
}

// UpdateCondition ...
func (c ControllerCondition) ConditionServiceUpdateCondition(ctx echo.Context, id int64, _ stub.ConditionServiceUpdateConditionParams) error {

	obj := &stub.ConditionServiceUpdateConditionJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	condition, err := c.endpoint.Condition.Update(ctx.Request().Context(), c.dto.Condition.UpdateCondition(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.ToCondition(condition)))
}

// GetConditionById ...
func (c ControllerCondition) ConditionServiceGetConditionById(ctx echo.Context, id int64, _ stub.ConditionServiceGetConditionByIdParams) error {

	condition, err := c.endpoint.Condition.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.ToCondition(condition)))
}

// GetConditionList ...
func (c ControllerCondition) ConditionServiceGetConditionList(ctx echo.Context, params stub.ConditionServiceGetConditionListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Condition.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Condition.ToListResult(items), total, pagination))
}

// DeleteCondition ...
func (c ControllerCondition) ConditionServiceDeleteCondition(ctx echo.Context, id int64) error {

	if err := c.endpoint.Condition.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// SearchCondition ...
func (c ControllerCondition) ConditionServiceSearchCondition(ctx echo.Context, params stub.ConditionServiceSearchConditionParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Condition.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Condition.ToSearchResult(items))
}
