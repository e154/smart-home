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

package endpoint

import (
	"context"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/gin-gonic/gin"
)

type GateEndpoint struct {
	*CommonEndpoint
}

func NewGateEndpoint(common *CommonEndpoint) *GateEndpoint {
	return &GateEndpoint{
		CommonEndpoint: common,
	}
}

func (d *GateEndpoint) GetSettings() (settings gate_client.Settings, err error) {
	settings, err = d.gate.GetSettings()
	return
}

func (d *GateEndpoint) UpdateSettings(settings gate_client.Settings) (err error) {
	if err = d.gate.UpdateSettings(settings); err != nil {
		return
	}

	return
}

func (d *GateEndpoint) GetMobileList(ctx *gin.Context) (list *gate_client.MobileList, err error) {
	list, err = d.gate.GetMobileList(ctx)
	return
}

func (d *GateEndpoint) DeleteMobile(token string, ctx context.Context) (list *gate_client.MobileList, err error) {
	list, err = d.gate.DeleteMobile(token, ctx)
	return
}

func (d *GateEndpoint) AddMobile(ctx context.Context) (list *gate_client.MobileList, err error) {
	list, err = d.gate.AddMobile(ctx)
	return
}
