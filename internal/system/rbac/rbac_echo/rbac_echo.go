// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023-2024, Filippov Alex
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

package rbac_echo

import (
	"net/http"
	"strings"

	"github.com/e154/smart-home/internal/api/controllers"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"

	"github.com/labstack/echo/v4"
)

type EchoAccessFilter struct {
	config        *m.AppConfig
	Authenticator plugins.Authorization
}

func NewEchoAccessFilter(config *m.AppConfig, authenticator plugins.Authorization) *EchoAccessFilter {
	return &EchoAccessFilter{
		config:        config,
		Authenticator: authenticator,
	}
}

func (f *EchoAccessFilter) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		requestURI := c.Request().URL
		method := strings.ToLower(c.Request().Method)

		var accessToken = f.getAccessToken(c)
		if accessToken == "" {
			return f.HTTP401(c, apperr.ErrUnauthorized)
		}

		user, root, err := f.Authenticator.AuthREST(c.Request().Context(), accessToken, requestURI, method)
		if err != nil {
			return f.HTTP401(c, apperr.ErrUnauthorized)
		}
		if user != nil {
			c.Set("currentUser", user)
			c.Set("root", f.config.RootMode || root)
			return next(c)
		}

		return f.HTTP401(c, apperr.ErrUnauthorized)
	}
}

func (f *EchoAccessFilter) getAccessToken(c echo.Context) (accessToken string) {
	accessToken = c.Request().Header.Get("authorization")
	if accessToken != "" {
		return
	}
	accessToken = c.Request().URL.Query().Get("access_token")
	return
}

func (f *EchoAccessFilter) HTTP401(ctx echo.Context, err error) error {
	e := apperr.GetError(err)
	if e != nil {
		return ctx.JSON(http.StatusUnauthorized, controllers.ResponseWithError(ctx, &controllers.ErrorBase{
			Code:    common.String(e.Code()),
			Message: common.String(e.Message()),
		}))
	}
	return ctx.JSON(http.StatusUnauthorized, controllers.ResponseWithError(ctx, &controllers.ErrorBase{
		Code: common.String("UNAUTHORIZED"),
	}))
}
