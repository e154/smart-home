package controllers

import (
	"github.com/gin-gonic/gin"
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
// @Param flow body models.NewFlow true "flow params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow [post]
// @Security ApiKeyAuth
func (c ControllerFlow) AddFlow(ctx *gin.Context) {


	resp := NewSuccess()
	resp.Send(ctx)
}

// Flow godoc
// @tags flow
// @Summary Show flow
// @Description Get flow by id
// @Produce json
// @Accept  json
// @Param id path int true "Flow ID"
// @Success 200 {object} models.Flow
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerFlow) GetFlowById(ctx *gin.Context) {


	resp := NewSuccess()
	resp.Send(ctx)
}

// Flow godoc
// @tags flow
// @Summary Update flow
// @Description Update flow by id
// @Produce json
// @Accept  json
// @Param  id path int true "Flow ID"
// @Param  flow body models.UpdateFlow true "Update flow"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow/{id} [Put]
// @Security ApiKeyAuth
func (c ControllerFlow) UpdateFlow(ctx *gin.Context) {

	resp := NewSuccess()
	resp.Send(ctx)
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
// @Success 200 {object} models.FlowListModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flows [Get]
// @Security ApiKeyAuth
func (c ControllerFlow) GetFlowList(ctx *gin.Context) {

	resp := NewSuccess()
	resp.Send(ctx)
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
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /flow/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerFlow) DeleteFlowById(ctx *gin.Context) {

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
// @Success 200 {object} models.SearchFlowResponse
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /flows/search [Get]
func (c ControllerFlow) Search(ctx *gin.Context) {

	resp := NewSuccess()
	resp.Send(ctx)
}
