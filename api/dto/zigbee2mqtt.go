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

package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Zigbee2mqtt ...
type Zigbee2mqtt struct{}

// NewZigbee2mqttDto ...
func NewZigbee2mqttDto() Zigbee2mqtt {
	return Zigbee2mqtt{}
}

// AddZigbee2MqttBridgeRequest ...
func (u Zigbee2mqtt) AddZigbee2MqttBridgeRequest(obj *api.NewZigbee2MqttRequest) (bridge *m.Zigbee2mqtt) {
	bridge = &m.Zigbee2mqtt{
		Name:       obj.Name,
		Login:      obj.Login,
		Password:   obj.Password,
		PermitJoin: obj.PermitJoin,
		BaseTopic:  obj.BaseTopic,
	}
	return
}

// AddZigbee2MqttBridgeResult ...
func (u Zigbee2mqtt) AddZigbee2MqttBridgeResult(bridge *m.Zigbee2mqtt) (obj *api.Zigbee2Mqtt) {
	obj = &api.Zigbee2Mqtt{
		ScanInProcess: false,
		LastScan:      nil,
		Networkmap:    "",
		Status:        "",
		Id:            bridge.Id,
		Name:          bridge.Name,
		Login:         bridge.Login,
		PermitJoin:    bridge.PermitJoin,
		BaseTopic:     bridge.BaseTopic,
		CreatedAt:     timestamppb.New(bridge.CreatedAt),
		UpdatedAt:     timestamppb.New(bridge.UpdatedAt),
	}
	return
}

// ToZigbee2mqttInfo ...
func (u Zigbee2mqtt) ToZigbee2mqttInfo(info *zigbee2mqtt.Zigbee2mqttBridge) (obj *api.Zigbee2Mqtt) {
	obj = &api.Zigbee2Mqtt{
		ScanInProcess: info.ScanInProcess,
		Networkmap:    info.Networkmap,
		Status:        info.Status,
		Id:            info.Id,
		Name:          info.Name,
		Login:         info.Login,
		PermitJoin:    info.PermitJoin,
		BaseTopic:     info.BaseTopic,
		CreatedAt:     timestamppb.New(info.CreatedAt),
		UpdatedAt:     timestamppb.New(info.UpdatedAt),
	}
	if info.LastScan != nil {
		obj.LastScan = timestamppb.New(*info.LastScan)
	}
	return
}

// UpdateBridgeByIdRequest ...
func (u Zigbee2mqtt) UpdateBridgeByIdRequest(obj *api.UpdateBridgeRequest) (bridge *m.Zigbee2mqtt) {
	bridge = &m.Zigbee2mqtt{
		Id:         obj.Id,
		Name:       obj.Name,
		Login:      obj.Login,
		Password:   obj.Password,
		PermitJoin: obj.PermitJoin,
		BaseTopic:  obj.BaseTopic,
	}
	return
}

// UpdateBridgeByIdResult ...
func (u Zigbee2mqtt) UpdateBridgeByIdResult(bridge *m.Zigbee2mqtt) (obj *api.Zigbee2Mqtt) {
	obj = &api.Zigbee2Mqtt{
		Id:         bridge.Id,
		Name:       bridge.Name,
		Login:      bridge.Login,
		PermitJoin: bridge.PermitJoin,
		BaseTopic:  bridge.BaseTopic,
		CreatedAt:  timestamppb.New(bridge.CreatedAt),
		UpdatedAt:  timestamppb.New(bridge.UpdatedAt),
	}
	return
}

// GetBridgeListResult ...
func (u Zigbee2mqtt) GetBridgeListResult(list []*zigbee2mqtt.Zigbee2mqttBridge, total uint64, pagination common.PageParams) (obj *api.GetBridgeListResult) {
	items := make([]*api.Zigbee2MqttShort, 0, len(list))
	for _, item := range list {
		items = append(items, &api.Zigbee2MqttShort{
			Id:         item.Id,
			Name:       item.Name,
			Login:      item.Login,
			PermitJoin: item.PermitJoin,
			BaseTopic:  item.BaseTopic,
			CreatedAt:  timestamppb.New(item.CreatedAt),
			UpdatedAt:  timestamppb.New(item.UpdatedAt),
		})
	}
	obj = &api.GetBridgeListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
	return
}

// ToZigbee2MqttDevice ...
func (u Zigbee2mqtt) ToZigbee2MqttDevice(device *m.Zigbee2mqttDevice) (obj *api.Zigbee2MqttDevice) {

	return &api.Zigbee2MqttDevice{
		Id:            device.Id,
		Zigbee2MqttId: device.Zigbee2mqttId,
		Name:          device.Name,
		Type:          device.Type,
		Model:         device.Model,
		Description:   device.Description,
		Manufacturer:  device.Manufacturer,
		Functions:     device.Functions,
		ImageUrl:      device.ImageUrl,
		Status:        device.Status,
		CreatedAt:     timestamppb.New(device.CreatedAt),
		UpdatedAt:     timestamppb.New(device.UpdatedAt),
	}
}

// SearchDevice ...
func (u Zigbee2mqtt) SearchDevice(list []*m.Zigbee2mqttDevice) (obj *api.SearchDeviceResult) {

	items := make([]*api.Zigbee2MqttDevice, 0, len(list))

	for _, i := range list {
		items = append(items, u.ToZigbee2MqttDevice(i))
	}

	return &api.SearchDeviceResult{
		Items: items,
	}
}

// ToListResult ...
func (u Zigbee2mqtt) ToListResult(list []*m.Zigbee2mqttDevice, total uint64, pagination common.PageParams) *api.DeviceListResult {

	items := make([]*api.Zigbee2MqttDevice, 0, len(list))

	for _, i := range list {
		items = append(items, u.ToZigbee2MqttDevice(i))
	}

	return &api.DeviceListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}
