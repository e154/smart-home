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
	"github.com/e154/smart-home/system/supervisor"

	"github.com/e154/smart-home/system/scheduler"
	"go.uber.org/fx"

	. "github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	_ "github.com/e154/smart-home/plugins"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/gate_client"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/initial/demo"
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
	supervisor      supervisor.Supervisor
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
	supervisor supervisor.Supervisor,
	automation automation.Automation,
	api *api.Api,
	gateClient *gate_client.GateClient,
	validation *validation.Validate,
	_ *logging_ws.LoggingWs,
	localMigrations *localMigrations.Migrations,
	demo *demo.Demos,
	_ *scheduler.Scheduler) *Initial {
	initial := &Initial{
		migrations:      migrations,
		adaptors:        adaptors,
		scriptService:   scriptService,
		accessList:      accessList,
		supervisor:      supervisor,
		automation:      automation,
		api:             api,
		gateClient:      gateClient,
		validation:      validation,
		localMigrations: localMigrations,
		demo:            demo,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return initial.Start(ctx)
		},
		OnStop: func(ctx context.Context) (err error) {
			return initial.Shutdown(ctx)
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

	v, err := n.adaptors.Variable.GetByName("initial_version")
	if err != nil {

		if errors.Is(err, apperr.ErrNotFound) {
			v = m.Variable{
				Name:  "initial_version",
				Value: fmt.Sprintf("%d", 1),
			}
			err = n.adaptors.Variable.Add(v)
			So(err, ShouldBeNil)
		}
	}

	oldVersion := v.Value
	So(err, ShouldBeNil)

	var currentVersion string
	if currentVersion, err = n.localMigrations.Up(context.TODO(), n.adaptors, oldVersion); err != nil {
		return
	}

	v.Value = currentVersion
	err = n.adaptors.Variable.Update(v)
	So(err, ShouldBeNil)
}

// Start ...
func (n *Initial) Start(ctx context.Context) (err error) {
	n.checkForUpgrade()
	_ = n.supervisor.Start(ctx)
	//_ = n.automation.Start()
	go func() {
		_ = n.api.Start()
	}()
	n.gateClient.Start()
	return
}

// Shutdown ...
func (n *Initial) Shutdown(ctx context.Context) (err error) {
	_ = n.api.Shutdown(ctx)
	return
}
