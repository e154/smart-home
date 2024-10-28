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
	"github.com/e154/smart-home/pkg/common"
	"github.com/labstack/echo/v4"
)

// ControllerEntityStorage ...
type ControllerEntityStorage struct {
	*ControllerCommon
}

// NewControllerEntityStorage ...
func NewControllerEntityStorage(common *ControllerCommon) *ControllerEntityStorage {
	return &ControllerEntityStorage{
		ControllerCommon: common,
	}
}

// GetEntityStorageList ...
func (c ControllerEntityStorage) EntityStorageServiceGetEntityStorageList(ctx echo.Context, params stub.EntityStorageServiceGetEntityStorageListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)

	var entityIds []common.EntityId
	if params.EntityId != nil {
		for _, item := range *params.EntityId {
			entityIds = append(entityIds, common.EntityId(item))
		}
	}

	items, total, err := c.endpoint.EntityStorage.GetList(ctx.Request().Context(), entityIds, pagination, params.StartDate, params.EndDate)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.EntityStorage.ToListResult(items), total, pagination))
}
