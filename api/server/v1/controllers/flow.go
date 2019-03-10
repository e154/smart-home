package controllers

import (
	"strconv"
	"github.com/gin-gonic/gin"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"github.com/e154/smart-home/api/server/v1/models"
)

type ControllerFlow struct {
	*ControllerCommon
}

func NewControllerFlow(common *ControllerCommon) *ControllerFlow {
	return &ControllerFlow{ControllerCommon: common}
}

// Flow godoc
// @tags flow
// @Summary Add new flow
// @Description
// @Produce json
// @Accept  json
// @Param flow body models.NewFlowModel true "flow params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow [post]
// @Security ApiKeyAuth
func (c ControllerFlow) Add(ctx *gin.Context) {

	flow := &models.NewFlowModel{}
	if err := ctx.ShouldBindJSON(&flow); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	id, errs, err := AddFlow(flow, c.adaptors, c.core)
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

// Flow godoc
// @tags flow
// @Summary Show flow
// @Description Get flow by id
// @Produce json
// @Accept  json
// @Param id path int true "Flow ID"
// @Success 200 {object} models.ResponseFlow
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerFlow) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	flow, err := GetFlowById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("flow", flow).Send(ctx)
}

// Flow godoc
// @tags flow
// @Summary Show flow
// @Description Get flow by id
// @Produce json
// @Accept  json
// @Param id path int true "Flow ID"
// @Success 200 {object} models.ResponseRedactorFlowModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow/{id}/redactor [Get]
// @Security ApiKeyAuth
func (c ControllerFlow) GetRedactor(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	flow, err := GetFlowRedactor(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("flow", flow).Send(ctx)
}

// Flow godoc
// @tags flow
// @Summary Update flow
// @Description Update flow by id
// @Produce json
// @Accept  json
// @Param  id path int true "Flow ID"
// @Param  flow body models.RedactorFlowModel true "Update flow"
// @Success 200 {object} models.ResponseRedactorFlowModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow/{id}/redactor [Put]
// @Security ApiKeyAuth
func (c ControllerFlow) UpdateRedactor(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	redactor := &models.RedactorFlowModel{}
	if err := ctx.ShouldBindJSON(&redactor); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	redactor.Id = int64(aid)

	result, errs, err := UpdateFlowRedactor(redactor, c.adaptors, c.core)
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
	resp.Item("flow", result).Send(ctx)
}

// Flow godoc
// @tags flow
// @Summary Update flow
// @Description Update flow by id
// @Produce json
// @Accept  json
// @Param  id path int true "Flow ID"
// @Param  flow body models.UpdateFlowModel true "Update flow"
// @Success 200 {object} models.ResponseFlow
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow/{id} [Put]
// @Security ApiKeyAuth
func (c ControllerFlow) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	flow := &models.UpdateFlowModel{}
	if err := ctx.ShouldBindJSON(&flow); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	flow.Id = int64(aid)

	result, errs, err := UpdateFlow(flow, c.adaptors, c.core)
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
	resp.Item("flow", result).Send(ctx)
}

// Flow godoc
// @tags flow
// @Summary Flow list
// @Description Get flow list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Success 200 {object} models.ResponseFlowList
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flows [Get]
// @Security ApiKeyAuth
func (c ControllerFlow) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetFlowList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
}

// Flow godoc
// @tags flow
// @Summary Delete flow
// @Description Delete flow by id
// @Produce json
// @Accept  json
// @Param  id path int true "Flow ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerFlow) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err = DeleteFlowById(int64(aid), c.adaptors, c.core); err != nil {
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

// Flow godoc
// @tags flow
// @Summary Search flow
// @Description Search flow by name
// @Produce json
// @Accept  json
// @Param query query string false "query"
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Success 200 {object} models.ResponseSearchFlow
// @Failure 401 "Unauthorized"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /flows/search [Get]
func (c ControllerFlow) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	flows, _, err := SearchFlow(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("flows", flows)
	resp.Send(ctx)
}
