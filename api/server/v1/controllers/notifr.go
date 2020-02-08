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
	"github.com/e154/smart-home/system/notify"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ControllerNotifr struct {
	*ControllerCommon
}

func NewControllerNotifr(common *ControllerCommon) *ControllerNotifr {
	return &ControllerNotifr{ControllerCommon: common}
}

// swagger:operation PUT /notifr/config notifyUpdateSettings
// ---
// parameters:
// - description: Update notifr params
//   in: body
//   name: notifr
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateNotifrConfig'
//     type: object
// summary: update notifr settings
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notifr
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
func (c ControllerNotifr) Update(ctx *gin.Context) {

	params := &models.UpdateNotifrConfig{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	settings := &notify.NotifyConfig{}
	_ = common.Copy(&settings, &params, common.JsonEngine)

	err := c.endpoint.Notify.UpdateSettings(settings)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation GET /notifr/config notifyGetSettings
// ---
// parameters:
// summary: get notifr settings
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notifr
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/NotifrConfig'
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
func (c ControllerNotifr) GetSettings(ctx *gin.Context) {

	settings, err := c.endpoint.Notify.GetSettings()
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.NotifrConfig{}
	_ = common.Copy(&result, &settings, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /notifrs notifrList
// ---
// summary: get notification list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notifr
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
//	   $ref: '#/responses/MessageDeliveryList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerNotifr) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	notifications, total, err := c.endpoint.MessageDelivery.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MessageDelivery, 0)
	_ = common.Copy(&result, &notifications)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /notifr/{id} notifrDeleteById
// ---
// parameters:
// - description: notification ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete notification by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notifr
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
func (c ControllerNotifr) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.MessageDelivery.Delete(int64(aid)); err != nil {
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

// swagger:operation POST /notifr/{id}/repeat notifyRepeatMessage
// ---
// parameters:
// - description: notification ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: repeat notification by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notifr
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
func (c ControllerNotifr) Repeat(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err = c.endpoint.Notify.Repeat(int64(aid)); err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation POST /notifr notifySendNewMessage
// ---
// parameters:
// - description: message object
//   in: body
//   required: true
//   schema:
//     $ref: '#/definitions/NewNotifrMessage'
//     type: object
// summary: repeat notification by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notifr
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
func (c ControllerNotifr) Send(ctx *gin.Context) {

	params := &models.NewNotifrMessage{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	message := &m.NewNotifrMessage{}
	_ = common.Copy(&message, &params, common.JsonEngine)

	err := c.endpoint.Notify.Send(message)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}
