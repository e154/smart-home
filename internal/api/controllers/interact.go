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

// ControllerInteract ...
type ControllerInteract struct {
	*ControllerCommon
}

// NewControllerInteract ...
func NewControllerInteract(common *ControllerCommon) *ControllerInteract {
	return &ControllerInteract{
		ControllerCommon: common,
	}
}

func (c ControllerInteract) InteractServiceEntityCallAction(ctx echo.Context, _ stub.InteractServiceEntityCallActionParams) error {

	obj := &stub.ApiEntityCallActionRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	attributes := dto.AttributeFromApi(obj.Attributes)
	err := c.endpoint.Interact.EntityCallAction(ctx.Request().Context(), obj.Id, obj.Name, obj.AreaId, obj.Tags, attributes.Serialize())
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}
