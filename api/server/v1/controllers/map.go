package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerMap struct {
	*ControllerCommon
}

func NewControllerMap(common *ControllerCommon) *ControllerMap {
	return &ControllerMap{ControllerCommon: common}
}

// swagger:operation POST /map mapAdd
// ---
// parameters:
// - description: map params
//   in: body
//   name: map
//   required: true
//   schema:
//     $ref: '#/definitions/NewMap'
//     type: object
// summary: add new map
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Map'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMap) Add(ctx *gin.Context) {

	params := &models.NewMap{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, _, errs, err := AddMap(params, c.adaptors, c.core)
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

// swagger:operation GET /map/{id} mapGetById
// ---
// parameters:
// - description: Map ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get map by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Map'
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
func (c ControllerMap) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	m, err := GetMapById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(m).Send(ctx)
}

// swagger:operation GET /map/{id}/full mapFullGetById
// ---
// parameters:
// - description: Map ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get map full by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapFull'
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
func (c ControllerMap) GetFullMap(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	m, err := GetFullMapById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(m).Send(ctx)
}

// swagger:operation PUT /map/{id} mapUpdateById
// ---
// parameters:
// - description: Map ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update map params
//   in: body
//   name: map
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateMap'
//     type: object
// summary: update map by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Map'
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
func (c ControllerMap) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateMap{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	result, errs, err := UpdateMap(params, c.adaptors, c.core)
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

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /maps mapList
// ---
// summary: get map list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
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
//	   $ref: '#/responses/MapList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMap) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetMapList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
	return
}

// swagger:operation DELETE /map/{id} mapDeleteById
// ---
// parameters:
// - description: Map ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete map by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
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
func (c ControllerMap) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteMapById(int64(aid), c.adaptors, c.core); err != nil {
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

// swagger:operation GET /maps/search mapSearch
// ---
// summary: search map
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
// parameters:
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
//	   $ref: '#/responses/MapSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMap) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	maps, _, err := SearchMap(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("maps", maps)
	resp.Send(ctx)
}
