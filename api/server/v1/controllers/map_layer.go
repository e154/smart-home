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

// swagger:operation POST /map_layer mapLayerAdd
// ---
// parameters:
// - description: layer params
//   in: body
//   name: map_layer
//   required: true
//   schema:
//     $ref: '#/definitions/NewMapLayer'
//     type: object
// summary: add new map layer
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_layer
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapLayer'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
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
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /map_layer/{id} mapLayerGetById
// ---
// parameters:
// - description: MapLayer ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get map layer by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_layer
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapLayer'
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
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /map_layer/{id} mapLayerUpdateById
// ---
// parameters:
// - description: MapLayer ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: layer params
//   in: body
//   name: map_layer
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateMapLayer'
//     type: object
// summary: update map layer
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_layer
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapLayer'
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

	result, errs, err := UpdateMapLayer(params, c.adaptors)
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

// swagger:operation PUT /map_layer/sort mapLayerUpdateById
// ---
// parameters:
// - description: sort params
//   in: body
//   name: sort_params
//   required: true
//   schema:
//     $ref: '#/definitions/SortMapLayer'
//     type: object
// summary: sort map layers
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_layer
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

// swagger:operation GET /map_layers mapLayerList
// ---
// summary: get map layer list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_layer
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
//	   $ref: '#/responses/MapLayerList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
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

// swagger:operation DELETE /map_layer/{id} mapLayerDeleteById
// ---
// parameters:
// - description: MapLayer ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete map layer by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_layer
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
