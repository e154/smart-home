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

type ControllerFlow struct {
	*ControllerCommon
}

func NewControllerFlow(common *ControllerCommon) *ControllerFlow {
	return &ControllerFlow{ControllerCommon: common}
}

// swagger:operation POST /flow flowAdd
// ---
// parameters:
// - description: flow params
//   in: body
//   name: flow
//   required: true
//   schema:
//     $ref: '#/definitions/NewFlow'
//     type: object
// summary: add new flow
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Flow'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerFlow) Add(ctx *gin.Context) {

	params := &models.NewFlow{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	flow := &m.Flow{}
	_ = common.Copy(&flow, &params)

	if params.Workflow.Id != 0 {
		flow.WorkflowId = params.Workflow.Id
	}

	if params.Scenario.Id != 0 {
		flow.WorkflowScenarioId = params.Scenario.Id
	}

	flow, errs, err := c.endpoint.Flow.Add(flow)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Flow{}
	_ = common.Copy(&result, &flow)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /flow/{id} flowGetById
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get flow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Flow'
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
func (c ControllerFlow) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	flow, err := c.endpoint.Flow.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Flow{}
	_ = common.Copy(&result, &flow, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /flow/{id}/redactor flowGetRedactor
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get flow redactor data by flow id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/RedactorFlow'
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
func (c ControllerFlow) GetRedactor(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	redactorFlow, err := c.endpoint.Flow.GetRedactor(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.RedactorFlow{}
	common.Copy(&result, &redactorFlow, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /flow/{id}/redactor flowUpdateRedactor
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: flow redactor params
//   in: body
//   name: flow
//   required: true
//   schema:
//     $ref: '#/definitions/RedactorFlow'
//     type: object
// summary: update flow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/RedactorFlow'
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
func (c ControllerFlow) UpdateRedactor(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.RedactorFlow{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	flowRedactor := &m.RedactorFlow{}
	common.Copy(&flowRedactor, &params, common.JsonEngine)

	flowRedactor, errs, err := c.endpoint.Flow.UpdateRedactor(flowRedactor)
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

	result := &models.RedactorFlow{}
	_ = common.Copy(&result, &flowRedactor)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /flow/{id} flowUpdateById
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update flow params
//   in: body
//   name: flow
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateFlow'
//     type: object
// summary: update flow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Flow'
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
func (c ControllerFlow) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateFlow{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	flow := &m.Flow{}
	if err = common.Copy(&flow, &params); err != nil {
		return
	}

	flow, errs, err := c.endpoint.Flow.Update(flow)
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

	result := &models.Flow{}
	common.Copy(&result, &flow)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /flows flowList
// ---
// summary: get flow list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
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
//	   $ref: '#/responses/FlowList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerFlow) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Flow.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.FlowShort, 0)
	_ = common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
}

// swagger:operation DELETE /flow/{id} flowDeleteById
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete flow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
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
func (c ControllerFlow) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err = c.endpoint.Flow.Delete(int64(aid)); err != nil {
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

// swagger:operation GET /flows/search flowSearch
// ---
// summary: search flow
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
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
//	   $ref: '#/responses/FlowSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerFlow) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	list, _, err := c.endpoint.Flow.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Flow, 0)
	_ = common.Copy(&result, &list)

	resp := NewSuccess()
	resp.Item("flows", result)
	resp.Send(ctx)
}
