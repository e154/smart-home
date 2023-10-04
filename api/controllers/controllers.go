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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/stream/handlers"
)

// Controllers ...
type Controllers struct {
	Auth              ControllerAuth
	Stream            ControllerStream
	User              ControllerUser
	Role              ControllerRole
	Script            ControllerScript
	Image             ControllerImage
	Plugin            ControllerPlugin
	Zigbee2mqtt       ControllerZigbee2mqtt
	Entity            ControllerEntity
	Action            ControllerAction
	Condition         ControllerCondition
	Trigger           ControllerTrigger
	Automation        ControllerAutomation
	Area              ControllerArea
	Dev               ControllerDeveloperTools
	Interact          ControllerInteract
	Logs              ControllerLogs
	Dashboard         ControllerDashboard
	DashboardCardItem ControllerDashboardCardItem
	DashboardCard     ControllerDashboardCard
	DashboardTab      ControllerDashboardTab
	Variable          ControllerVariable
	EntityStorage     ControllerEntityStorage
	Metric            ControllerMetric
	Backup            ControllerBackup
	MessageDelivery   ControllerMessageDelivery
	Index             ControllerIndex
	Mqtt              ControllerMqtt
	Media             ControllerMedia
}

// NewControllers ...
func NewControllers(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	command *endpoint.Endpoint,
	_ *handlers.EventHandler,
	appConfig *m.AppConfig) *Controllers {
	common := NewControllerCommon(adaptors, accessList, command, appConfig)
	return &Controllers{
		Auth:              NewControllerAuth(common),
		Stream:            NewControllerStream(common),
		User:              NewControllerUser(common),
		Role:              NewControllerRole(common),
		Script:            NewControllerScript(common),
		Image:             NewControllerImage(common),
		Plugin:            NewControllerPlugin(common),
		Zigbee2mqtt:       NewControllerZigbee2mqtt(common),
		Entity:            NewControllerEntity(common),
		Action:            NewControllerAction(common),
		Condition:         NewControllerCondition(common),
		Trigger:           NewControllerTrigger(common),
		Automation:        NewControllerAutomation(common),
		Area:              NewControllerArea(common),
		Dev:               NewControllerDeveloperTools(common),
		Interact:          NewControllerInteract(common),
		Logs:              NewControllerLogs(common),
		Dashboard:         NewControllerDashboard(common),
		DashboardCardItem: NewControllerDashboardCardItem(common),
		DashboardCard:     NewControllerDashboardCard(common),
		DashboardTab:      NewControllerDashboardTab(common),
		Variable:          NewControllerVariable(common),
		EntityStorage:     NewControllerEntityStorage(common),
		Metric:            NewControllerMetric(common),
		Backup:            NewControllerBackup(common),
		MessageDelivery:   NewControllerMessageDelivery(common),
		Index:             NewControllerIndex(common),
		Mqtt:              NewControllerMqtt(common),
		Media:             NewControllerMedia(common),
	}
}
