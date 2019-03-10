package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerMapLayer struct {
	*ControllerCommon
}

func NewControllerMapLayer(common *ControllerCommon) *ControllerMapLayer {
	return &ControllerMapLayer{ControllerCommon: common}
}

// MapLayer godoc
// @tags map_layer
// @Summary Add new map_layer
// @Description
// @Produce json
// @Accept  json
// @Param map_layer body models.NewMapLayer true "map_layer params"
// @Success 200 {object} models.MapLayer
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_layer [post]
// @Security ApiKeyAuth
func (c ControllerMapLayer) Add(ctx *gin.Context) {

	params := &models.NewMapLayer{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, _, errs, err := AddMapLayer(params, c.adaptors)
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
	resp.Item("map_layer", result).Send(ctx)
}

// MapLayer godoc
// @tags map_layer
// @Summary Show map_layer
// @Description Get map_layer by id
// @Produce json
// @Accept  json
// @Param id path int true "MapLayer ID"
// @Success 200 {object} models.MapLayer
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_layer/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerMapLayer) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, err := GetMapLayerById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("map_layer", result).Send(ctx)
}

// MapLayer godoc
// @tags map_layer
// @Summary Update map_layer
// @Description Update map_layer by id
// @Produce json
// @Accept  json
// @Param  id path int true "MapLayer ID"
// @Param  map_layer body models.UpdateMapLayer true "Update map_layer"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_layer/{id} [Put]
// @Security ApiKeyAuth
func (c ControllerMapLayer) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateMapLayer{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	_, errs, err := UpdateMapLayer(params, c.adaptors)
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

// MapLayer godoc
// @tags map_layer
// @Summary Sort map_layer
// @Description Sort map_layer by id
// @Produce json
// @Accept  json
// @Param  map_layer body models.SortMapLayer true "Sort map_layer"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_layer/{id} [Put]
// @Security ApiKeyAuth
func (c ControllerMapLayer) Sort(ctx *gin.Context) {

	params := make([]*models.SortMapLayer, 0)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := SortMapLayers(params, c.adaptors); err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// MapLayer godoc
// @tags map_layer
// @Summary MapLayer list
// @Description Get map_layer list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Success 200 {object} models.MapListModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_layers [Get]
// @Security ApiKeyAuth
func (c ControllerMapLayer) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetMapLayerList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
	return
}

// MapLayer godoc
// @tags map_layer
// @Summary Delete map_layer
// @Description Delete map_layer by id
// @Produce json
// @Accept  json
// @Param  id path int true "MapLayer ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_layer/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerMapLayer) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteMapLayer(int64(aid), c.adaptors); err != nil {
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
