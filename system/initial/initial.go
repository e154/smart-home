// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"encoding/hex"
	"errors"
	"fmt"
	"go.uber.org/fx"
	"gorm.io/gorm"

	. "github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/encryptor"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	_ "github.com/e154/smart-home/plugins"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/gate/client"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/initial/demo"
	localMigrations "github.com/e154/smart-home/system/initial/local_migrations"
	"github.com/e154/smart-home/system/logging_ws"
	"github.com/e154/smart-home/system/media"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
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
	gateClient      *client.GateClient
	validation      *validation.Validate
	localMigrations *localMigrations.Migrations
	demo            *demo.Demos
	eventBus        bus.Bus
	db              *gorm.DB
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
	gateClient *client.GateClient,
	validation *validation.Validate,
	_ *logging_ws.LoggingWs,
	localMigrations *localMigrations.Migrations,
	demo *demo.Demos,
	_ *scheduler.Scheduler,
	_ *media.Media,
	db *gorm.DB,
	eventBus bus.Bus) *Initial {
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
		db:              db,
		eventBus:        eventBus,
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

	const name = "initialVersion"
	v, err := n.adaptors.Variable.GetByName(context.Background(), name)
	if err != nil {

		if errors.Is(err, apperr.ErrNotFound) {
			v = m.Variable{
				Name:   name,
				Value:  fmt.Sprintf("%d", 1),
				System: true,
			}
			err = n.adaptors.Variable.Add(context.Background(), v)
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
	err = n.adaptors.Variable.Update(context.Background(), v)
	So(err, ShouldBeNil)
}

// Start ...
func (n *Initial) Start(ctx context.Context) (err error) {
	n.checkForUpgrade()

	variable, _ := n.adaptors.Variable.GetByName(ctx, "encryptor")
	val, _ := hex.DecodeString(variable.Value)
	encryptor.SetKey(val)

	_ = n.eventBus.Subscribe("system/models/variables/+", n.eventHandler)

	_ = n.gateClient.Start()
	_ = n.supervisor.Start(ctx)
	_ = n.automation.Start()
	go func() {
		_ = n.api.Start()
	}()
	return
}

// Shutdown ...
func (n *Initial) Shutdown(ctx context.Context) (err error) {
	_ = n.eventBus.Unsubscribe("system/models/variables/+", n.eventHandler)
	_ = n.api.Shutdown(ctx)
	return
}

func (n *Initial) eventHandler(_ string, message interface{}) {
	switch v := message.(type) {
	case events.EventUpdatedVariableModel:
		if v.Name == "timezone" && v.Value != "" {
			log.Infof("update database timezone to: \"%s\"", v.Value)
			n.db.Exec(`SET TIME ZONE '?';`, v.Value)
		}
	}
}
