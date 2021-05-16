// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

// ControllerAlexa ...
type ControllerAlexa struct {
	*ControllerCommon
}

// NewControllerAlexa ...
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
//     $ref: '#/definitions/NewAlexaSkill'
//     type: object
// summary: add alexa skill
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - alexa
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/AlexaSkill'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAlexa) Add(ctx *gin.Context) {

	params := &models.NewAlexaSkill{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	alexa := &m.AlexaSkill{}
	common.Copy(&alexa, &params, common.JsonEngine)

	alexa, errs, err := c.endpoint.AlexaSkill.Add(alexa)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.AlexaSkill{}
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
// summary: get alexa skill by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - alexa
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/AlexaSkill'
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

	alexa, err := c.endpoint.AlexaSkill.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.AlexaSkill{}
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
// - description: Update alexa skill
//   in: body
//   name: alexa
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateAlexaSkill'
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
//       $ref: '#/definitions/AlexaSkill'
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

	params := &models.UpdateAlexaSkill{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	alexa := &m.AlexaSkill{}
	common.Copy(&alexa, &params, common.JsonEngine)

	alexa, errs, err := c.endpoint.AlexaSkill.Update(alexa)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.AlexaSkill{}
	common.Copy(&result, &alexa, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /alexas alexaList
// ---
// summary: get alexa skill list
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
//	   $ref: '#/responses/AlexaSkillList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAlexa) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.AlexaSkill.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.AlexaSkillShort, 0)
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
// summary: delete alexa skill by id
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

	if err := c.endpoint.AlexaSkill.Delete(int64(aid)); err != nil {
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
