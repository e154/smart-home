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
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

type ControllerAlexa struct {
	*ControllerCommon
}

func NewControllerAlexa(common *ControllerCommon) *ControllerAlexa {
	return &ControllerAlexa{ControllerCommon: common}
}

// swagger:operation POST /alexa alexaAdd
// ---
// parameters:
// - description: alexa params
//   in: body
//   name: alexa
//   required: true
//   schema:
//     $ref: '#/definitions/NewAlexaApplication'
//     type: object
// summary: add new alexa
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - alexa
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/AlexaApplication'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAlexa) Add(ctx *gin.Context) {

	params := &models.NewAlexaApplication{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	alexa := &m.AlexaApplication{}
	common.Copy(&alexa, &params, common.JsonEngine)

	alexa, errs, err := c.endpoint.AlexaApplication.Add(alexa)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.AlexaApplication{}
	if err = common.Copy(&result, &alexa, common.JsonEngine); err != nil {
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /alexa/{id} alexaGetById
// ---
// parameters:
// - description: Alexa ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get alexa by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - alexa
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/AlexaApplication'
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
func (c ControllerAlexa) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	alexa, err := c.endpoint.AlexaApplication.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.AlexaApplication{}
	common.Copy(&result, &alexa, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /alexa/{id} alexaUpdateById
// ---
// parameters:
// - description: Alexa ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update alexa params
//   in: body
//   name: alexa
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateAlexaApplication'
//     type: object
// summary: update alexa by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - alexa
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/AlexaApplication'
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
func (c ControllerAlexa) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateAlexaApplication{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	alexa := &m.AlexaApplication{}
	common.Copy(&alexa, &params, common.JsonEngine)

	alexa, errs, err := c.endpoint.AlexaApplication.Update(alexa)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.AlexaApplication{}
	common.Copy(&result, &alexa, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /alexas alexaList
// ---
// summary: get alexa list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - alexa
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
//	   $ref: '#/responses/AlexaApplicationList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAlexa) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.AlexaApplication.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.AlexaApplicationShort, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /alexa/{id} alexaDeleteById
// ---
// parameters:
// - description: Alexa ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete alexa by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - alexa
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
func (c ControllerAlexa) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.AlexaApplication.Delete(int64(aid)); err != nil {
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
