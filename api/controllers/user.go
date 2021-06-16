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
)

type ControllerUser struct {
	*ControllerCommon
}

func NewControllerUser(common *ControllerCommon) ControllerUser {
	return ControllerUser{
		ControllerCommon: common,
	}
}

func (c ControllerUser) AddUser(context.Context, *api.NewtUser) (*api.UserFull, error) {

	var userFull = &api.UserFull{}
	return userFull, nil
}

func (c ControllerUser) GetUserById(_ context.Context, req *api.GetUserByIdRequest) (*api.UserFull, error) {

	var userFull = &api.UserFull{}

	user, err := c.endpoint.User.GetById(req.Id)
	if err != nil {
		return nil, err
	}

	_ = common.Copy(userFull, user, common.JsonEngine)

	return userFull, nil
}
