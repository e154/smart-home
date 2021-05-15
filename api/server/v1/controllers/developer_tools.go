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
	"github.com/gin-gonic/gin"
)

// ControllerDeveloperTools ...
type ControllerDeveloperTools struct {
	*ControllerCommon
}

// NewControllerDeveloperTools ...
func NewControllerDeveloperTools(common *ControllerCommon) *ControllerDeveloperTools {
	return &ControllerDeveloperTools{ControllerCommon: common}
}

// swagger:operation GET /developer_tools/states stateList
// ---
// summary: get state list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - developer_tools
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
//	   $ref: '#/responses/DeveloperToolsStateList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeveloperTools) GetStateList(ctx *gin.Context) {

	_, _, _, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.DeveloperTools.StateList()
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]models.DeveloperToolsState, 0, len(items))
	for _, item := range items {
		state := models.DeveloperToolsState{
			Id:         item.Id.String(),
			Attributes: item.Attributes.Serialize(),
		}
		if item.State != nil {
			state.State = common.String(item.State.Name)
		}
		result = append(result, state)
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation PATCH /developer_tools/state stateUpdate
// ---
// parameters:
// - description: Update state params
//   in: body
//   name: state
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateDeveloperToolsState'
//     type: object
// summary: update state
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - developer_tools
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/DeveloperToolsState'
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
func (c ControllerDeveloperTools) UpdateState(ctx *gin.Context) {

	state := models.UpdateDeveloperToolsState{}
	if err := ctx.ShouldBindJSON(&state); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	errs, err := c.endpoint.DeveloperTools.UpdateState(state.Id, state.State, state.Attributes)
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
	resp.SetData(state).Send(ctx)
}
