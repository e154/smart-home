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

// ControllerZigbee2mqtt ...
type ControllerZigbee2mqtt struct {
	*ControllerCommon
}

// NewControllerZigbee2mqtt ...
func NewControllerZigbee2mqtt(common *ControllerCommon) *ControllerZigbee2mqtt {
	return &ControllerZigbee2mqtt{
		ControllerCommon: common,
	}
}

// AddZigbee2MqttBridge ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceAddZigbee2mqttBridge(ctx echo.Context, _ stub.Zigbee2mqttServiceAddZigbee2mqttBridgeParams) error {

	obj := &stub.ApiNewZigbee2mqttRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	bridge := c.dto.Zigbee2mqtt.AddZigbee2MqttBridgeRequest(obj)

	result, err := c.endpoint.Zigbee2mqtt.AddBridge(ctx.Request().Context(), bridge)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Zigbee2mqtt.AddZigbee2MqttBridgeResult(result)))
}

// GetZigbee2MqttBridge ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceGetZigbee2mqttBridge(ctx echo.Context, id int64) error {

	bridge, err := c.endpoint.Zigbee2mqtt.GetBridgeById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Zigbee2mqtt.ToZigbee2mqttInfo(bridge)))
}

// UpdateBridgeById ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceUpdateBridgeById(ctx echo.Context, id int64, _ stub.Zigbee2mqttServiceUpdateBridgeByIdParams) error {

	obj := &stub.Zigbee2mqttServiceUpdateBridgeByIdJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	bridge := c.dto.Zigbee2mqtt.UpdateBridgeByIdRequest(obj, id)

	bridge, err := c.endpoint.Zigbee2mqtt.UpdateBridge(ctx.Request().Context(), bridge)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Zigbee2mqtt.UpdateBridgeByIdResult(bridge)))
}

// GetBridgeList ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceGetBridgeList(ctx echo.Context, params stub.Zigbee2mqttServiceGetBridgeListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Zigbee2mqtt.GetBridgeList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Zigbee2mqtt.GetBridgeListResult(items), total, pagination))
}

// DeleteBridgeById ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceDeleteBridgeById(ctx echo.Context, id int64) error {

	err := c.endpoint.Zigbee2mqtt.Delete(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// ResetBridgeById ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceResetBridgeById(ctx echo.Context, id int64) error {

	err := c.endpoint.Zigbee2mqtt.ResetBridge(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// DeviceBan ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceDeviceBan(ctx echo.Context, _ stub.Zigbee2mqttServiceDeviceBanParams) error {

	obj := &stub.ApiDeviceBanRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	err := c.endpoint.Zigbee2mqtt.DeviceBan(ctx.Request().Context(), obj.Id, obj.FriendlyName)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// DeviceWhitelist ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceDeviceWhitelist(ctx echo.Context, params stub.Zigbee2mqttServiceDeviceWhitelistParams) error {

	obj := &stub.ApiDeviceWhitelistRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	err := c.endpoint.Zigbee2mqtt.DeviceWhitelist(ctx.Request().Context(), obj.Id, obj.FriendlyName)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// DeviceRename ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceDeviceRename(ctx echo.Context, params stub.Zigbee2mqttServiceDeviceRenameParams) error {

	obj := &stub.ApiDeviceRenameRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	err := c.endpoint.Zigbee2mqtt.DeviceRename(ctx.Request().Context(), obj.FriendlyName, obj.NewName)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// SearchDevice ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceSearchDevice(ctx echo.Context, params stub.Zigbee2mqttServiceSearchDeviceParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Zigbee2mqtt.SearchDevice(ctx.Request().Context(), search)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Zigbee2mqtt.SearchDevice(items))
}

// Networkmap ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceNetworkmap(ctx echo.Context, id int64) error {

	networkMap, err := c.endpoint.Zigbee2mqtt.Networkmap(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, &stub.ApiNetworkmapResponse{Networkmap: networkMap}))
}

// UpdateNetworkmap ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceUpdateNetworkmap(ctx echo.Context, id int64) error {

	err := c.endpoint.Zigbee2mqtt.UpdateNetworkmap(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// DeviceList ...
func (c ControllerZigbee2mqtt) Zigbee2mqttServiceDeviceList(ctx echo.Context, id int64, params stub.Zigbee2mqttServiceDeviceListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Zigbee2mqtt.DeviceList(ctx.Request().Context(), id, pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Zigbee2mqtt.ToListResult(items), total, pagination))
}
