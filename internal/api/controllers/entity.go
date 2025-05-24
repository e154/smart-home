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
	dto2 "github.com/e154/smart-home/internal/api/dto"
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/pkg/common"
	"github.com/labstack/echo/v4"
)

// ControllerEntity ...
type ControllerEntity struct {
	*ControllerCommon
}

// NewControllerEntity ...
func NewControllerEntity(common *ControllerCommon) *ControllerEntity {
	return &ControllerEntity{
		ControllerCommon: common,
	}
}

// AddEntity ...
func (c ControllerEntity) EntityServiceAddEntity(ctx echo.Context, _ stub.EntityServiceAddEntityParams) error {

	obj := &stub.ApiNewEntityRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	entity, err := c.endpoint.Entity.Add(ctx.Request().Context(), c.dto.Entity.AddEntity(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto2.ToEntity(entity)))
}

// ImportEntity ...
func (c ControllerEntity) EntityServiceImportEntity(ctx echo.Context, _ stub.EntityServiceImportEntityParams) error {

	obj := &stub.ApiEntity{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	err := c.endpoint.Entity.Import(ctx.Request().Context(), c.dto.Entity.ImportEntity(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// UpdateEntity ...
func (c ControllerEntity) EntityServiceUpdateEntity(ctx echo.Context, id string, _ stub.EntityServiceUpdateEntityParams) error {

	obj := &stub.EntityServiceUpdateEntityJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	entity, err := c.endpoint.Entity.Update(ctx.Request().Context(), c.dto.Entity.UpdateEntity(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto2.ToEntity(entity)))
}

// GetEntity ...
func (c ControllerEntity) EntityServiceGetEntity(ctx echo.Context, id string) error {

	entity, err := c.endpoint.Entity.GetById(ctx.Request().Context(), common.EntityId(id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto2.ToEntity(entity)))
}

// GetEntityList ...
func (c ControllerEntity) EntityServiceGetEntityList(ctx echo.Context, params stub.EntityServiceGetEntityListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Entity.List(ctx.Request().Context(), pagination, params.Query, params.Plugin, params.Area, params.Tags)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Entity.ToListResult(items), total, pagination))
}

// DeleteEntity ...
func (c ControllerEntity) EntityServiceDeleteEntity(ctx echo.Context, id string) error {

	if err := c.endpoint.Entity.Delete(ctx.Request().Context(), common.EntityId(id)); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// SearchEntity ...
func (c ControllerEntity) EntityServiceSearchEntity(ctx echo.Context, params stub.EntityServiceSearchEntityParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Entity.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Entity.ToSearchResult(items))
}

func (c ControllerEntity) EntityServiceEnabledEntity(ctx echo.Context, id string) error {
	if err := c.endpoint.Entity.Enable(ctx.Request().Context(), common.EntityId(id)); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

func (c ControllerEntity) EntityServiceDisabledEntity(ctx echo.Context, id string) error {
	if err := c.endpoint.Entity.Disable(ctx.Request().Context(), common.EntityId(id)); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// GetStatistic ...
func (c ControllerEntity) EntityServiceGetStatistic(ctx echo.Context) error {

	statistic, err := c.endpoint.Entity.Statistic(ctx.Request().Context())
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto2.GetStatistic(statistic)))
}
