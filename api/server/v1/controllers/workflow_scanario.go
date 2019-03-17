package controllers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/copier"
)

type ControllerWorkflowScenario struct {
	*ControllerCommon
}

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
//	   $ref: '#/responses/NewObjectSuccess'
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

	_, id, errs, err := AddWorkflowScenario(scenario, c.adaptors)
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
	resp.Item("id", id).Send(ctx)
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

	workflowScenario, err := GetWorkflowScenarioById(int64(workflowId), int64(scenarioId), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(workflowScenario).Send(ctx)
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

	workflowScenario := &m.WorkflowScenario{
		Id:         params.Id,
		Name:       params.Name,
		WorkflowId: int64(workflowId),
		SystemName: params.SystemName,
		Scripts:    make([]*m.Script, 0),
	}

	for _, s := range params.Scripts {
		script := &m.Script{}
		if err = copier.Copy(&script, &s); err != nil {
			log.Error(err.Error())
		}
		workflowScenario.Scripts = append(workflowScenario.Scripts, script)
	}

	_, errs, err := WorkflowUpdateWorkflowScenario(workflowScenario, c.adaptors)
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
	resp.Send(ctx)
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

	items, _, err := GetWorkflowScenarioList(int64(workflowId), c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("scenarios", items).Send(ctx)
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

	if err := DeleteWorkflowScenarioById(int64(workflowScenarioId), c.adaptors); err != nil {
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
	scenarios, _, err := SearchWorkflowScenario(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("scenarios", scenarios)
	resp.Send(ctx)
}
