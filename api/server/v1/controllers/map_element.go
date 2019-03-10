package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerMapElement struct {
	*ControllerCommon
}

func NewControllerMapElement(common *ControllerCommon) *ControllerMapElement{
	return &ControllerMapElement{ControllerCommon: common}
}

// MapElement godoc
// @tags map_element
// @Summary Add new map_element
// @Description
// @Produce json
// @Accept  json
// @Param map_element body models.NewMapElement true "map_element params"
// @Success 200 {object} models.MapElement
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_element [post]
// @Security ApiKeyAuth
func (c ControllerMapElement) Add(ctx *gin.Context) {

	params := &models.NewMapElement{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, _, errs, err := AddMapElement(params, c.adaptors)
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
	resp.Item("map_element", result).Send(ctx)
}

// MapElement godoc
// @tags map_element
// @Summary Show map_element
// @Description Get map_element by id
// @Produce json
// @Accept  json
// @Param id path int true "MapElement ID"
// @Success 200 {object} models.MapElement
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_element/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerMapElement) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, err := GetMapElementById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("map_element", result).Send(ctx)
}

// MapElement godoc
// @tags map_element
// @Summary Update map_element
// @Description Update map_element by id
// @Produce json
// @Accept  json
// @Param  id path int true "MapElement ID"
// @Param  map_element body models.UpdateMapElement true "Update map_element"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_element/{id} [Put]
// @Security ApiKeyAuth
func (c ControllerMapElement) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateMapElement{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	_, errs, err := UpdateMapElement(params, c.adaptors)
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

// MapElement godoc
// @tags map_element
// @Summary Sort map_element
// @Description Sort map_element by id
// @Produce json
// @Accept  json
// @Param  map_element body models.SortMapElement true "Sort map_element"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_element/{id} [Put]
// @Security ApiKeyAuth
func (c ControllerMapElement) Sort(ctx *gin.Context) {

	params := make([]*models.SortMapElement, 0)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := SortMapElements(params, c.adaptors); err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// MapElement godoc
// @tags map_element
// @Summary MapElement list
// @Description Get map_element list
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
// @Router /map_elements [Get]
// @Security ApiKeyAuth
func (c ControllerMapElement) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetMapElementList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
	return
}

// MapElement godoc
// @tags map_element
// @Summary Delete map_element
// @Description Delete map_element by id
// @Produce json
// @Accept  json
// @Param  id path int true "MapElement ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /map_element/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerMapElement) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteMapElement(int64(aid), c.adaptors); err != nil {
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
