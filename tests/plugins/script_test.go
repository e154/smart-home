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

package plugins

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/script"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/initial/local_migrations"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
)

func TestScript(t *testing.T) {

	const (
		plugActionOnOffSourceScript = `
entityAction = (entityId, actionName)->
  #print '---action on/off--'
  Done entityId, actionName
`
	)

	Convey("script", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			eventBus.Purge()
			scriptService.Purge()

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			err = local_migrations.NewMigrationPlugins(adaptors).Up(context.TODO(), nil)
			So(err, ShouldBeNil)

			go mqttServer.Start()

			// add scripts
			// ------------------------------------------------

			plugActionOnOffScript, err := AddScript("plug script", plugActionOnOffSourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			plugEnt := GetNewScript("script.1", []*m.Script{})
			plugEnt.PluginName = script.EntityScript
			plugEnt.Actions = []*m.EntityAction{
				{
					Name:        "ON",
					Description: "включить",
					Script:      plugActionOnOffScript,
				},
				{
					Name:        "OFF",
					Description: "выключить",
					Script:      plugActionOnOffScript,
				},
			}
			err = adaptors.Entity.Add(plugEnt)
			So(err, ShouldBeNil)

			// automation
			// ------------------------------------------------

			var counter atomic.Int32
			scriptService.PushFunctions("Done", func(entityId, action string) {
				counter.Inc()
			})

			pluginManager.Start()
			automation.Reload()
			entityManager.SetPluginManager(pluginManager)
			entityManager.LoadEntities()
			go zigbee2mqtt.Start()

			defer func() {
				_ = mqttServer.Shutdown()
				zigbee2mqtt.Shutdown()
				entityManager.Shutdown()
				_ = automation.Shutdown()
				pluginManager.Shutdown()
			}()

			entityManager.CallAction(plugEnt.Id, "ON", nil)
			entityManager.CallAction(plugEnt.Id, "OFF", nil)
			entityManager.CallAction(plugEnt.Id, "NULL", nil)
			time.Sleep(time.Millisecond * 500)

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 2)

		})
	})
}