// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/labstack/echo/v4"
)

// ControllerMqtt ...
type ControllerMqtt struct {
	*ControllerCommon
}

// NewControllerMqtt ...
func NewControllerMqtt(common *ControllerCommon) *ControllerMqtt {
	return &ControllerMqtt{
		ControllerCommon: common,
	}
}

// GetClientById ...
func (c ControllerMqtt) MqttServiceGetClientById(ctx echo.Context, id string) error {

	client, err := c.endpoint.Mqtt.GetClientById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Mqtt.GetClientById(client)))
}

// GetClientList ...
func (c ControllerMqtt) MqttServiceGetClientList(ctx echo.Context, params stub.MqttServiceGetClientListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Mqtt.GetClientList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Mqtt.ToListResult(items), int64(total), pagination))
}

// GetSubscriptionList ...
func (c ControllerMqtt) MqttServiceGetSubscriptionList(ctx echo.Context, params stub.MqttServiceGetSubscriptionListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Mqtt.GetSubscriptionList(ctx.Request().Context(), params.ClientId, pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Mqtt.GetSubscriptionList(items), int64(total), pagination))
}
