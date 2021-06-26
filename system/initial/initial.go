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
	. "github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api"

	//"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/gate_client"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/initial/env1"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/fx"
	"strconv"

	_ "github.com/e154/smart-home/plugins"
)

var (
	log = common.MustGetLogger("initial")
)

var (
	currentVersion = 3
)

// Initial ...
type Initial struct {
	migrations    *migrations.Migrations
	adaptors      *Adaptors
	scriptService scripts.ScriptService
	accessList    access_list.AccessListService
	entityManager entity_manager.EntityManager
	pluginManager common.PluginManager
	automation    automation.Automation
	api           *api.Api
	metrics       *metrics.MetricManager
	gateClient    *gate_client.GateClient
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
	metrics *metrics.MetricManager,
	gateClient *gate_client.GateClient) *Initial {
	initial := &Initial{
		migrations:    migrations,
		adaptors:      adaptors,
		scriptService: scriptService,
		accessList:    accessList,
		entityManager: entityManager,
		pluginManager: pluginManager,
		automation:    automation,
		api:           api,
		metrics:       metrics,
		gateClient:    gateClient,
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

	n.migrations.Purge()

	log.Info("complete")
}

// InstallDemoData ...
func (n *Initial) InstallDemoData() {

	log.Info("install demo data")

	tx := n.adaptors.Begin()

	env1.InstallDemoData(tx, n.accessList, n.scriptService)

	err := tx.Variable.Add(m.Variable{
		Name:  "initial_version",
		Value: fmt.Sprintf("%d", currentVersion),
	})
	if err != nil {
		tx.Rollback()
	}
	So(err, ShouldBeNil)

	tx.Commit()

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

		if errors.Is(err, ErrRecordNotFound) {
			v = m.Variable{
				Name:  "initial_version",
				Value: fmt.Sprintf("%d", 1),
			}
			err = tx.Variable.Add(v)
			So(err, ShouldBeNil)
		}

		// create
		env1.Create(tx, n.accessList, n.scriptService)
	}

	oldVersion, err := strconv.Atoi(v.Value)
	So(err, ShouldBeNil)

	// upgrade
	env1.Upgrade(oldVersion, tx, n.accessList, n.scriptService)

	if oldVersion >= currentVersion {
		tx.Commit()
		return
	}

	v.Value = fmt.Sprintf("%d", currentVersion)
	err = tx.Variable.Update(v)
	So(err, ShouldBeNil)

	tx.Commit()
}

// Start ...
func (n *Initial) Start() {

	n.checkForUpgrade()
	n.metrics.Start()
	n.pluginManager.Start()
	n.entityManager.LoadEntities(n.pluginManager)
	n.automation.Start()
	n.api.Start()
	n.gateClient.Start()
}

// Shutdown ...
func (n *Initial) Shutdown() {

}
