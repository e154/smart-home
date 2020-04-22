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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
)

// ControllersV1 ...
type ControllersV1 struct {
	Index            *ControllerIndex
	Node             *ControllerNode
	Swagger          *ControllerSwagger
	Script           *ControllerScript
	Workflow         *ControllerWorkflow
	Device           *ControllerDevice
	Role             *ControllerRole
	User             *ControllerUser
	Auth             *ControllerAuth
	DeviceAction     *ControllerDeviceAction
	DeviceState      *ControllerDeviceState
	Map              *ControllerMap
	MapLayer         *ControllerMapLayer
	MapElement       *ControllerMapElement
	Image            *ControllerImage
	WorkflowScenario *ControllerWorkflowScenario
	Flow             *ControllerFlow
	Log              *ControllerLog
	Gate             *ControllerGate
	MapZone          *ControllerMapZone
	Template         *ControllerTemplate
	TemplateItem     *ControllerTemplateItem
	Notifr           *ControllerNotifr
	Mqtt             *ControllerMqtt
	Version          *ControllerVersion
	Zigbee2mqtt      *ControllerZigbee2mqtt
	MapDeviceHistory *ControllerMapDeviceHistory
	Alexa            *ControllerAlexa
	Metric           *ControllerMetric
}

// NewControllersV1 ...
func NewControllersV1(adaptors *adaptors.Adaptors,
	core *core.Core,
	scriptService *scripts.ScriptService,
	accessList *access_list.AccessListService,
	command *endpoint.Endpoint) *ControllersV1 {
	common := NewControllerCommon(adaptors, core, accessList, command)
	return &ControllersV1{
		Index:            NewControllerIndex(common),
		Node:             NewControllerNode(common),
		Swagger:          NewControllerSwagger(common),
		Script:           NewControllerScript(common, scriptService),
		Workflow:         NewControllerWorkflow(common),
		Device:           NewControllerDevice(common),
		Role:             NewControllerRole(common),
		User:             NewControllerUser(common),
		Auth:             NewControllerAuth(common),
		DeviceAction:     NewControllerDeviceAction(common),
		DeviceState:      NewControllerDeviceState(common),
		Map:              NewControllerMap(common),
		MapLayer:         NewControllerMapLayer(common),
		MapElement:       NewControllerMapElement(common),
		Image:            NewControllerImage(common),
		WorkflowScenario: NewControllerWorkflowScenario(common),
		Flow:             NewControllerFlow(common),
		Log:              NewControllerLog(common),
		Gate:             NewControllerGate(common),
		MapZone:          NewControllerMapZone(common),
		Template:         NewControllerTemplate(common),
		TemplateItem:     NewControllerTemplateItem(common),
		Notifr:           NewControllerNotifr(common),
		Mqtt:             NewControllerMqtt(common),
		Version:          NewControllerVersion(common),
		Zigbee2mqtt:      NewControllerZigbee2mqtt(common),
		MapDeviceHistory: NewControllerMapDeviceHistory(common),
		Alexa:            NewControllerAlexa(common),
		Metric:           NewControllerMetric(common),
	}
}
