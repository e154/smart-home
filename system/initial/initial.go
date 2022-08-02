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

package initial

import (
	"context"
	"errors"
	"fmt"

	"github.com/e154/smart-home/system/initial/demo"
	"go.uber.org/fx"

	. "github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	_ "github.com/e154/smart-home/plugins"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/gate_client"
	. "github.com/e154/smart-home/system/initial/assertions"
	localMigrations "github.com/e154/smart-home/system/initial/local_migrations"
	"github.com/e154/smart-home/system/logging_ws"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/validation"
)

var (
	log = logger.MustGetLogger("initial")
)

// Initial ...
type Initial struct {
	migrations      *migrations.Migrations
	adaptors        *Adaptors
	scriptService   scripts.ScriptService
	accessList      access_list.AccessListService
	entityManager   entity_manager.EntityManager
	pluginManager   common.PluginManager
	automation      automation.Automation
	api             *api.Api
	gateClient      *gate_client.GateClient
	validation      *validation.Validate
	localMigrations *localMigrations.Migrations
	demo            *demo.Demos
}

// NewInitial ...
func NewInitial(lc fx.Lifecycle,
	migrations *migrations.Migrations,
	adaptors *Adaptors,
	scriptService scripts.ScriptService,
	accessList access_list.AccessListService,
	entityManager entity_manager.EntityManager,
	pluginManager common.PluginManager,
	automation automation.Automation,
	api *api.Api,
	gateClient *gate_client.GateClient,
	validation *validation.Validate,
	_ *logging_ws.LoggingWs,
	localMigrations *localMigrations.Migrations,
	demo *demo.Demos) *Initial {
	initial := &Initial{
		migrations:      migrations,
		adaptors:        adaptors,
		scriptService:   scriptService,
		accessList:      accessList,
		entityManager:   entityManager,
		pluginManager:   pluginManager,
		automation:      automation,
		api:             api,
		gateClient:      gateClient,
		validation:      validation,
		localMigrations: localMigrations,
		demo:            demo,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			initial.Start()
			return nil
		},
		OnStop: func(ctx context.Context) (err error) {
			initial.Shutdown()
			return nil
		},
	})
	return initial
}

// Reset ...
func (n *Initial) Reset() {

	log.Info("full reset")

	_ = n.migrations.Purge()

	log.Info("complete")
}

// InstallDemoData ...
func (n *Initial) InstallDemoData() {

	log.Info("install demo data")

	tx := n.adaptors.Begin()

	// install demo
	_ = n.demo.InstallByName(context.TODO(), tx, "example1")

	_ = tx.Commit()

	log.Info("complete")
}

// checkForUpgrade ...
func (n *Initial) checkForUpgrade() {

	defer func() {
		fmt.Println("")
	}()

	tx := n.adaptors.Begin()

	v, err := tx.Variable.GetByName("initial_version")
	if err != nil {

		if errors.Is(err, common.ErrNotFound) {
			v = m.Variable{
				Name:  "initial_version",
				Value: fmt.Sprintf("%d", 1),
			}
			err = tx.Variable.Add(v)
			So(err, ShouldBeNil)
		}
	}

	oldVersion := v.Value
	So(err, ShouldBeNil)

	var currentVersion string
	if currentVersion, err = n.localMigrations.Up(context.TODO(), tx, oldVersion); err != nil {
		return
	}

	v.Value = currentVersion
	err = tx.Variable.Update(v)
	So(err, ShouldBeNil)

	_ = tx.Commit()
}

// Start ...
func (n *Initial) Start() {

	n.checkForUpgrade()
	n.entityManager.SetPluginManager(n.pluginManager)
	n.pluginManager.Start()
	_ = n.automation.Start()
	go func() { _ = n.api.Start() }()
	n.gateClient.Start()
}

// Shutdown ...
func (n *Initial) Shutdown() {

}
