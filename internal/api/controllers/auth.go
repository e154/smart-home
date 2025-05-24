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
// You should have received c copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package controllers

import (
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/apperr"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/labstack/echo/v4"
)

// ControllerAuth ...
type ControllerAuth struct {
	*ControllerCommon
}

// NewControllerAuth ...
func NewControllerAuth(common *ControllerCommon) *ControllerAuth {
	return &ControllerAuth{
		ControllerCommon: common,
	}
}

// Signin ...
func (c ControllerAuth) AuthServiceSignin(ctx echo.Context) error {

	username, pass, _ := c.parseBasicAuth(ctx.Request().Header.Get("authorization"))

	var user *m.User
	var accessToken string

	var ip string
	if _ip := ctx.Request().Header.Get("ip"); _ip != "" {
		if ok, _ := c.validation.ValidVar(_ip, "ip", "required,ipv4"); ok {
			ip = _ip
		}
	}

	var err error
	if user, accessToken, err = c.endpoint.Auth.SignIn(ctx.Request().Context(), username, pass, ip); err != nil {
		return c.ERROR(ctx, apperr.ErrUnauthorized)
	}

	currentUser := &stub.ApiCurrentUser{}
	_ = common.Copy(&currentUser, &user, common.JsonEngine)

	resp := &stub.ApiSigninResponse{
		CurrentUser: currentUser,
		AccessToken: accessToken,
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, resp))
}

// Signout ...
func (c ControllerAuth) AuthServiceSignout(ctx echo.Context) error {

	currentUser, err := c.currentUser(ctx)
	if err != nil {
		return c.ERROR(ctx, apperr.ErrUnauthorized)
	}

	if err = c.endpoint.Auth.SignOut(ctx.Request().Context(), currentUser); err != nil {
		return c.ERROR(ctx, apperr.ErrUnauthorized)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// AccessList ...
func (c ControllerAuth) AuthServiceAccessList(ctx echo.Context) error {

	currentUser, err := c.currentUser(ctx)
	if err != nil {
		return c.ERROR(ctx, apperr.ErrUnauthorized)
	}

	accessList, err := c.endpoint.Auth.AccessList(ctx.Request().Context(), currentUser, c.accessList)
	if err != nil {
		return c.ERROR(ctx, apperr.ErrUnauthorized)
	}

	resp := &stub.ApiAccessListResponse{
		AccessList: c.dto.Role.ToAccessListResult(accessList),
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, resp))
}

// PasswordReset ...
func (c ControllerAuth) AuthServicePasswordReset(ctx echo.Context, _ stub.AuthServicePasswordResetParams) error {

	obj := &stub.ApiPasswordResetRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	if err := c.endpoint.Auth.PasswordReset(ctx.Request().Context(), obj.Email, obj.Token, obj.NewPassword); err != nil {
		return c.ERROR(ctx, apperr.ErrUserNotFound)
	}
	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}
