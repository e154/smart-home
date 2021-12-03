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

package environments

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/validation"
)

// IEnvironment ...
type IEnvironment interface {
	InstallDemoData(*adaptors.Adaptors, access_list.AccessListService, scripts.ScriptService)
	Create(*adaptors.Adaptors, access_list.AccessListService, scripts.ScriptService, *validation.Validate)
	Upgrade(int, *adaptors.Adaptors, access_list.AccessListService, scripts.ScriptService, *validation.Validate)
}

var environments = map[string]IEnvironment{}

// Register ...
func Register(name string, env IEnvironment) {
	if _, ok := environments[name]; ok {
		panic("duplicated environment: " + name)
	}
	environments[name] = env
}

// InstallDemoData ...
func InstallDemoData(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	scriptService scripts.ScriptService) {
	for _, env := range environments {
		env.InstallDemoData(adaptors, accessList, scriptService)
	}
}

// Create ...
func Create(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	scriptService scripts.ScriptService,
	validation *validation.Validate) {
	for _, env := range environments {
		env.Create(adaptors, accessList, scriptService, validation)
	}
}

// Upgrade ...
func Upgrade(oldVersion int,
	adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	scriptService scripts.ScriptService,
	validation *validation.Validate) {
	for _, env := range environments {
		env.Upgrade(oldVersion, adaptors, accessList, scriptService, validation)
	}
}
