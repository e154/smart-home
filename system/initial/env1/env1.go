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

package env1

import (
	. "github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/access_list"
	. "github.com/e154/smart-home/system/scripts"
)

// env1
//
// node1 + node2
// script1 + script2 + script3
// device1
// 		+ child device2
// device3
//
func InstallDemoData(adaptors *Adaptors,
	accessList *access_list.AccessListService,
	scriptService *ScriptService) {

	// images
	// ------------------------------------------------
	imageList := NewImageManager(adaptors).Create()

	// roles
	// ------------------------------------------------
	NewRoleManager(adaptors, accessList).Create()

	// nodes
	// ------------------------------------------------
	node1, _ := NewNodeManager(adaptors).Create()

	// scripts
	// ------------------------------------------------
	scripts := NewScriptManager(adaptors, scriptService).Create()

	// devices
	// ------------------------------------------------
	devList, deviceActions, deviceStates := devices(node1, adaptors, scripts)

	// workflow
	// ------------------------------------------------
	addWorkflow(adaptors, deviceActions, scripts)

	// maps
	// ------------------------------------------------
	addMaps(adaptors, scripts, devList, imageList, deviceActions, deviceStates)

	// templates
	// ------------------------------------------------
	NewTemplateManager(adaptors).Create()
}

func Create(adaptors *Adaptors,
	accessList *access_list.AccessListService,
	scriptService *ScriptService) {

	NewImageManager(adaptors).Create()
	NewRoleManager(adaptors, accessList).Create()
	NewTemplateManager(adaptors).Create()
	NewNodeManager(adaptors).Create()
}

func Upgrade(oldVersion int,
	adaptors *Adaptors,
	accessList *access_list.AccessListService,
	scriptService *ScriptService) {

	NewImageManager(adaptors).Upgrade(oldVersion)
	NewTemplateManager(adaptors).Upgrade(oldVersion)
	NewScriptManager(adaptors, scriptService).Upgrade(oldVersion)
	NewRoleManager(adaptors, accessList).Upgrade(oldVersion)
}
