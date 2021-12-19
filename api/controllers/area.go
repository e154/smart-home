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

// ControllerArea ...
type ControllerArea struct {
	*ControllerCommon
}

func (c ControllerArea) AddArea(ctx context.Context, request *api.NewAreaRequest) (*api.Area, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}

func (c ControllerArea) UpdateAreaByName(ctx context.Context, request *api.UpdateAreaRequest) (*api.Area, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}

func (c ControllerArea) GetAreaByName(ctx context.Context, request *api.GetAreaRequest) (*api.Area, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}

func (c ControllerArea) GetAreaList(ctx context.Context, request *api.GetAreaListRequest) (*api.GetAreaListResult, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}

func (c ControllerArea) DeleteAreaByName(ctx context.Context, request *api.DeleteAreaRequest) (*emptypb.Empty, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}

// NewControllerArea ...
func NewControllerArea(common *ControllerCommon) ControllerArea {
	return ControllerArea{
		ControllerCommon: common,
	}
}

