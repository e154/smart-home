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
// summary: update map element by id
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
//	   $ref: '#/responses/MapElementList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
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
