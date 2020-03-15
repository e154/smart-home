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
	"strconv"
)

type ControllerMapDevice struct {
	*ControllerCommon
}

func NewControllerMapDevice(common *ControllerCommon) *ControllerMapDevice {
	return &ControllerMapDevice{ControllerCommon: common}
}

func (c *ControllerMapDevice) GetHistory(ctx *gin.Context) {

	var limit, offset int
	var id int64

	if ctx.Request.URL.Query().Get("limit") != "" {
		limit, _ = strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
	}

	if ctx.Request.URL.Query().Get("offset") != "" {
		offset, _ = strconv.Atoi(ctx.Request.URL.Query().Get("offset"))
	}

	if ctx.Request.URL.Query().Get("id") != "" {
		_id, _ := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
		id = int64(_id)
	}

	items, total, err := c.endpoint.MapDeviceHistory.GetAllByDeviceId(id, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MapDeviceHistory, 0)
	_ = common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
}
