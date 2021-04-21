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

package plugins

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/plugin_manager"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"
	"testing"
	"time"
)

func TestScene(t *testing.T) {

	const (
		sceneSourceScript = `
sceneEvent = (args)->
  #print '---action on/off--'
  Done args
`
	)

	Convey("scene", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			plugins *plugins.Loader,
			entityManager *entity_manager.EntityManager,
			zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt,
			mqttServer *mqtt.Mqtt,
			automation *automation.Automation,
			eventBus *event_bus.EventBus,
			pluginManager *plugin_manager.PluginManager) {

			err := migrations.Purge()
			So(err, ShouldBeNil)

			go mqttServer.Start()

			// add scripts
			// ------------------------------------------------

			sceneScript := &m.Script{
				Lang:        common.ScriptLangCoffee,
				Name:        "romantic",
				Source:      sceneSourceScript,
				Description: "romantic scene",
			}

			engineSceneScript, err := scriptService.NewEngine(sceneScript)
			So(err, ShouldBeNil)
			err = engineSceneScript.Compile()
			So(err, ShouldBeNil)

			sceneScript.Id, err = adaptors.Script.Add(sceneScript)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			romanticEnt := GetNewScene(fmt.Sprintf("scene.romantic"), []m.Script{*sceneScript,})
			err = adaptors.Entity.Add(romanticEnt)
			So(err, ShouldBeNil)

			// automation
			// ------------------------------------------------

			var counter atomic.Int32
			scriptService.PushFunctions("Done", func(args string) {
				counter.Inc()
			})

			plugins.Register()
			pluginManager.Start()
			automation.Reload()
			entityManager.LoadEntities()
			go zigbee2mqtt.Start()

			time.Sleep(time.Millisecond * 500)
			entityManager.CallScene(romanticEnt.Id, nil)
			time.Sleep(time.Millisecond * 500)

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 1)

			mqttServer.Shutdown()
			zigbee2mqtt.Shutdown()
			entityManager.Shutdown()
			automation.Shutdown()
			pluginManager.Shutdown()
		})
	})
}
