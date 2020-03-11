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

package initial

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/initial/env1"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
)

var (
	log = common.MustGetLogger("initial")
)

type InitialService struct {
	migrations    *migrations.Migrations
	adaptors      *adaptors.Adaptors
	scriptService *scripts.ScriptService
	accessList    *access_list.AccessListService
}

func NewInitialService(migrations *migrations.Migrations,
	adaptors *adaptors.Adaptors,
	scriptService *scripts.ScriptService,
	accessList *access_list.AccessListService) *InitialService {
	return &InitialService{
		migrations:    migrations,
		adaptors:      adaptors,
		scriptService: scriptService,
		accessList:    accessList,
	}
}

func (n *InitialService) Reset() {

	log.Info("full reset")

	n.migrations.Purge()

	log.Info("complete")
}

func (n *InitialService) InstallDemoData() {

	log.Info("install demo data")

	n.migrations.Purge()

	env1.Init(n.adaptors, n.accessList, n.scriptService)

	log.Info("complete")
}
