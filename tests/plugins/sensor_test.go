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
	context2 "context"
	"fmt"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSensor(t *testing.T) {

	const sensorSourceScript = `
checkStatus =->
    res = http.get("http://%s:%d/?t=12345678")
    if res.error 
        Actor.setState
            'new_state': 'ERROR'
        return
    p = JSON.parse(res.body)
    attrs =
        paid_rewards: p.user.paid_rewards

    Actor.setState
        new_state: 'ENABLED'
        attribute_values: attrs
        storage_save: true

entityAction = (entityId, actionName)->
    switch actionName
        when 'CHECK' then checkStatus()
`

	const response1 = `{"user": {"paid_rewards": 9.08}}`

	Convey("sensor", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			eventBus.Purge()
			scriptService.Purge()

			err := migrations.Purge()
			ctx.So(err, ShouldBeNil)

			// bind convey
			RegisterConvey(scriptService, ctx)

			// register plugins
			err = AddPlugin(adaptors, "sensor")
			ctx.So(err, ShouldBeNil)

			// test server
			var port = GetPort()
			var host = "127.0.0.1"

			// add scripts
			// ------------------------------------------------

			plugScript, err := AddScript("sensor script", fmt.Sprintf(sensorSourceScript, host, port), adaptors, scriptService)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------

			sensorEnt := GetNewSensor("device1")
			sensorEnt.Actions = []*m.EntityAction{
				{
					Name:        "CHECK",
					Description: "condition check",
					Script:      plugScript,
				},
			}
			sensorEnt.States = []*m.EntityState{
				{
					Name:        "ENABLED",
					Description: "enabled state",
				},
				{
					Name:        "DISABLED",
					Description: "disabled state",
				},
				{
					Name:        "ERROR",
					Description: "error state",
				},
			}
			sensorEnt.Attributes = m.Attributes{
				"paid_rewards": {
					Name: "paid_rewards",
					Type: common.AttributeFloat,
				},
			}
			err = adaptors.Entity.Add(sensorEnt)
			ctx.So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(&m.EntityStorage{
				EntityId:   sensorEnt.Id,
				Attributes: sensorEnt.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			// ------------------------------------------------
			pluginManager.Start()
			entityManager.SetPluginManager(pluginManager)
			entityManager.LoadEntities()

			defer func() {
				entityManager.Shutdown()
				pluginManager.Shutdown()
			}()

			time.Sleep(time.Millisecond * 500)

			t.Run("sensor stats", func(t *testing.T) {
				Convey("stats", t, func(ctx C) {

					ctx2, cancel := context2.WithCancel(context2.Background())
					go func() { _ = MockHttpServer(ctx2, host, port, []byte(response1)) }()
					time.Sleep(time.Millisecond * 500)
					entityManager.CallAction(sensorEnt.Id, "CHECK", nil)
					time.Sleep(time.Second)
					cancel()

					lastState, err := adaptors.EntityStorage.GetLastByEntityId(sensorEnt.Id)
					ctx.So(err, ShouldBeNil)

					ctx.So(lastState.Attributes["paid_rewards"], ShouldEqual, 9.08)
				})
			})

		})
	})
}
