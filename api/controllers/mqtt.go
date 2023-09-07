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
	"context"

	"github.com/e154/smart-home/api/stub/api"
)

// ControllerMqtt ...
type ControllerMqtt struct {
	*ControllerCommon
}

// NewControllerMqtt ...
func NewControllerMqtt(common *ControllerCommon) ControllerMqtt {
	return ControllerMqtt{
		ControllerCommon: common,
	}
}

// GetClientById ...
func (c ControllerMqtt) GetClientById(ctx context.Context, req *api.GetClientRequest) (*api.Client, error) {

	client, err := c.endpoint.Mqtt.GetClientById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Mqtt.GetClientById(client), nil
}

// GetClientList ...
func (c ControllerMqtt) GetClientList(ctx context.Context, req *api.PaginationRequest) (*api.GetClientListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.Mqtt.GetClientList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Mqtt.ToListResult(items, uint64(total), pagination), nil
}

// GetSubscriptionList ...
func (c ControllerMqtt) GetSubscriptionList(ctx context.Context, req *api.SubscriptionPaginationRequest) (*api.GetSubscriptionListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.Mqtt.GetSubscriptionList(ctx, req.ClientId, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Mqtt.GetSubscriptionList(items, uint64(total), pagination), nil
}
