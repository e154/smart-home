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
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// ControllerUser ...
type ControllerUser struct {
	*ControllerCommon
}

// NewControllerUser ...
func NewControllerUser(common *ControllerCommon) *ControllerUser {
	return &ControllerUser{ControllerCommon: common}
}

// swagger:operation POST /user userAdd
// ---
// parameters:
// - description: user params
//   in: body
//   name: user
//   required: true
//   schema:
//     $ref: '#/definitions/NewUser'
// summary: add new user
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - user
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/UserFull'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerUser) Add(ctx *gin.Context) {

	var params models.NewUser
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	var currentUser *m.User
	if user, ok := ctx.Get("currentUser"); ok {
		currentUser = user.(*m.User)
	}

	user := &m.User{}
	common.Copy(&user, &params, common.JsonEngine)

	if params.Password == params.PasswordRepeat {
		user.SetPass(params.Password)
	}

	user, errs, err := c.endpoint.User.Add(user, currentUser)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.UserFull{}
	common.Copy(&result, &user)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /user/{id} userGetById
// ---
// parameters:
// - description: User ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get user by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - user
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/UserFull'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerUser) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	user, err := c.endpoint.User.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.UserFull{}
	common.Copy(&result, &user)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /user/{id} userUpdateById
// ---
// parameters:
// - description: User ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update user params
//   in: body
//   name: user
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateUser'
//     type: object
// summary: update user by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - user
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/UserFull'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerUser) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateUser{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	user := &m.User{}
	common.Copy(&user, &params, common.JsonEngine)

	err400 := NewError(400)
	if params.Password != params.PasswordRepeat {
		err400.AddFieldf("password", FieldNotValid)
	}

	if params.Password != "" {
		user.SetPass(params.Password)
	}

	if err400.Errors() {
		err400.Send(ctx)
		return
	}

	user, errs, err := c.endpoint.User.Update(user)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.UserFull{}
	common.Copy(&result, &user)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /user/{id}/update_status userUpdateStatus
// ---
// parameters:
// - description: User ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: status params
//   in: body
//   name: user_status
//   required: true
//   schema:
//     $ref: '#/definitions/UserUpdateStatusRequest'
// summary: update user status
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - user
// responses:
//   "200":
//     $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerUser) UpdateStatus(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n := &models.UserUpdateStatusRequest{}
	if err := ctx.ShouldBindJSON(n); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err = c.endpoint.User.UpdateStatus(int64(aid), n.Status); err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation GET /users userList
// ---
// summary: get user list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - user
// parameters:
// - default: 10
//   description: limit
//   in: query
//   name: limit
//   required: true
//   type: integer
// - default: 0
//   description: offset
//   in: query
//   name: offset
//   required: true
//   type: integer
// - default: DESC
//   description: order
//   in: query
//   name: order
//   type: string
// - default: id
//   description: sort_by
//   in: query
//   name: sort_by
//   type: string
// responses:
//   "200":
//	   $ref: '#/responses/UserList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerUser) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.User.GetList(limit, offset, order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.UserShot, 0)
	common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /user/{id} userDeleteById
// ---
// parameters:
// - description: User ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete user by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - user
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerUser) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.User.Delete(int64(aid)); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}
