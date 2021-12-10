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
	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerZigbee2mqtt ...
type ControllerZigbee2mqtt struct {
	*ControllerCommon
}

// NewControllerZigbee2mqtt ...
func NewControllerZigbee2mqtt(common *ControllerCommon) ControllerZigbee2mqtt {
	return ControllerZigbee2mqtt{
		ControllerCommon: common,
	}
}

// AddZigbee2MqttBridge ...
func (c ControllerZigbee2mqtt) AddZigbee2MqttBridge(ctx context.Context, req *api.NewtZigbee2MqttRequest) (*api.Zigbee2Mqtt, error) {

	bridge := c.dto.Zigbee2mqtt.AddZigbee2MqttBridgeRequest(req)

	result, errs, err := c.endpoint.Zigbee2mqtt.AddBridge(ctx, bridge)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Zigbee2mqtt.AddZigbee2MqttBridgeResult(result), nil
}

// GetZigbee2MqttBridge ...
func (c ControllerZigbee2mqtt) GetZigbee2MqttBridge(ctx context.Context, req *api.GetBridgeRequest) (*api.Zigbee2MqttInfo, error) {

	info, err := c.endpoint.Zigbee2mqtt.GetBridgeById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Zigbee2mqtt.ToZigbee2mqttInfo(info), nil
}

// UpdateBridgeById ...
func (c ControllerZigbee2mqtt) UpdateBridgeById(ctx context.Context, req *api.UpdateBridgeRequest) (*api.Zigbee2Mqtt, error) {

	bridge := c.dto.Zigbee2mqtt.UpdateBridgeByIdRequest(req)

	bridge, errs, err := c.endpoint.Zigbee2mqtt.UpdateBridge(ctx, bridge)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Zigbee2mqtt.UpdateBridgeByIdResult(bridge), nil
}

// GetBridgeList ...
func (c ControllerZigbee2mqtt) GetBridgeList(ctx context.Context, req *api.GetBridgeListRequest) (*api.GetBridgeListResult, error) {

	pagination := c.Pagination(req.Limit, req.Offset, req.Order, req.SortBy)
	items, total, err := c.endpoint.Zigbee2mqtt.GetBridgeList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Zigbee2mqtt.GetBridgeListResult(items, uint64(total), req.Limit, req.Offset), nil
}

// DeleteBridgeById ...
func (c ControllerZigbee2mqtt) DeleteBridgeById(ctx context.Context, req *api.DeleteBridgeRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

// ResetBridgeById ...
func (c ControllerZigbee2mqtt) ResetBridgeById(ctx context.Context, req *api.ResetBridgeRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

// DeviceBan ...
func (c ControllerZigbee2mqtt) DeviceBan(ctx context.Context, req *api.DeviceBanRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

// DeviceWhitelist ...
func (c ControllerZigbee2mqtt) DeviceWhitelist(ctx context.Context, req *api.DeviceWhitelistRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

// DeviceRename ...
func (c ControllerZigbee2mqtt) DeviceRename(ctx context.Context, req *api.DeviceRenameRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

// SearchDevice ...
func (c ControllerZigbee2mqtt) SearchDevice(ctx context.Context, req *api.SearchDeviceRequest) (*api.SearchDeviceResult, error) {
	panic("implement me")
}

// Networkmap ...
func (c ControllerZigbee2mqtt) Networkmap(ctx context.Context, req *api.NetworkmapRequest) (*api.NetworkmapResponse, error) {
	panic("implement me")
}

// UpdateNetworkmap ...
func (c ControllerZigbee2mqtt) UpdateNetworkmap(ctx context.Context, req *api.NetworkmapRequest) (*emptypb.Empty, error) {
	panic("implement me")
}
