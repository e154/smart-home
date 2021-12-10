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
func (u Zigbee2mqtt) AddZigbee2MqttBridgeRequest(obj *api.NewtZigbee2MqttRequest) (bridge *m.Zigbee2mqtt) {
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

// ToZigbee2mqttInfo ...
func (u Zigbee2mqtt) ToZigbee2mqttInfo(info *zigbee2mqtt.Zigbee2mqttInfo) (obj *api.Zigbee2MqttInfo) {
	obj = &api.Zigbee2MqttInfo{
		ScanInProcess: info.ScanInProcess,
		Networkmap:    info.Networkmap,
		Status:        info.Status,
		Model:         u.AddZigbee2MqttBridgeResult(&info.Model),
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

func (u Zigbee2mqtt) GetBridgeListResult(list []*zigbee2mqtt.Zigbee2mqttInfo, total, limit, offset uint64) (obj *api.GetBridgeListResult) {
	items := make([]*api.Zigbee2MqttInfo, 0, len(list))
	for _, i := range list {
		items = append(items, u.ToZigbee2mqttInfo(i))
	}
	obj = &api.GetBridgeListResult{
		Items: items,
		Meta: &api.GetBridgeListResult_Meta{
			Limit:        limit,
			ObjectsCount: total,
			Offset:       offset,
		},
	}
	return
}
