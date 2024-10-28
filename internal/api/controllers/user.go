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
	"github.com/e154/smart-home/pkg/apperr"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/labstack/echo/v4"
)

// ControllerUser ...
type ControllerUser struct {
	*ControllerCommon
}

// NewControllerUser ...
func NewControllerUser(common *ControllerCommon) *ControllerUser {
	return &ControllerUser{
		ControllerCommon: common,
	}
}

// AddUser ...
func (c ControllerUser) UserServiceAddUser(ctx echo.Context, _ stub.UserServiceAddUserParams) error {

	obj := &stub.ApiNewtUserRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	user := c.dto.User.AddUserRequest(obj)

	if obj.Password == obj.PasswordRepeat {
		_ = user.SetPass(obj.Password)
	}

	currentUser, err := c.currentUser(ctx)
	if err != nil {
		return c.ERROR(ctx, apperr.ErrInternal)
	}

	var userF *m.User
	userF, err = c.endpoint.User.Add(ctx.Request().Context(), user, currentUser)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.User.ToUserFull(userF)))
}

// GetUserById ...
func (c ControllerUser) UserServiceGetUserById(ctx echo.Context, id int64) error {

	user, err := c.endpoint.User.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.User.ToUserFull(user)))
}

// UpdateUserById ...
func (c ControllerUser) UserServiceUpdateUserById(ctx echo.Context, id int64, _ stub.UserServiceUpdateUserByIdParams) error {

	obj := &stub.UserServiceUpdateUserByIdJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	user := c.dto.User.UpdateUserByIdRequest(obj, id)

	if obj.Password != obj.PasswordRepeat {
		return c.ERROR(ctx, apperr.ErrPassNotValid)
	}

	if obj.Password != "" {
		_ = user.SetPass(obj.Password)
	}

	user, err := c.endpoint.User.Update(ctx.Request().Context(), user)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.User.ToUserFull(user)))
}

// GetUserList ...
func (c ControllerUser) UserServiceGetUserList(ctx echo.Context, params stub.UserServiceGetUserListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.User.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.User.ToListResult(items), total, pagination))
}

// DeleteUserById ...
func (c ControllerUser) UserServiceDeleteUserById(ctx echo.Context, id int64) error {

	if err := c.endpoint.User.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}
