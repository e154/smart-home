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

type MobileControllersV1 struct {
	Auth     *ControllerAuth
	Workflow *ControllerWorkflow
	Map      *ControllerMap
}

func NewMobileControllersV1(adaptors *adaptors.Adaptors,
	core *core.Core,
	scriptService *scripts.ScriptService,
	accessList *access_list.AccessListService,
	command *endpoint.Endpoint) *MobileControllersV1 {
	common := NewControllerCommon(adaptors, core, accessList, command)
	return &MobileControllersV1{
		Auth:     NewControllerAuth(common),
		Workflow: NewControllerWorkflow(common),
		Map:      NewControllerMap(common),
	}
}
