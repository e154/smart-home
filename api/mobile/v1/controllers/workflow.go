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
	"github.com/gin-gonic/gin"
	"strconv"
)

// ControllerWorkflow ...
type ControllerWorkflow struct {
	*ControllerCommon
}

// NewControllerWorkflow ...
func NewControllerWorkflow(common *ControllerCommon) *ControllerWorkflow {
	return &ControllerWorkflow{ControllerCommon: common}
}

// swagger:operation GET /workflow/{id} workflowGetById
// ---
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get workflow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Workflow'
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
func (c ControllerWorkflow) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		//log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	workflow, err := c.endpoint.Workflow.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Workflow{}
	common.Copy(&result, &workflow, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /workflows workflowList
// ---
// summary: get workflow list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow
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
//	   $ref: '#/responses/WorkflowList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerWorkflow) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Workflow.GetList(int64(limit), int64(offset), order, sortBy, true)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.WorkflowShort, 0)
	_ = common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// Workflow godoc
// swagger:operation PUT /workflow/{id}/update_scenario workflowUpdateScenario
// ---
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update workflow scenario params
//   in: body
//   name: workflowUpdateWorkflowScenario
//   required: true
//   schema:
//     $ref: '#/definitions/WorkflowUpdateWorkflowScenario'
// summary: update workflow scenario
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow
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
func (c ControllerWorkflow) UpdateScenario(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	workflowScenario := &models.WorkflowUpdateWorkflowScenario{}
	if err := ctx.ShouldBindJSON(&workflowScenario); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	err = c.endpoint.Workflow.UpdateScenario(int64(aid), workflowScenario.WorkflowScenarioId)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)

	return
}
