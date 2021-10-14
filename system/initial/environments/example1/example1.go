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

package example1

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/initial/environments"
	"github.com/e154/smart-home/system/scripts"
)

type Example1 struct{}

// Demo ...
func (e Example1) InstallDemoData(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	scriptService scripts.ScriptService) {

	scripts := NewScriptManager(adaptors, scriptService).Create()
	entities := NewEntityManager(adaptors).Create(scripts)
	NewTriggerManager(adaptors).Create(scripts, entities)
}

// Create ...
func (e Example1) Create(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	scriptService scripts.ScriptService) {


}

// Upgrade ...
func (e Example1) Upgrade(oldVersion int,
	adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	scriptService scripts.ScriptService) {

}

func init() {
	environments.Register("example1", &Example1{})
}

