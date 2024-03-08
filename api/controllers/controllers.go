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
	"github.com/e154/smart-home/endpoint"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/stream/handlers"
	"github.com/e154/smart-home/system/validation"
)

// Controllers ...
type Controllers struct {
	*ControllerAuth
	*ControllerStream
	*ControllerUser
	*ControllerRole
	*ControllerScript
	*ControllerTag
	*ControllerImage
	*ControllerPlugin
	*ControllerZigbee2mqtt
	*ControllerEntity
	*ControllerAction
	*ControllerCondition
	*ControllerTrigger
	*ControllerAutomation
	*ControllerArea
	*ControllerDeveloperTools
	*ControllerInteract
	*ControllerLogs
	*ControllerDashboard
	*ControllerDashboardCardItem
	*ControllerDashboardCard
	*ControllerDashboardTab
	*ControllerVariable
	*ControllerEntityStorage
	*ControllerMetric
	*ControllerBackup
	*ControllerMessageDelivery
	*ControllerIndex
	*ControllerMqtt
	*ControllerMedia
	*ControllerWebdav
}

// NewControllers ...
func NewControllers(
	accessList access_list.AccessListService,
	endpoint *endpoint.Endpoint,
	_ *handlers.EventHandler,
	appConfig *m.AppConfig,
	validation *validation.Validate) *Controllers {
	common := NewControllerCommon(endpoint, accessList, appConfig, validation)
	return &Controllers{
		ControllerAuth:              NewControllerAuth(common),
		ControllerStream:            NewControllerStream(common),
		ControllerUser:              NewControllerUser(common),
		ControllerRole:              NewControllerRole(common),
		ControllerScript:            NewControllerScript(common),
		ControllerTag:               NewControllerTag(common),
		ControllerImage:             NewControllerImage(common),
		ControllerPlugin:            NewControllerPlugin(common),
		ControllerZigbee2mqtt:       NewControllerZigbee2mqtt(common),
		ControllerEntity:            NewControllerEntity(common),
		ControllerAction:            NewControllerAction(common),
		ControllerCondition:         NewControllerCondition(common),
		ControllerTrigger:           NewControllerTrigger(common),
		ControllerAutomation:        NewControllerAutomation(common),
		ControllerArea:              NewControllerArea(common),
		ControllerDeveloperTools:    NewControllerDeveloperTools(common),
		ControllerInteract:          NewControllerInteract(common),
		ControllerLogs:              NewControllerLogs(common),
		ControllerDashboard:         NewControllerDashboard(common),
		ControllerDashboardCardItem: NewControllerDashboardCardItem(common),
		ControllerDashboardCard:     NewControllerDashboardCard(common),
		ControllerDashboardTab:      NewControllerDashboardTab(common),
		ControllerVariable:          NewControllerVariable(common),
		ControllerEntityStorage:     NewControllerEntityStorage(common),
		ControllerMetric:            NewControllerMetric(common),
		ControllerBackup:            NewControllerBackup(common),
		ControllerMessageDelivery:   NewControllerMessageDelivery(common),
		ControllerIndex:             NewControllerIndex(common),
		ControllerMqtt:              NewControllerMqtt(common),
		ControllerMedia:             NewControllerMedia(common),
		ControllerWebdav:            NewControllerWebdav(common),
	}
}
