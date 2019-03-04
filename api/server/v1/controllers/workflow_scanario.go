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

// WorkflowScenario godoc
// @tags workflow_scenario
// @Summary Add new workflow_scenario
// @Description
// @Produce json
// @Accept  json
// @Param workflow_scenario body models.NewWorkflowScenario true "workflow_scenario params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow/{id}/scenario [post]
// @Security ApiKeyAuth
func (c ControllerWorkflowScenario) AddWorkflowScenario(ctx *gin.Context) {

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

// WorkflowScenario godoc
// @tags workflow_scenario
// @Summary Show workflow_scenario
// @Description Get workflow_scenario by id
// @Produce json
// @Accept  json
// @Param id path int true "WorkflowScenario ID"
// @Success 200 {object} models.WorkflowScenario
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow/{id}/scenario/{scenario_id} [Get]
// @Security ApiKeyAuth
func (c ControllerWorkflowScenario) GetWorkflowScenarioById(ctx *gin.Context) {

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

// WorkflowScenario godoc
// @tags workflow_scenario
// @Summary Update workflow_scenario
// @Description Update workflow_scenario by id
// @Produce json
// @Accept  json
// @Param  id path int true "WorkflowScenario ID"
// @Param  workflow_scenario body models.WorkflowUpdateWorkflowScenario true "Update workflow_scenario"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow/{id}/scenario/{scenario_id} [Put]
// @Security ApiKeyAuth
func (c ControllerWorkflowScenario) UpdateWorkflowScenario(ctx *gin.Context) {

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

// WorkflowScenario godoc
// @tags workflow_scenario
// @Summary WorkflowScenario list
// @Description Get workflow_scenario list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Success 200 {object} models.WorkflowScenarioListModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow/{id}/scenarios [Get]
// @Security ApiKeyAuth
func (c ControllerWorkflowScenario) GetWorkflowScenarioList(ctx *gin.Context) {

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

// WorkflowScenario godoc
// @tags workflow_scenario
// @Summary Delete workflow_scenario
// @Description Delete workflow_scenario by id
// @Produce json
// @Accept  json
// @Param  id path int true "WorkflowScenario ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow/{id}/scenario/{scenario_id} [Delete]
// @Security ApiKeyAuth
func (c ControllerWorkflowScenario) DeleteWorkflowScenarioById(ctx *gin.Context) {

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

// WorkflowScenario godoc
// @tags workflow
// @Summary Search workflow
// @Description Search workflow by name
// @Produce json
// @Accept  json
// @Param query query string false "query"
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Success 200 {object} models.SearchWorkflowScenarioResponse
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /workflow/{id}/scenarios/search [Get]
func (c ControllerWorkflowScenario) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	workflows, _, err := SearchWorkflowScenario(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("workflows", workflows)
	resp.Send(ctx)
}
