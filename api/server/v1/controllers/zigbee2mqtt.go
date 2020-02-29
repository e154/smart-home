// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package controllers

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ControllerZigbee2mqtt struct {
	*ControllerCommon
}

func NewControllerZigbee2mqtt(common *ControllerCommon) *ControllerZigbee2mqtt {
	return &ControllerZigbee2mqtt{ControllerCommon: common}
}

// swagger:operation POST /zigbee2mqtt bridgeAdd
// ---
// parameters:
// - description: zigbee2mqtt params
//   in: body
//   name: zigbee2mqtt
//   required: true
//   schema:
//     $ref: '#/definitions/NewZigbee2mqtt'
//     type: object
// summary: add new zigbee2mqtt
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Zigbee2mqtt'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerZigbee2mqtt) Add(ctx *gin.Context) {

	params := &models.NewZigbee2mqtt{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if params.Password != params.PasswordRepeat {
		NewError(400, "bad password repeat").Send(ctx)
		return
	}

	bridge := &m.Zigbee2mqtt{}
	common.Copy(&bridge, &params, common.JsonEngine)

	bridge, errs, err := c.endpoint.Zigbee2mqtt.Add(bridge)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Zigbee2mqtt{}
	if err = common.Copy(&result, &bridge, common.JsonEngine); err != nil {
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /zigbee2mqtt/{id} bridgeGetById
// ---
// parameters:
// - description: get zigbee2mqtt bridge by ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get zigbee2mqtt bridge by ID
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Zigbee2mqttInfo'
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
func (c ControllerZigbee2mqtt) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	bridge, err := c.endpoint.Zigbee2mqtt.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Zigbee2mqttInfo{}
	common.Copy(&result, &bridge, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /zigbee2mqtt/{id} bridgeUpdateById
// ---
// parameters:
// - description: Zigbee2mqtt ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update bridge params
//   in: body
//   name: bridge
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateZigbee2mqtt'
//     type: object
// summary: update bridge by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Zigbee2mqtt'
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
func (c ControllerZigbee2mqtt) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateZigbee2mqtt{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if params.Password != nil && params.PasswordRepeat != nil {
		if *params.Password != *params.PasswordRepeat {
			NewError(400, "bad password repeat").Send(ctx)
			return
		}
	}

	bridge := &m.Zigbee2mqtt{}
	common.Copy(&bridge, &params, common.JsonEngine)
	bridge.Id = int64(aid)

	bridge, errs, err := c.endpoint.Zigbee2mqtt.Update(bridge)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Zigbee2mqtt{}
	common.Copy(&result, &bridge, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /zigbee2mqtts bridgeList
// ---
// summary: get bridge list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
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
//	   $ref: '#/responses/Zigbee2mqttList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerZigbee2mqtt) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Zigbee2mqtt.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Zigbee2mqttInfo, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /zigbee2mqtt/{id} bridgeDeleteById
// ---
// parameters:
// - description: Zigbee2mqtt ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete bridge by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
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
func (c ControllerZigbee2mqtt) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.Zigbee2mqtt.Delete(int64(aid)); err != nil {
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

// swagger:operation POST /zigbee2mqtt/{id}/reset bridgeResetById
// ---
// parameters:
// - description: reset zigbee2mqtt bridge by ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: reset zigbee2mqtt bridge by ID
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
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
func (c ControllerZigbee2mqtt) Reset(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.Zigbee2mqtt.ResetBridge(int64(aid)); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
}

// swagger:operation POST /zigbee2mqtt/{id}/device_ban banDevice
// ---
// parameters:
// - description: set device to ban
//   in: path
//   name: id
//   required: true
//   type: integer
// - name: friendly_name
//   in: query
//   description: device id
//   required: true
//   type: text
// summary: set device to ban
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
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
func (c ControllerZigbee2mqtt) DeviceBan(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	friendlyName := ctx.Param("friendly_name")
	if friendlyName == "" {
		NewError(400, "friendly_name is required param").Send(ctx)
		return
	}

	if err := c.endpoint.Zigbee2mqtt.DeviceBan(int64(aid), friendlyName); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
}

// swagger:operation POST /zigbee2mqtt/{id}/device_whitelist deviceWhitelist
// ---
// parameters:
// - description: set device by id to white list
//   in: path
//   name: id
//   required: true
//   type: integer
// - name: friendly_name
//   in: query
//   description: device id
//   required: true
//   type: text
// summary: set device by id to white list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
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
func (c ControllerZigbee2mqtt) DeviceWhitelist(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	friendlyName := ctx.Param("friendly_name")
	if friendlyName == "" {
		NewError(400, "friendly_name is required param").Send(ctx)
		return
	}

	if err := c.endpoint.Zigbee2mqtt.DeviceWhitelist(int64(aid), friendlyName); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
}

// swagger:operation GET /zigbee2mqtt/{id}/networkmap Networkmap
// ---
// parameters:
// - description: get network map
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get network map
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
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
func (c ControllerZigbee2mqtt) Networkmap(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	var networkmap string
	if networkmap, err = c.endpoint.Zigbee2mqtt.Networkmap(int64(aid)); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	ctx.String(200, networkmap)
}

// swagger:operation POST /zigbee2mqtt/{id}/update_networkmap UpdateNetworkmap
// ---
// parameters:
// - description: update network map
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: reset update network map
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
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
func (c ControllerZigbee2mqtt) UpdateNetworkmap(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.Zigbee2mqtt.UpdateNetworkmap(int64(aid)); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
}

// swagger:operation PATCH /zigbee2mqtts/device_rename deviceRename
// ---
// parameters:
// - name: friendly_name
//   in: query
//   description: device id
//   required: true
//   type: string
// - name: name
//   in: query
//   description: new name
//   required: true
//   type: string
// summary: rename device
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
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
func (c ControllerZigbee2mqtt) DeviceRename(ctx *gin.Context) {

	friendlyName := ctx.Query("friendly_name")
	if friendlyName == "" {
		NewError(400, "friendly_name is required param").Send(ctx)
		return
	}

	name := ctx.Query("name")
	if name == "" {
		NewError(400, "name is required param").Send(ctx)
		return
	}

	if err := c.endpoint.Zigbee2mqtt.DeviceRename(friendlyName, name); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
}

// swagger:operation GET /zigbee2mqtts/search_device zigbee2mqttSearchDevice
// ---
// summary: search device by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - zigbee2mqtt
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
//	   $ref: '#/responses/Zigbee2mqttDeviceSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerZigbee2mqtt) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.Zigbee2mqtt.SearchDevice(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Zigbee2mqttDeviceShort, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Item("devices", result)
	resp.Send(ctx)
}
