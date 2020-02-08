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
	"github.com/e154/smart-home/api/mobile/v1/models"
	"github.com/e154/smart-home/common"
	"github.com/gin-gonic/gin"
)

type ControllerMap struct {
	*ControllerCommon
}

func NewControllerMap(common *ControllerCommon) *ControllerMap {
	return &ControllerMap{ControllerCommon: common}
}

// swagger:operation GET /map/active_elements mapGetActiveElements
// ---
// summary: get active map elements
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
//     description: OK
//     schema:
//       $ref: '#/responses/MapActiveElementList'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMap) GetActiveElements(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Map.GetActiveElements(sortBy, order, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MapElement, 0)
	_ = common.Copy(&result, &items, common.JsonEngine)

	//fmt.Println("<<<<<<<<<")
	//debug.Println(result)
	//fmt.Println(">>>>>>>>>>")

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
}
