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

func TestZone(t *testing.T) {

	Convey("zone", t, func(ctx C) {
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

			zoneEnt := GetNewZone()
			err = adaptors.Entity.Add(zoneEnt)
			So(err, ShouldBeNil)

			serial := zoneEnt.Attributes.Serialize()
			_, err = adaptors.EntityStorage.Add(m.EntityStorage{
				EntityId:   "zone.home",
				Attributes: serial,
			})
			So(err, ShouldBeNil)

			// ------------------------------------------------
			wgAdd := sync.WaitGroup{}
			wgAdd.Add(1)
			wgUpdate := sync.WaitGroup{}
			wgUpdate.Add(1)
			eventBus.Subscribe(event_bus.TopicEntities, func(_ string, msg interface{}) {

				switch v := msg.(type) {
				case event_bus.EventStateChanged:
					if v.Type != "zone" {
						return
					}

					attr := v.NewState.Attributes
					ctx.So(attr["elevation"].Value, ShouldEqual, 10)
					ctx.So(attr["lat"].Value, ShouldEqual, 10.881)
					ctx.So(attr["lon"].Value, ShouldEqual, 107.570)
					ctx.So(attr["timezone"].Value, ShouldEqual, 7)
					wgUpdate.Done()

				case event_bus.EventAddedNewEntity:
					if v.Type != "zone" {
						return
					}

					attr := v.Attributes
					ctx.So(attr["elevation"].Value, ShouldEqual, 150)
					ctx.So(attr["lat"].Value, ShouldEqual, 54.9022)
					ctx.So(attr["lon"].Value, ShouldEqual, 83.0335)
					ctx.So(attr["timezone"].Value, ShouldEqual, 7)
					wgAdd.Done()

				default:
					//fmt.Printf("new event: %v\n", reflect.TypeOf(v).String())
				}
			})

			// ------------------------------------------------

			pluginManager.Start()
			automation.Reload()
			entityManager.LoadEntities(pluginManager)
			go zigbee2mqtt.Start()

			defer func() {
				mqttServer.Shutdown()
				zigbee2mqtt.Shutdown()
				entityManager.Shutdown()
				automation.Shutdown()
				pluginManager.Shutdown()
			}()

			//...
			wgAdd.Wait()
			entityManager.SetState(zoneEnt.Id, entity_manager.EntityStateParams{
				AttributeValues: m.EntityAttributeValue{
					"elevation": 10,
					"lat":       10.881,
					"lon":       107.570,
					"timezone":  7,
				},
			})

			wgUpdate.Wait()
			time.Sleep(time.Millisecond * 500)
		})
	})
}
