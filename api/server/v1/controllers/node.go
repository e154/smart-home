package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerNode struct {
	*ControllerCommon
}

func NewControllerNode(common *ControllerCommon) *ControllerNode {
	return &ControllerNode{ControllerCommon: common}
}

// Node godoc
// @tags node
// @Summary Add new node
// @Description
// @Produce json
// @Accept  json
// @Param node body models.NewNode true "node params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /node [post]
// @Security ApiKeyAuth
func (c ControllerNode) AddNode(ctx *gin.Context) {

	params := &models.NewNode{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	_, id, errs, err := AddNode(params, c.adaptors, c.core)
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

// Node godoc
// @tags node
// @Summary Show node
// @Description Get node by id
// @Produce json
// @Accept  json
// @Param id path int true "Node ID"
// @Success 200 {object} models.NodeModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /node/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerNode) GetNodeById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	node, err := GetNodeById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(node).Send(ctx)
}

// Node godoc
// @tags node
// @Summary Update node
// @Description Update node by id
// @Produce json
// @Accept  json
// @Param  id path int true "Node ID"
// @Param  node body models.UpdateNodeModel true "Update node"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /node/{id} [Put]
// @Security ApiKeyAuth
func (c ControllerNode) UpdateNode(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n := &models.UpdateNodeModel{}
	if err := ctx.ShouldBindJSON(&n); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n.Id = int64(aid)

	_, errs, err := UpdateNode(n, c.adaptors, c.core)
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

// Node godoc
// @tags node
// @Summary Node list
// @Description Get node list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Success 200 {object} models.NodeListModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /nodes [Get]
// @Security ApiKeyAuth
func (c ControllerNode) GetNodeList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetNodeList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
	return
}

// Node godoc
// @tags node
// @Summary Delete node
// @Description Delete node by id
// @Produce json
// @Accept  json
// @Param  id path int true "Node ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /node/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerNode) DeleteNodeById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteNodeById(int64(aid), c.adaptors, c.core); err != nil {
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

// NodeModel godoc
// @tags node
// @Summary Search node
// @Description Search node by name
// @Produce json
// @Accept  json
// @Param query query string false "query"
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Success 200 {object} models.SearchNodeResponse
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /nodes/search [Get]
func (c ControllerNode) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	nodes, _, err := SearchNode(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("nodes", nodes)
	resp.Send(ctx)
}
