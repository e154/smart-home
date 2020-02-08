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
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

type ControllerMap struct {
	*ControllerCommon
}

func NewControllerMap(common *ControllerCommon) *ControllerMap {
	return &ControllerMap{ControllerCommon: common}
}

// swagger:operation POST /map mapAdd
// ---
// parameters:
// - description: map params
//   in: body
//   name: map
//   required: true
//   schema:
//     $ref: '#/definitions/NewMap'
//     type: object
// summary: add new map
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Map'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMap) Add(ctx *gin.Context) {

	params := &models.NewMap{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	m := &m.Map{}
	common.Copy(&m, &params)

	m, errs, err := c.endpoint.Map.Add(m)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Map{}
	common.Copy(&result, &m, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /map/{id} mapGetById
// ---
// parameters:
// - description: Map ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get map by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Map'
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
func (c ControllerMap) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	m, err := c.endpoint.Map.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Map{}
	common.Copy(&result, &m, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /map/{id}/full mapFullGetById
// ---
// parameters:
// - description: Map ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get map full by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapFull'
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
func (c ControllerMap) GetFullMap(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	m, err := c.endpoint.Map.GetFullById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.MapFull{}
	common.Copy(&result, &m, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /map/{id} mapUpdateById
// ---
// parameters:
// - description: Map ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update map params
//   in: body
//   name: map
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateMap'
//     type: object
// summary: update map by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Map'
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
func (c ControllerMap) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateMap{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	m := &m.Map{}
	common.Copy(&m, &params, common.JsonEngine)

	m, errs, err := c.endpoint.Map.Update(m)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	result := &models.Map{}
	common.Copy(&result, &m, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /maps mapList
// ---
// summary: get map list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
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
//	   $ref: '#/responses/MapList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMap) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Map.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Map, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /map/{id} mapDeleteById
// ---
// parameters:
// - description: Map ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete map by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
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
func (c ControllerMap) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.Map.Delete(int64(aid)); err != nil {
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

// swagger:operation GET /maps/search mapSearch
// ---
// summary: search map
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
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
//	   $ref: '#/responses/MapSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMap) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.Map.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Map, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Item("maps", result)
	resp.Send(ctx)
}
