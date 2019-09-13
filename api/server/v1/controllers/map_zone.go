package controllers

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/gin-gonic/gin"
)

type ControllerMapZone struct {
	*ControllerCommon
}

func NewControllerMapZone(common *ControllerCommon) *ControllerMapZone {
	return &ControllerMapZone{ControllerCommon: common}
}

// swagger:operation POST /map_zone mapZoneAdd
// ---
// parameters:
// - description: map_zone params
//   in: body
//   name: map_zone
//   required: true
//   schema:
//     $ref: '#/definitions/NewMapZone'
//     type: object
// summary: add new map_zone
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_zone
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MapZone'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMapZone) Add(ctx *gin.Context) {

	params := &models.NewMapZone{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	zone := &m.MapZone{}
	_ = common.Copy(&zone, &params, common.JsonEngine)

	zone, errs, err := c.endpoint.MapZone.Add(zone)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.MapZone{}
	if err = common.Copy(&result, &zone, common.JsonEngine); err != nil {
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation DELETE /map_zone/{name} mapZoneDeleteByName
// ---
// parameters:
// - description: MapZone Name
//   in: path
//   name: name
//   required: true
//   type: text
// summary: delete map_zone by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_zone
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
func (c ControllerMapZone) Delete(ctx *gin.Context) {

	name := ctx.Param("name")
	if err := c.endpoint.MapZone.Delete(name); err != nil {
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

// swagger:operation GET /map_zone/search mapZoneSearch
// ---
// summary: search map_zone
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_zone
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
//	   $ref: '#/responses/MapZoneSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMapZone) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.MapZone.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MapZone, 0)
	_ = common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Item("zones", result)
	resp.Send(ctx)
}
