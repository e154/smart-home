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
	"fmt"
	. "github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/initial/env1"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	"strconv"
)

var (
	log = common.MustGetLogger("initial")
)

var (
	currentVersion = 2
)

type InitialService struct {
	migrations    *migrations.Migrations
	adaptors      *Adaptors
	scriptService *scripts.ScriptService
	accessList    *access_list.AccessListService
}

func NewInitialService(migrations *migrations.Migrations,
	adaptors *Adaptors,
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

	env1.InstallDemoData(n.adaptors, n.accessList, n.scriptService)

	err := n.adaptors.Variable.Add(&m.Variable{
		Name:  "initial_version",
		Value: fmt.Sprintf("%d", currentVersion),
	})
	So(err, ShouldBeNil)

	log.Info("complete")
}

func (n *InitialService) Start() {

	defer func() {
		fmt.Println("")
	}()

	v, err := n.adaptors.Variable.GetByName("initial_version")
	if err != nil {

		if err.Error() == ErrRecordNotFound.Error() {
			v = &m.Variable{
				Name:  "initial_version",
				Value: fmt.Sprintf("%d", 1),
			}
			err = n.adaptors.Variable.Add(v)
			So(err, ShouldBeNil)
		}

		// create
		env1.Create(n.adaptors, n.accessList, n.scriptService)
	}

	oldVersion, err := strconv.Atoi(v.Value)
	So(err, ShouldBeNil)

	// upgrade
	env1.Upgrade(oldVersion, n.adaptors, n.accessList, n.scriptService)

	if oldVersion >= currentVersion {
		return
	}

	v.Value = fmt.Sprintf("%d", currentVersion)
	err = n.adaptors.Variable.Update(v)
	So(err, ShouldBeNil)
}
