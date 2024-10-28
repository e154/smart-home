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

package dto

import (
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/pkg/models"
)

// Zigbee2mqtt ...
type Zigbee2mqtt struct{}

// NewZigbee2mqttDto ...
func NewZigbee2mqttDto() Zigbee2mqtt {
	return Zigbee2mqtt{}
}

// AddZigbee2MqttBridgeRequest ...
func (u Zigbee2mqtt) AddZigbee2MqttBridgeRequest(obj *stub.ApiNewZigbee2mqttRequest) (bridge *models.Zigbee2mqtt) {
	bridge = &models.Zigbee2mqtt{
		Name:       obj.Name,
		Login:      obj.Login,
		Password:   obj.Password,
		PermitJoin: obj.PermitJoin,
		BaseTopic:  obj.BaseTopic,
	}
	return
}

// AddZigbee2MqttBridgeResult ...
func (u Zigbee2mqtt) AddZigbee2MqttBridgeResult(bridge *models.Zigbee2mqtt) (obj stub.ApiZigbee2mqtt) {
	obj = stub.ApiZigbee2mqtt{
		ScanInProcess: false,
		LastScan:      nil,
		Networkmap:    "",
		Status:        "",
		Id:            bridge.Id,
		Name:          bridge.Name,
		Login:         bridge.Login,
		PermitJoin:    bridge.PermitJoin,
		BaseTopic:     bridge.BaseTopic,
		CreatedAt:     bridge.CreatedAt,
		UpdatedAt:     bridge.UpdatedAt,
	}
	return
}

// ToZigbee2mqttInfo ...
func (u Zigbee2mqtt) ToZigbee2mqttInfo(bridge *models.Zigbee2mqtt) (obj *stub.ApiZigbee2mqtt) {
	obj = &stub.ApiZigbee2mqtt{
		Id:         bridge.Id,
		Name:       bridge.Name,
		Login:      bridge.Login,
		PermitJoin: bridge.PermitJoin,
		BaseTopic:  bridge.BaseTopic,
		CreatedAt:  bridge.CreatedAt,
		UpdatedAt:  bridge.UpdatedAt,
	}
	if bridge.Info != nil {
		obj.ScanInProcess = bridge.Info.ScanInProcess
		obj.Networkmap = bridge.Info.Networkmap
		obj.Status = bridge.Info.Status
		obj.LastScan = bridge.Info.LastScan
	}
	return
}

// UpdateBridgeByIdRequest ...
func (u Zigbee2mqtt) UpdateBridgeByIdRequest(obj *stub.Zigbee2mqttServiceUpdateBridgeByIdJSONBody, id int64) (bridge *models.Zigbee2mqtt) {
	bridge = &models.Zigbee2mqtt{
		Id:         id,
		Name:       obj.Name,
		Login:      obj.Login,
		Password:   obj.Password,
		PermitJoin: obj.PermitJoin,
		BaseTopic:  obj.BaseTopic,
	}
	return
}

// UpdateBridgeByIdResult ...
func (u Zigbee2mqtt) UpdateBridgeByIdResult(bridge *models.Zigbee2mqtt) (obj *stub.ApiZigbee2mqtt) {
	obj = &stub.ApiZigbee2mqtt{
		Id:         bridge.Id,
		Name:       bridge.Name,
		Login:      bridge.Login,
		PermitJoin: bridge.PermitJoin,
		BaseTopic:  bridge.BaseTopic,
		CreatedAt:  bridge.CreatedAt,
		UpdatedAt:  bridge.UpdatedAt,
	}
	return
}

// GetBridgeListResult ...
func (u Zigbee2mqtt) GetBridgeListResult(list []*models.Zigbee2mqtt) []*stub.ApiZigbee2mqttShort {
	items := make([]*stub.ApiZigbee2mqttShort, 0, len(list))
	for _, item := range list {
		items = append(items, &stub.ApiZigbee2mqttShort{
			Id:         item.Id,
			Name:       item.Name,
			Login:      item.Login,
			PermitJoin: item.PermitJoin,
			BaseTopic:  item.BaseTopic,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		})
	}

	return items
}

// ToZigbee2MqttDevice ...
func (u Zigbee2mqtt) ToZigbee2MqttDevice(device *models.Zigbee2mqttDevice) (obj *stub.ApiZigbee2mqttDevice) {

	return &stub.ApiZigbee2mqttDevice{
		Id:            device.Id,
		Zigbee2mqttId: device.Zigbee2mqttId,
		Name:          device.Name,
		Type:          device.Type,
		Model:         device.Model,
		Description:   device.Description,
		Manufacturer:  device.Manufacturer,
		Functions:     device.Functions,
		ImageUrl:      device.ImageUrl,
		Status:        device.Status,
		CreatedAt:     device.CreatedAt,
		UpdatedAt:     device.UpdatedAt,
	}
}

// SearchDevice ...
func (u Zigbee2mqtt) SearchDevice(list []*models.Zigbee2mqttDevice) (obj *stub.ApiSearchDeviceResult) {

	items := make([]stub.ApiZigbee2mqttDevice, 0, len(list))

	for _, i := range list {
		items = append(items, *u.ToZigbee2MqttDevice(i))
	}

	return &stub.ApiSearchDeviceResult{
		Items: items,
	}
}

// ToListResult ...
func (u Zigbee2mqtt) ToListResult(list []*models.Zigbee2mqttDevice) []*stub.ApiZigbee2mqttDevice {

	items := make([]*stub.ApiZigbee2mqttDevice, 0, len(list))

	for _, i := range list {
		items = append(items, u.ToZigbee2MqttDevice(i))
	}

	return items
}
