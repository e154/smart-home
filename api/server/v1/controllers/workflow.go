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

// Workflow godoc
// @tags workflow
// @Summary Add new workflow
// @Description
// @Produce json
// @Accept  json
// @Param workflow body models.NewWorkflow true "workflow params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow [post]
func (c ControllerWorkflow) AddWorkflow(ctx *gin.Context) {

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

	_, id, errs, err := AddWorkflow(n, c.adaptors, c.core)
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

// Workflow godoc
// @tags workflow
// @Summary Show workflow
// @Description Get workflow by id
// @Produce json
// @Accept  json
// @Param id path int true "Workflow ID"
// @Success 200 {object} models.ResponseWorkflow
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow/{id} [Get]
func (c ControllerWorkflow) GetWorkflowById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
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
	resp.Item("workflow", workflow).Send(ctx)
}

// Workflow godoc
// @tags workflow
// @Summary Update workflow
// @Description Update workflow by id
// @Produce json
// @Accept  json
// @Param  id path int true "Workflow ID"
// @Param  workflow body models.UpdateWorkflow true "Update workflow"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow/{id} [Put]
func (c ControllerWorkflow) UpdateWorkflow(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n := &m.Workflow{}
	if err := ctx.ShouldBindJSON(&n); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n.Id = int64(aid)

	_, errs, err := UpdateWorkflow(n, c.adaptors, c.core)
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

// Workflow godoc
// @tags workflow
// @Summary Workflow list
// @Description Get workflow list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Success 200 {object} models.ResponseWorkflowList
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow [Get]
func (c ControllerWorkflow) GetWorkflowList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetWorkflowList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, int(total), items).Send(ctx)
	return
}

// Workflow godoc
// @tags workflow
// @Summary Delete workflow
// @Description Delete workflow by id
// @Produce json
// @Accept  json
// @Param  id path int true "Workflow ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /workflow/{id} [Delete]
func (c ControllerWorkflow) DeleteWorkflowById(ctx *gin.Context) {

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
