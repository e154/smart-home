package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerWorkflow struct {
	*ControllerCommon
}

func NewControllerWorkflow(common *ControllerCommon) *ControllerWorkflow {
	return &ControllerWorkflow{ControllerCommon: common}
}

// swagger:operation POST /workflow workflowAdd
// ---
// parameters:
// - description: workflow params
//   in: body
//   name: workflow
//   required: true
//   schema:
//     $ref: '#/definitions/NewWorkflow'
//     type: object
// summary: add new workflow
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
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerWorkflow) Add(ctx *gin.Context) {

	var params models.NewWorkflow
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n := &m.Workflow{
		Name:        params.Name,
		Description: params.Description,
		Status:      params.Status,
	}

	result, errs, err := AddWorkflow(n, c.adaptors, c.core)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
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

	workflow, err := GetWorkflowById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(workflow).Send(ctx)
}

// Workflow godoc
// swagger:operation PUT /workflow/{id} workflowUpdateById
// ---
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update workflow params
//   in: body
//   name: workflow
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateWorkflow'
//     type: object
// summary: update workflow by id
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
func (c ControllerWorkflow) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateWorkflow{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	result, errs, err := UpdateWorkflow(params, c.adaptors, c.core)
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
	items, total, err := GetWorkflowList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
	return
}

// swagger:operation DELETE /workflow/{id} workflowDeleteById
// ---
// parameters:
// - description: Workflow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete workflow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow
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
func (c ControllerWorkflow) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteWorkflowById(int64(aid), c.adaptors, c.core); err != nil {
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

// swagger:operation GET /workflows/search workflowSearch
// ---
// summary: search workflow
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - workflow
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
//	   $ref: '#/responses/WorkflowSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerWorkflow) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	workflows, _, err := SearchWorkflow(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("workflows", workflows)
	resp.Send(ctx)
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

	err = UpdateWorkflowScenario(int64(aid), workflowScenario.WorkflowScenarioId, c.adaptors, c.core)
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