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
	"github.com/gin-gonic/gin"
	"strconv"
)

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
)

// ControllerMapDeviceHistory ...
type ControllerMapDeviceHistory struct {
	*ControllerCommon
}

// NewControllerMapDeviceHistory ...
func NewControllerMapDeviceHistory(common *ControllerCommon) *ControllerMapDeviceHistory {
	return &ControllerMapDeviceHistory{ControllerCommon: common}
}

// swagger:operation GET /history/map MapDeviceHistoryList
// ---
// summary: get map device history
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map_device_history
// parameters:
// - name: map_id
//   description: map id
//   in: query
//   type: integer
//   required: true
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
//     $ref: '#/responses/MapDeviceHistoryList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMapDeviceHistory) GetList(ctx *gin.Context) {

	_, orderBy, sort, limit, offset := c.list(ctx)

	var mapId int64

	if ctx.Request.URL.Query().Get("map_id") != "" {
		_mapId, _ := strconv.Atoi(ctx.Request.URL.Query().Get("map_id"))
		mapId = int64(_mapId)
	}

	items, total, err := c.endpoint.MapDeviceHistory.ListByMapId(mapId, limit, offset, orderBy, sort)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MapDeviceHistory, 0)
	common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}
