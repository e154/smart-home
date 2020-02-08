// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	"net/http"
	"github.com/e154/smart-home/common"
)

type ControllerAuth struct {
	*ControllerCommon
}

func NewControllerAuth(common *ControllerCommon) *ControllerAuth {
	return &ControllerAuth{ControllerCommon: common}
}

// swagger:operation POST /signin authSignin
// ---
// summary: sign in
// description:
// security:
// - BasicAuth: []
// tags:
// - auth
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/AuthSignInResponse'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAuth) SignIn(ctx *gin.Context) {

	email, password, ok := ctx.Request.BasicAuth()
	if !ok {
		err := NewError(400, "bad request")
		if email == "" {
			err.AddField("common.field_not_blank", "The field can't be empty", "email")
		}
		if password == "" {
			err.AddField("common.field_not_blank", "The field can't be empty", "password")
		}
		err.Send(ctx)
		return
	}

	user, accessToken, err := c.endpoint.Auth.SignIn(email, password, ctx.ClientIP())
	if err != nil {
		code := 500
		switch err.Error() {
		case "user not found":
			code = 401
		case "password not valid":
			code = 403
		}

		NewError(code, err.Error()).Send(ctx)
		return
	}

	currentUser := &models.CurrentUser{}
	_ = common.Copy(&currentUser, &user, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(&models.AuthSignInResponse{
		AccessToken: accessToken,
		CurrentUser: currentUser,
	}).Send(ctx)
}

// swagger:operation POST /signout authSignout
// ---
// summary: sign out
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - auth
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAuth) SignOut(ctx *gin.Context) {

	user, err := c.getUser(ctx)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	if err := c.endpoint.Auth.SignOut(user); err != nil {
		NewError(500, err.Error()).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
}

// swagger:operation POST /recovery authRecovery
// ---
// summary: recovery access
// description:
// tags:
// - auth
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAuth) Recovery(ctx *gin.Context) {
	ctx.String(http.StatusOK, "operation Recovery has not yet been implemented")
}

// swagger:operation POST /reset authReset
// ---
// summary: reset access
// description:
// tags:
// - auth
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAuth) Reset(ctx *gin.Context) {
	ctx.String(http.StatusOK, "operation Reset has not yet been implemented")
}

// swagger:operation GET /access_list authGetAccessList
// ---
// summary: get user access list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - auth
// responses:
//   "200":
//     description: OK
//     schema:
//       type: object
//       properties:
//         access_list:
//           $ref: '#/definitions/AccessList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAuth) AccessList(ctx *gin.Context) {

	user, err := c.getUser(ctx)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	accessList, err := c.endpoint.Auth.AccessList(user, c.accessList)
	if err != nil {
		NewError(500, err.Error()).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(&map[string]interface{}{
		"access_list": accessList,
	}).Send(ctx)
}
