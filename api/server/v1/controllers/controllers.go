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
	"github.com/e154/smart-home/system/scripts"
)

// ControllersV1 ...
type ControllersV1 struct {
	Index        *ControllerIndex
	Node         *ControllerNode
	Swagger      *ControllerSwagger
	Script       *ControllerScript
	Role         *ControllerRole
	User         *ControllerUser
	Auth         *ControllerAuth
	Image        *ControllerImage
	Log          *ControllerLog
	Template     *ControllerTemplate
	TemplateItem *ControllerTemplateItem
	Notifr       *ControllerNotifr
	Version      *ControllerVersion
	Zigbee2mqtt  *ControllerZigbee2mqtt
	Entity       *ControllerEntity
}

// NewControllersV1 ...
func NewControllersV1(adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	accessList access_list.AccessListService,
	command *endpoint.Endpoint) *ControllersV1 {
	common := NewControllerCommon(adaptors, accessList, command)
	return &ControllersV1{
		Index:        NewControllerIndex(common),
		Node:         NewControllerNode(common),
		Swagger:      NewControllerSwagger(common),
		Script:       NewControllerScript(common, scriptService),
		Role:         NewControllerRole(common),
		User:         NewControllerUser(common),
		Auth:         NewControllerAuth(common),
		Image:        NewControllerImage(common),
		Log:          NewControllerLog(common),
		Template:     NewControllerTemplate(common),
		TemplateItem: NewControllerTemplateItem(common),
		Notifr:       NewControllerNotifr(common),
		Version:      NewControllerVersion(common),
		Zigbee2mqtt:  NewControllerZigbee2mqtt(common),
		Entity:       NewControllerEntity(common),
	}
}
