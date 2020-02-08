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
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
)

type ControllerRole struct {
	*ControllerCommon
}

func NewControllerRole(common *ControllerCommon) *ControllerRole {
	return &ControllerRole{ControllerCommon: common}
}

// swagger:operation POST /role roleAdd
// ---
// parameters:
// - description: role params
//   in: body
//   name: role
//   required: true
//   schema:
//     $ref: '#/definitions/NewRole'
//     type: object
// summary: add new role
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Role'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) Add(ctx *gin.Context) {

	var params models.NewRole
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	role := &m.Role{
		Name:        params.Name,
		Description: params.Description,
	}

	role, errs, err := c.endpoint.Role.Add(role)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Role{}
	common.Copy(&result, &role, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /role/{name} roleGetById
// ---
// parameters:
// - description: Role name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: get role by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Role'
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
func (c ControllerRole) GetByName(ctx *gin.Context) {

	name := ctx.Param("name")
	role, err := c.endpoint.Role.GetByName(name)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Role{}
	common.Copy(&result, &role, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /role/{name}/access_list roleGetById
// ---
// parameters:
// - description: Role name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: get access list by role name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//     description: OK
//     schema:
//       type: object
//       properties:
//         access_list:
//           $ref: '#/definitions/AccessList'
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
func (c ControllerRole) GetAccessList(ctx *gin.Context) {

	name := ctx.Param("name")
	accessList, err := c.endpoint.Role.GetAccessList(name, c.accessList)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := make(models.AccessList)
	common.Copy(&result, accessList, common.JsonEngine)

	resp := NewSuccess()
	resp.Item("access_list", result).Send(ctx)
}

// swagger:operation PUT /role/{name}/access_list roleUpdateById
// ---
// parameters:
// - description: Role name
//   in: path
//   name: name
//   required: true
//   type: string
// - description: Update access list params
//   in: body
//   name: access_list_diff
//   required: true
//   schema:
//     $ref: '#/definitions/AccessListDiff'
// summary: update role access list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
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
func (c ControllerRole) UpdateAccessList(ctx *gin.Context) {

	accessListDif := make(map[string]map[string]bool)
	if err := ctx.ShouldBindJSON(&accessListDif); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	name := ctx.Param("name")
	if err := c.endpoint.Role.UpdateAccessList(name, accessListDif); err != nil {
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

// swagger:operation PUT /role/{name} roleUpdateById
// ---
// parameters:
// - description: Role ID
//   in: path
//   name: name
//   required: true
//   type: string
// - description: Update role params
//   in: body
//   name: role
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateRole'
// summary: update role by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Role'
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
func (c ControllerRole) Update(ctx *gin.Context) {

	name := ctx.Param("name")
	params := &models.UpdateRole{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	params.Name = name

	role := &m.Role{}
	common.Copy(&role, &params)

	role, errs, err := c.endpoint.Role.Update(role)
	if len(errs) > 0 {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Role{}
	common.Copy(&result, &role, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /roles roleList
// ---
// summary: get role list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
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
//	   $ref: '#/responses/RoleList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Role.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Role, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /role/{name} roleDeleteById
// ---
// parameters:
// - description: Role ID
//   in: path
//   name: name
//   required: true
//   type: string
// summary: delete role by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
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
func (c ControllerRole) Delete(ctx *gin.Context) {

	name := ctx.Param("name")
	_, err := c.adaptors.Role.GetByName(name)
	if err != nil {
		NewError(404, err).Send(ctx)
		return
	}

	if err := c.endpoint.Role.Delete(name); err != nil {
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

// swagger:operation GET /roles/search roleSearch
// ---
// summary: search role
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// parameters:
// - description: query
//   in: query
//   name: query
//   type: string
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
// responses:
//   "200":
//	   $ref: '#/responses/RoleSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.Role.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Role, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Item("roles", result)
	resp.Send(ctx)
}
