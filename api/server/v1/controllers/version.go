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

type ControllerVersion struct {
	*ControllerCommon
}

func NewControllerVersion(common *ControllerCommon) *ControllerVersion {
	return &ControllerVersion{ControllerCommon: common}
}

// swagger:operation GET /version getServerVersion
// ---
// summary: get server version
// description:
// tags:
// - version
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Version'
func (c ControllerVersion) Version(ctx *gin.Context) {

	serverVersion := c.endpoint.Version.ServerVersion()

	result := &models.Version{}
	common.Copy(&result, &serverVersion)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}
