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
	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerArea ...
type ControllerArea struct {
	*ControllerCommon
}

// NewControllerArea ...
func NewControllerArea(common *ControllerCommon) ControllerArea {
	return ControllerArea{
		ControllerCommon: common,
	}
}

func (c ControllerArea) AddArea(ctx context.Context, req *api.NewAreaRequest) (*api.Area, error) {

	area := c.dto.Area.AddArea(req)

	area, errs, err := c.endpoint.Area.Add(ctx, area)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Area.ToArea(area), nil
}

func (c ControllerArea) UpdateArea(ctx context.Context, req *api.UpdateAreaRequest) (*api.Area, error) {

	area := c.dto.Area.UpdateArea(req)

	area, errs, err := c.endpoint.Area.Update(ctx, area)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Area.ToArea(area), nil
}

func (c ControllerArea) GetAreaById(ctx context.Context, req *api.GetAreaRequest) (*api.Area, error) {

	area, err := c.endpoint.Area.GetById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Area.ToArea(area), nil
}

func (c ControllerArea) GetAreaList(ctx context.Context, req *api.GetAreaListRequest) (*api.GetAreaListResult, error) {

	pagination := c.Pagination(req.Limit, req.Offset, req.Order, req.SortBy)
	items, total, err := c.endpoint.Area.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Area.ToListResult(items, uint64(total), pagination), nil
}

func (c ControllerArea) DeleteArea(ctx context.Context, req *api.DeleteAreaRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Area.Delete(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// SearchArea ...
func (c ControllerArea) SearchArea(ctx context.Context, req *api.SearchAreaRequest) (*api.SearchAreaResult, error) {

	search := c.Search(req.Query, req.Limit, req.Offset)
	items, _, err := c.endpoint.Area.Search(ctx, search.Query, search.Limit, search.Offset)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Area.ToSearchResult(items), nil
}
