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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/initial/env1"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
)

func TestNode(t *testing.T) {

	Convey("node", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			env1.NewPluginManager(adaptors).Create()

			go mqttServer.Start()

			// add entity
			// ------------------------------------------------

			nodeEnt := GetNewNode()
			err = adaptors.Entity.Add(nodeEnt)
			So(err, ShouldBeNil)

			// ------------------------------------------------
			wgAdd := sync.WaitGroup{}
			wgAdd.Add(1)
			wgUpdate := sync.WaitGroup{}
			wgUpdate.Add(1)
			eventBus.Subscribe(event_bus.TopicEntities, func(msg interface{}) {

			})

			// ------------------------------------------------
			pluginManager.Start()
			automation.Reload()
			entityManager.LoadEntities(pluginManager)
			go zigbee2mqtt.Start()

			//...
			wgAdd.Wait()
			entityManager.SetState(nodeEnt.Id, entity_manager.EntityStateParams{
				AttributeValues: m.EntityAttributeValue{
					"elevation": 10,
					"lat":       10.881,
					"lon":       107.570,
					"timezone":  7,
				},
			})

			wgUpdate.Wait()
			time.Sleep(time.Millisecond * 500)

			mqttServer.Shutdown()
			zigbee2mqtt.Shutdown()
			entityManager.Shutdown()
			automation.Shutdown()
			pluginManager.Shutdown()
		})
	})
}
