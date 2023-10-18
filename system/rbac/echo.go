// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package rbac

import (
	"context"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/controllers"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/jwt_manager"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"strings"
)

var (
	log = logger.MustGetLogger("rbac")
)

// EchoAccessFilter ...
type EchoAccessFilter struct {
	adaptors            *adaptors.Adaptors
	jwtManager          jwt_manager.JwtManager
	accessListService   access_list.AccessListService
	internalServerError error
	config              *m.AppConfig
}

// NewEchoAccessFilter ...
func NewEchoAccessFilter(adaptors *adaptors.Adaptors,
	jwtManager jwt_manager.JwtManager,
	accessListService access_list.AccessListService,
	config *m.AppConfig) *EchoAccessFilter {
	return &EchoAccessFilter{
		adaptors:          adaptors,
		jwtManager:        jwtManager,
		accessListService: accessListService,
		config:            config,
	}
}

func (f *EchoAccessFilter) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		requestURI := c.Request().URL
		method := strings.ToLower(c.Request().Method)

		if f.config.GodMode {
			return next(c)
		}

		// get access token from meta
		var accessToken = f.getAccessToken(c)
		if accessToken == "" {
			return f.HTTP401(c, apperr.ErrUnauthorized)
		}

		claims, err := f.jwtManager.Verify(accessToken)
		if err != nil {
			return f.HTTP401(c, apperr.ErrUnauthorized)
		}

		// if id == 1 is admin
		if claims.UserId == 1 || claims.RoleName == "admin" {
			if err = f.getUser(claims.UserId, c); err != nil {
				return f.HTTP401(c, apperr.ErrUnauthorized)
			}
			return next(c)
		}

		var accessList access_list.AccessList
		if accessList, err = f.accessListService.GetFullAccessList(c.Request().Context(), claims.RoleName); err != nil {
			return f.HTTP401(c, apperr.ErrUnauthorized)
		}

		// check access filter
		if ret := f.accessDecision(requestURI.Path, method, accessList); ret {
			if err = f.getUser(claims.UserId, c); err != nil {
				return f.HTTP401(c, apperr.ErrUnauthorized)
			}
			return next(c)
		}

		log.Warnf(fmt.Sprintf("access denied: role(%s) [%s] url(%s)", claims.RoleName, method, requestURI.Path))

		c.Error(f.HTTP401(c, apperr.ErrUnauthorized))

		return nil
	}
}

func (f *EchoAccessFilter) getUser(userId int64, c echo.Context) error {
	user, err := f.adaptors.User.GetById(context.Background(), userId)
	if err != nil {
		return err
	}
	c.Set("currentUser", user)
	return nil
}

func (f *EchoAccessFilter) accessDecision(params, method string, accessList access_list.AccessList) bool {

	for _, action := range []string{"/stream/[\\wW\\.0-9]+/channel/[0-9]+/mse", "/v1/image/upload"} {
		if ok, _ := regexp.MatchString(action, params); ok {
			return true
		}
	}

	return false
}

func (f *EchoAccessFilter) getAccessToken(c echo.Context) (accessToken string) {
	accessToken = c.Request().Header.Get("authorization")
	if accessToken != "" {
		return
	}
	accessToken = c.Request().URL.Query().Get("access_token")
	return
}

// HTTP401 ...
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
