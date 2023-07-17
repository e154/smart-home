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

	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerEntity ...
type ControllerEntity struct {
	*ControllerCommon
}

// NewControllerEntity ...
func NewControllerEntity(common *ControllerCommon) ControllerEntity {
	return ControllerEntity{
		ControllerCommon: common,
	}
}

// AddEntity ...
func (c ControllerEntity) AddEntity(ctx context.Context, req *api.NewEntityRequest) (*api.Entity, error) {

	entity := c.dto.Entity.AddEntity(req)

	entity, errs, err := c.endpoint.Entity.Add(ctx, entity)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Entity.ToEntity(entity), nil
}

// ImportEntity ...
func (c ControllerEntity) ImportEntity(ctx context.Context, req *api.Entity) (*emptypb.Empty, error) {

	entity := c.dto.Entity.ImportEntity(req)

	err := c.endpoint.Entity.Import(ctx, entity)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// UpdateEntity ...
func (c ControllerEntity) UpdateEntity(ctx context.Context, req *api.UpdateEntityRequest) (*api.Entity, error) {

	entity := c.dto.Entity.UpdateEntity(req)

	entity, errs, err := c.endpoint.Entity.Update(ctx, entity)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Entity.ToEntity(entity), nil
}

// GetEntity ...
func (c ControllerEntity) GetEntity(ctx context.Context, req *api.GetEntityRequest) (*api.Entity, error) {

	entity, err := c.endpoint.Entity.GetById(ctx, common.EntityId(req.Id))
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Entity.ToEntity(entity), nil
}

// GetEntityList ...
func (c ControllerEntity) GetEntityList(ctx context.Context, req *api.PaginationRequest) (*api.GetEntityListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.Entity.List(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Entity.ToListResult(items, uint64(total), pagination), nil
}

// DeleteEntity ...
func (c ControllerEntity) DeleteEntity(ctx context.Context, req *api.DeleteEntityRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Entity.Delete(ctx, common.EntityId(req.Id)); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// SearchEntity ...
func (c ControllerEntity) SearchEntity(ctx context.Context, req *api.SearchRequest) (*api.SearchEntityResult, error) {

	search := c.Search(req.Query, req.Limit, req.Offset)
	items, _, err := c.endpoint.Entity.Search(ctx, search.Query, search.Limit, search.Offset)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Entity.ToSearchResult(items), nil
}

func (c ControllerEntity) EnabledEntity(ctx context.Context, req *api.EnableEntityRequest) (*emptypb.Empty, error) {
	if err := c.endpoint.Entity.Enable(ctx, common.EntityId(req.Id)); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

func (c ControllerEntity) DisabledEntity(ctx context.Context, req *api.DisableEntityRequest) (*emptypb.Empty, error) {
	if err := c.endpoint.Entity.Disable(ctx, common.EntityId(req.Id)); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}
