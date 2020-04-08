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

// ControllerWorkflowScenario ...
type ControllerWorkflowScenario struct {
	*ControllerCommon
}

// NewControllerWorkflowScenario ...
func NewControllerWorkflowScenario(common *ControllerCommon) *ControllerWorkflowScenario {
	return &ControllerWorkflowScenario{ControllerCommon: common}
}

// swagger:operation POST /workflow/{id}/scenario workflowScenarioAdd
// ---
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: workflow scenario params
//   in: body
//   name: workflow_scenario
//   required: true
//   schema:
//     $ref: '#/definitions/NewWorkflowScenario'
// summary: add new workflow scenario
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow_scenario
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/WorkflowScenario'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerWorkflowScenario) Add(ctx *gin.Context) {

	var params models.NewWorkflowScenario
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	scenario := &m.WorkflowScenario{
		Name:       params.Name,
		SystemName: params.SystemName,
		WorkflowId: params.WorkflowId,
	}

	scenario, errs, err := c.endpoint.WorkflowScenario.Add(scenario)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.WorkflowScenario{}
	common.Copy(&result, &scenario, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /workflow/{id}/scenario/{scenario_id} workflowScenarioGetById
// ---
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: WorkflowScenario ID
//   in: path
//   name: scenario_id
//   required: true
//   type: integer
// summary: get workflow scenario by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow_scenario
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/WorkflowScenario'
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
func (c ControllerWorkflowScenario) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	workflowId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}
	id = ctx.Param("scenario_id")
	scenarioId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	workflowScenario, err := c.endpoint.WorkflowScenario.GetById(int64(workflowId), int64(scenarioId))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.WorkflowScenario{}
	common.Copy(&result, &workflowScenario, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /workflow/{id}/scenario/{scenario_id} workflowScenarioUpdateById
// ---
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: WorkflowScenario ID
//   in: path
//   name: scenario_id
//   required: true
//   type: integer
// - description: Update workflow scenario params
//   in: body
//   name: workflowScenario
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateWorkflowScenario'
//     type: object
// summary: update workflow scenario by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow_scenario
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/WorkflowScenario'
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
func (c ControllerWorkflowScenario) Update(ctx *gin.Context) {

	workflowId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateWorkflowScenario{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.WorkflowId = int64(workflowId)
	workflowScenario := &m.WorkflowScenario{}
	common.Copy(&workflowScenario, &params, common.JsonEngine)

	workflowScenario, errs, err := c.endpoint.WorkflowScenario.Update(workflowScenario)
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

	result := &models.WorkflowScenario{}
	common.Copy(&result, &workflowScenario, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /workflow/{id}/scenarios workflowScenarioList
// ---
// summary: get workflow scenario list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow_scenario
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// responses:
//   "200":
//	   $ref: '#/responses/WorkflowScenarios'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerWorkflowScenario) GetList(ctx *gin.Context) {

	id := ctx.Param("id")
	workflowId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	items, _, err := c.endpoint.WorkflowScenario.GetList(int64(workflowId))
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.WorkflowScenario, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Item("scenarios", result).Send(ctx)
	return
}

// swagger:operation DELETE /workflow/{id}/scenario/{scenario_id} workflowScenarioDeleteById
// ---
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Workflow scenario ID
//   in: path
//   name: scenario_id
//   required: true
//   type: integer
// summary: delete workflow scenario by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow_scenario
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
func (c ControllerWorkflowScenario) Delete(ctx *gin.Context) {

	id := ctx.Param("scenario_id")
	workflowScenarioId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.WorkflowScenario.Delete(int64(workflowScenarioId)); err != nil {
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

// swagger:operation GET /workflow/{id}/scenarios/search workflowScenarioSearch
// ---
// summary: search workflow scenario
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow_scenario
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
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
//	   $ref: '#/responses/WorkflowScenarioSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerWorkflowScenario) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.WorkflowScenario.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.WorkflowScenario, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Item("scenarios", result).Send(ctx)
}
