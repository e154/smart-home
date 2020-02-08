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

type ControllerNode struct {
	*ControllerCommon
}

func NewControllerNode(common *ControllerCommon) *ControllerNode {
	return &ControllerNode{ControllerCommon: common}
}

// swagger:operation POST /node nodeAdd
// ---
// parameters:
// - description: node params
//   in: body
//   name: node
//   required: true
//   schema:
//     $ref: '#/definitions/NewNode'
//     type: object
// summary: add new node
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - node
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Node'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerNode) Add(ctx *gin.Context) {

	params := &models.NewNode{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if params.Password != params.PasswordRepeat {
		NewError(400, "bad password repeat").Send(ctx)
		return
	}

	node := &m.Node{}
	common.Copy(&node, &params, common.JsonEngine)

	node, errs, err := c.endpoint.Node.Add(node)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Node{}
	if err = common.Copy(&result, &node, common.JsonEngine); err != nil {
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /node/{id} nodeGetById
// ---
// parameters:
// - description: Node ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get node by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - node
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Node'
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
func (c ControllerNode) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	node, err := c.endpoint.Node.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Node{}
	common.Copy(&result, &node, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /node/{id} nodeUpdateById
// ---
// parameters:
// - description: Node ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update node params
//   in: body
//   name: node
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateNode'
//     type: object
// summary: update node by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - node
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Node'
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
func (c ControllerNode) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateNode{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if params.Password != params.PasswordRepeat {
		NewError(400, "bad password repeat").Send(ctx)
		return
	}

	params.Id = int64(aid)

	node := &m.Node{}
	common.Copy(&node, &params, common.JsonEngine)

	node, errs, err := c.endpoint.Node.Update(node)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Node{}
	common.Copy(&result, &node, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /nodes nodeList
// ---
// summary: get node list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - node
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
//	   $ref: '#/responses/NodeList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerNode) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Node.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Node, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /node/{id} nodeDeleteById
// ---
// parameters:
// - description: Node ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete node by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - node
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
func (c ControllerNode) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.Node.Delete(int64(aid)); err != nil {
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

// swagger:operation GET /nodes/search nodeSearch
// ---
// summary: search node
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - node
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
//	   $ref: '#/responses/NodeSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerNode) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.Node.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Node, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Item("nodes", result)
	resp.Send(ctx)
}
