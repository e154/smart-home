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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/notify"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("endpoint")
)

type Endpoint struct {
	Auth             *AuthEndpoint
	Device           *DeviceEndpoint
	DeviceAction     *DeviceActionEndpoint
	DeviceState      *DeviceStateEndpoint
	Flow             *FlowEndpoint
	Image            *ImageEndpoint
	Log              *LogEndpoint
	Map              *MapEndpoint
	MapElement       *MapElementEndpoint
	MapLayer         *MapLayerEndpoint
	MapZone          *MapZoneEndpoint
	Node             *NodeEndpoint
	Role             *RoleEndpoint
	Script           *ScriptEndpoint
	Workflow         *WorkflowEndpoint
	WorkflowScenario *WorkflowScenarioEndpoint
	User             *UserEndpoint
	Gate             *GateEndpoint
	Template         *TemplateEndpoint
	Notify           *NotifyEndpoint
	MessageDelivery  *MessageDeliveryEndpoint
	Mqtt             *MqttEndpoint
	Version          *VersionEndpoint
	Zigbee2mqtt      *Zigbee2mqttEndpoint
}

func NewEndpoint(adaptors *adaptors.Adaptors,
	core *core.Core,
	scriptService *scripts.ScriptService,
	accessList *access_list.AccessListService,
	gate *gate_client.GateClient,
	notify *notify.Notify,
	mqtt *mqtt.Mqtt,
	zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt) *Endpoint {
	common := NewCommonEndpoint(adaptors, core, accessList, scriptService, gate, notify, mqtt, zigbee2mqtt)
	return &Endpoint{
		Auth:             NewAuthEndpoint(common),
		Device:           NewDeviceEndpoint(common),
		DeviceAction:     NewDeviceActionEndpoint(common),
		DeviceState:      NewDeviceStateEndpoint(common),
		Flow:             NewFlowEndpoint(common),
		Image:            NewImageEndpoint(common),
		Log:              NewLogEndpoint(common),
		Map:              NewMapEndpoint(common),
		MapElement:       NewMapElementEndpoint(common),
		MapLayer:         NewMapLayerEndpoint(common),
		Node:             NewNodeEndpoint(common),
		Role:             NewRoleEndpoint(common),
		Script:           NewScriptEndpoint(common),
		Workflow:         NewWorkflowEndpoint(common),
		WorkflowScenario: NewWorkflowScenarioEndpoint(common),
		User:             NewUserEndpoint(common),
		Gate:             NewGateEndpoint(common),
		MapZone:          NewMapZoneEndpoint(common),
		Template:         NewTemplateEndpoint(common),
		Notify:           NewNotifyEndpoint(common),
		MessageDelivery:  NewMessageDeliveryEndpoint(common),
		Mqtt:             NewMqttEndpoint(common),
		Version:          NewVersionEndpoint(common),
		Zigbee2mqtt:      NewZigbee2mqttEndpoint(common),
	}
}
