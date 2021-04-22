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
)

// ControllerEntity ...
type ControllerEntity struct {
	*ControllerCommon
}

// NewControllerEntity ...
func NewControllerEntity(common *ControllerCommon) *ControllerEntity {
	return &ControllerEntity{ControllerCommon: common}
}

// swagger:operation POST /entity entityAdd
// ---
// parameters:
// - description: entity params
//   in: body
//   name: entity
//   required: true
//   schema:
//     $ref: '#/definitions/NewEntity'
//     type: object
// summary: add new entity
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - entity
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Entity'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerEntity) Add(ctx *gin.Context) {

	params := &models.NewEntity{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	entity := &m.Entity{}
	common.Copy(&entity, &params, common.JsonEngine)

	entity, errs, err := c.endpoint.Entity.Add(entity)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Entity{}
	if err = common.Copy(&result, &entity, common.JsonEngine); err != nil {
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /entity/{id} entityGetById
// ---
// parameters:
// - description: Entity ID
//   in: path
//   name: id
//   required: true
//   type: string
// summary: get entity by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - entity
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Entity'
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
func (c ControllerEntity) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		NewError(400, "Bad request").Send(ctx)
		return
	}

	entityId := common.EntityId(id)

	entity, err := c.endpoint.Entity.GetById(entityId)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Entity{}
	common.Copy(&result, &entity, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /entity/{id} entityUpdateById
// ---
// parameters:
// - description: Entity ID
//   in: path
//   name: id
//   required: true
//   type: string
// - description: Update entity params
//   in: body
//   name: entity
//   required: true
//   schema:
//     $ref: '#/definitions/Update'
//     type: object
// summary: update entity by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - entity
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Entity'
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
func (c ControllerEntity) Update(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		NewError(400, "Bad request").Send(ctx)
		return
	}

	params := &models.UpdateEntity{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	entity := &m.Entity{}
	common.Copy(&entity, &params, common.JsonEngine)
	entity.Id = common.EntityId(id)

	entity, errs, err := c.endpoint.Entity.Update(entity)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Entity{}
	common.Copy(&result, &entity, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /entities entityList
// ---
// summary: get entity list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - entity
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
//	   $ref: '#/responses/EntityList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerEntity) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Entity.List(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Entity, 0)
	common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /entity/{id} entityDeleteById
// ---
// parameters:
// - description: Entity ID
//   in: path
//   name: id
//   required: true
//   type: string
// summary: delete entity by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - entity
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
func (c ControllerEntity) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		NewError(400, "Bad request").Send(ctx)
		return
	}

	if err := c.endpoint.Entity.Delete(common.EntityId(id)); err != nil {
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

// swagger:operation GET /entities/search entitySearch
// ---
// summary: search entity
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - entity
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
//	   $ref: '#/responses/EntitySearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerEntity) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.Entity.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Entity, 0)
	common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Item("entities", result)
	resp.Send(ctx)
}
