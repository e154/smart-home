package controllers

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

type ControllerMapElement struct {
	*ControllerCommon
}

func NewControllerMapElement(common *ControllerCommon) *ControllerMapElement {
	return &ControllerMapElement{ControllerCommon: common}
}

// swagger:operation POST /map_element mapElementAdd
// ---
// parameters:
// - description: element params
//   in: body
//   name: map_element
//   required: true
//   schema:
//     $ref: '#/definitions/NewMapElement'
//     type: object
// summary: add new map element
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_element
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapElement'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMapElement) Add(ctx *gin.Context) {

	params := &models.NewMapElement{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	mapElement := &m.MapElement{}
	common.Copy(&mapElement, &params, common.JsonEngine)

	mapElement, errs, err := c.command.MapElement.Add(mapElement)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.MapElement{}
	common.Copy(&result, &mapElement, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /map_element/{id} mapElementGetById
// ---
// parameters:
// - description: MapElement ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get map element by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_element
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapElement'
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
func (c ControllerMapElement) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	mapElement, err := c.command.MapElement.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.MapElement{}
	common.Copy(&result, &mapElement, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /map_element/{id} mapElementUpdateById
// ---
// parameters:
// - description: MapElement ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update map element params
//   in: body
//   name: map_element
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateMapElement'
//     type: object
// summary: full update map element by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_element
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapElement'
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
func (c ControllerMapElement) UpdateFull(ctx *gin.Context) {

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

	mapElement := &m.MapElement{}
	common.Copy(&mapElement, &params, common.JsonEngine)

	mapElement, errs, err := c.command.MapElement.Update(mapElement)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	result := &models.MapElement{}
	common.Copy(&result, &mapElement, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /map_element/{id}/element_only mapElementUpdateById
// ---
// parameters:
// - description: MapElement ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update map element params
//   in: body
//   name: map_element
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateMapElement'
//     type: object
// summary: update map element only by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_element
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapElement'
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
func (c ControllerMapElement) UpdateElement(ctx *gin.Context) {

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

	mapElement := &m.MapElement{}
	common.Copy(&mapElement, &params, common.JsonEngine)

	mapElement, errs, err := c.command.MapElement.UpdateElement(mapElement)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	result := &models.MapElement{}
	common.Copy(&result, &mapElement, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /map_element/sort mapElementUpdateById
// ---
// parameters:
// - description: sort params
//   in: body
//   name: sort_params
//   required: true
//   schema:
//     $ref: '#/definitions/SortMapElement'
//     type: object
// summary: sort map elements
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_element
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
func (c ControllerMapElement) Sort(ctx *gin.Context) {

	params := make([]*m.SortMapElement, 0)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.command.MapElement.Sort(params); err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation GET /map_elements mapElementList
// ---
// summary: get map element list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_element
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
//     description: OK
//     schema:
//       $ref: '#/definitions/MapElement'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMapElement) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.command.MapElement.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MapElement, 0)
	common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /map_element/{id} mapElementDeleteById
// ---
// parameters:
// - description: MapElement ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete map element by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_element
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
func (c ControllerMapElement) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.command.MapElement.Delete(int64(aid)); err != nil {
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
