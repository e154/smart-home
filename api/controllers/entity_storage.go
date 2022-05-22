// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"context"
	"github.com/e154/smart-home/common"

	"github.com/e154/smart-home/api/stub/api"
)

// ControllerEntityStorage ...
type ControllerEntityStorage struct {
	*ControllerCommon
}

// NewControllerEntityStorage ...
func NewControllerEntityStorage(common *ControllerCommon) ControllerEntityStorage {
	return ControllerEntityStorage{
		ControllerCommon: common,
	}
}

// GetEntityStorageList ...
func (c ControllerEntityStorage) GetEntityStorageList(ctx context.Context, req *api.GetEntityStorageRequest) (*api.GetEntityStorageResult, error) {

	entity, err := c.endpoint.Entity.GetById(ctx, common.EntityId(req.EntityId))
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.EntityStorage.GetList(ctx, common.EntityId(req.EntityId), pagination, req.StartDate, req.EndDate)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.EntityStorage.List(items, uint64(total), pagination, entity), nil
}
