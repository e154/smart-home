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

// ControllerInteract ...
type ControllerInteract struct {
	*ControllerCommon
}

// NewControllerInteract ...
func NewControllerInteract(common *ControllerCommon) ControllerInteract {
	return ControllerInteract{
		ControllerCommon: common,
	}
}

func (c ControllerInteract) EntityCallAction(ctx context.Context, req *api.EntityCallActionRequest) (*emptypb.Empty, error) {

	errs, err := c.endpoint.Interact.EntityCallAction(ctx, req.Id, req.Name, nil)
	if err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return &emptypb.Empty{}, nil
}