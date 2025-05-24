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

	"github.com/e154/smart-home/internal/api"
	"github.com/e154/smart-home/internal/endpoint"
	"github.com/e154/smart-home/internal/system/automation"
	"github.com/e154/smart-home/internal/system/gate/client"
	. "github.com/e154/smart-home/internal/system/initial/assertions"
	"github.com/e154/smart-home/internal/system/initial/demo"
	localMigrations "github.com/e154/smart-home/internal/system/initial/local_migrations"
	"github.com/e154/smart-home/internal/system/logging_ws"
	"github.com/e154/smart-home/internal/system/rbac/access_list"
	"github.com/e154/smart-home/internal/system/terminal"
	"github.com/e154/smart-home/internal/system/validation"
	. "github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/common/encryptor"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/logger"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scheduler"
	"github.com/e154/smart-home/pkg/scripts"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/e154/bus"
	_ "github.com/e154/smart-home/internal/plugins"
	"github.com/e154/smart-home/version"
)

var (
	log = logger.MustGetLogger("initial")
)

// Initial ...
type Initial struct {
	adaptors        *Adaptors
	scriptService   scripts.ScriptService
	accessList      access_list.AccessListService
	supervisor      plugins.Supervisor
	automation      automation.Automation
	endpoint        *endpoint.Endpoint
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
	adaptors *Adaptors,
	scriptService scripts.ScriptService,
	accessList access_list.AccessListService,
	supervisor plugins.Supervisor,
	automation automation.Automation,
	endpoint *endpoint.Endpoint,
	api *api.Api,
	gateClient *client.GateClient,
	validation *validation.Validate,
	_ *logging_ws.LoggingWs,
	localMigrations *localMigrations.Migrations,
	demo *demo.Demos,
	_ scheduler.Scheduler,
	db *gorm.DB,
	eventBus bus.Bus,
	_ *terminal.Terminal) *Initial {
	initial := &Initial{
		adaptors:        adaptors,
		scriptService:   scriptService,
		accessList:      accessList,
		supervisor:      supervisor,
		automation:      automation,
		endpoint:        endpoint,
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

// InstallDemoData ...
func (n *Initial) InstallDemoData() {

	log.Info("install demo data")

	err := n.adaptors.Transaction.Do(context.Background(), func(ctx context.Context) error {
		if err := n.demo.InstallByName(ctx, "example1"); err != nil {
			log.Errorf("\n\nmigration '%s' ... error", err.Error())
			return err
		}
		return nil
	})

	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("complete")
}

// checkForUpgrade ...
func (n *Initial) checkForUpgrade() {

	const name = "initialVersion"
	v, err := n.adaptors.Variable.GetByName(context.Background(), name)
	if err != nil {

		if errors.Is(err, apperr.ErrNotFound) {
			v = m.Variable{
				Name:   name,
				Value:  "*local_migrations.MigrationInit",
				System: true,
			}
			err = n.adaptors.Variable.CreateOrUpdate(context.Background(), v)
			So(err, ShouldBeNil)
		}
	}

	oldVersion := v.Value
	So(err, ShouldBeNil)

	var currentVersion string
	defer func() {
		fmt.Println("")

		v.Value = currentVersion
		err = n.adaptors.Variable.CreateOrUpdate(context.Background(), v)
		So(err, ShouldBeNil)
	}()

	if currentVersion, err = n.localMigrations.Up(context.TODO(), n.adaptors, oldVersion); err != nil {
		return
	}
}

// Start ...
func (n *Initial) Start(ctx context.Context) (err error) {
	n.checkForUpgrade()

	variable, _ := n.adaptors.Variable.GetByName(ctx, "encryptor")
	val, _ := hex.DecodeString(variable.Value)
	encryptor.SetKey(val)

	_ = n.eventBus.Subscribe("system/models/variables/+", n.eventHandler)
	_ = n.eventBus.Subscribe("system", n.eventHandler)

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
	_ = n.eventBus.Unsubscribe("system", n.eventHandler)
	_ = n.api.Shutdown(ctx)
	return
}

func (n *Initial) eventHandler(_ string, message interface{}) {
	switch v := message.(type) {
	case events.EventGetServerVersion:
		n.eventBus.Publish("system", events.EventServerVersion{
			Common: events.Common{
				Owner:     events.OwnerSystem,
				SessionID: v.SessionID,
				User:      v.User,
			},
			Version: version.GetVersion(),
		})
	case events.EventUpdatedVariableModel:
		if v.Name == "timezone" && v.Value != "" {
			log.Infof("update database timezone to: \"%s\"", v.Value)
			n.db.Exec(`SET TIME ZONE '?';`, v.Value)
		}
	}
}
