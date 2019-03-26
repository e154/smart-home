package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
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

	mapLayer := &m.MapLayer{}
	common.Copy(&mapLayer, &params)

	if params.Map != nil && params.Map.Id != 0 {
		mapLayer.MapId = params.Map.Id
	}

	mapLayer, errs, err := c.command.MapLayer.Add(mapLayer)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.MapLayer{}
	common.Copy(&result, &mapLayer)

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

	mapLayer, err := c.command.MapLayer.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.MapLayer{}
	common.Copy(&result, &mapLayer)

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

	mapLayer := &m.MapLayer{}
	common.Copy(&mapLayer, &params)

	if params.Map != nil && params.Map.Id != 0 {
		mapLayer.MapId = params.Map.Id
	}

	mapLayer, errs, err := c.command.MapLayer.Update(mapLayer)
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

	result := &models.MapLayer{}
	common.Copy(&result, &mapLayer)

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

	params := make([]*m.SortMapLayer, 0)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.command.MapLayer.Sort(params); err != nil {
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
	items, total, err := c.command.MapLayer.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MapLayer, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
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

	if err := c.command.MapLayer.Delete(int64(aid)); err != nil {
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
