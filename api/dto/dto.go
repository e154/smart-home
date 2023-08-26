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

// Dto ...
type Dto struct {
	Role              Role
	User              User
	Image             Image
	Script            Script
	Plugin            Plugin
	Entity            Entity
	Zigbee2mqtt       Zigbee2mqtt
	Area              Area
	Action            Action
	Condition         Condition
	Trigger           Trigger
	Automation        Automation
	Log               Log
	Dashboard         Dashboard
	DashboardCard     DashboardCard
	DashboardCardItem DashboardCardItem
	DashboardTab      DashboardTab
	Variable          Variable
	EntityStorage     EntityStorage
	MessageDelivery   MessageDelivery
}

// NewDto ...
func NewDto() Dto {
	return Dto{
		Role:              NewRoleDto(),
		User:              NewUserDto(),
		Image:             NewImageDto(),
		Script:            NewScriptDto(),
		Plugin:            NewPluginDto(),
		Entity:            NewEntityDto(),
		Zigbee2mqtt:       NewZigbee2mqttDto(),
		Area:              NewAreaDto(),
		Action:            NewActionDto(),
		Condition:         NewConditionDto(),
		Trigger:           NewTriggerDto(),
		Automation:        NewAutomationDto(),
		Log:               NewLogDto(),
		Dashboard:         NewDashboardDto(),
		DashboardCard:     NewDashboardCardDto(),
		DashboardCardItem: NewDashboardCardItemDto(),
		DashboardTab:      NewDashboardTabDto(),
		Variable:          NewVariableDto(),
		EntityStorage:     NewEntityStorageDto(),
		MessageDelivery:   NewMessageDeliveryDto(),
	}
}
