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
	"github.com/e154/smart-home/internal/api/stub"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/labstack/echo/v4"
)

// ControllerTag ...
type ControllerTag struct {
	*ControllerCommon
}

// NewControllerTag ...
func NewControllerTag(common *ControllerCommon) *ControllerTag {
	return &ControllerTag{
		ControllerCommon: common,
	}
}

// TagServiceGetTagById ...
func (c ControllerTag) TagServiceGetTagById(ctx echo.Context, id int64) error {

	tag, err := c.endpoint.Tag.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, tag))
}

// TagServiceUpdateTagById ...
func (c ControllerTag) TagServiceUpdateTagById(ctx echo.Context, id int64, _ stub.TagServiceUpdateTagByIdParams) error {

	obj := &stub.TagServiceUpdateTagByIdJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	tag, err := c.endpoint.Tag.Update(ctx.Request().Context(), &m.Tag{Id: id, Name: obj.Name})
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, tag))
}

// TagServiceDeleteTagById ...
func (c ControllerTag) TagServiceDeleteTagById(ctx echo.Context, id int64) error {

	if err := c.endpoint.Tag.DeleteTagById(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// TagServiceGetTagList ...
func (c ControllerTag) TagServiceGetTagList(ctx echo.Context, params stub.TagServiceGetTagListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Tag.GetList(ctx.Request().Context(), pagination, params.Query, params.Tags)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Tag.ToListResult(items), total, pagination))
}

// TagServiceSearchTag ...
func (c ControllerTag) TagServiceSearchTag(ctx echo.Context, params stub.TagServiceSearchTagParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Tag.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Tag.ToSearchResult(items))
}
