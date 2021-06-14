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
	"github.com/gin-gonic/gin"
)

// ControllerPlugin ...
type ControllerPlugin struct {
	*ControllerCommon
}

// NewControllerPlugin ...
func NewControllerPlugin(common *ControllerCommon) *ControllerPlugin {
	return &ControllerPlugin{ControllerCommon: common}
}

// swagger:operation POST /plugin/{name}/enable pluginEnable
// ---
// parameters:
// - description: Plugin name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: enable plugin by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - plugin
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Plugin'
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
func (c ControllerPlugin) Enable(ctx *gin.Context) {

	err := c.endpoint.Plugin.Enabled(ctx.Param("name"))

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation POST /plugin/{name}/disable pluginDisable
// ---
// parameters:
// - description: Plugin name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: disable plugin by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - plugin
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Plugin'
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
func (c ControllerPlugin) Disable(ctx *gin.Context) {

	err := c.endpoint.Plugin.Disable(ctx.Param("name"))

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation GET /plugins pluginList
// ---
// summary: get plugin list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - plugin
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
//	   $ref: '#/responses/PluginList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerPlugin) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Plugin.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Plugin, 0)
	common.Copy(&result, items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
}

// swagger:operation GET /plugin/{name}/options pluginOptions
// ---
// parameters:
// - description: Plugin name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: get plugin by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - plugin
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Plugin'
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
func (c ControllerPlugin) GetOptions(ctx *gin.Context) {

	pluginOptions, err := c.endpoint.Plugin.GetOptions(ctx.Param("name"))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.PluginOption{}
	common.Copy(&result, &pluginOptions, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}
