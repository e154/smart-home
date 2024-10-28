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

package sensor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/e154/smart-home/pkg/adaptors"
	commonPkg "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/e154/bus"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSensor(t *testing.T) {

	const sensorSourceScript = `
checkStatus =->
    res = HTTP.get("http://%s:%d/?t=12345678")
    if res.error 
        EntitySetStateName ENTITY_ID, 'ERROR'
        return
    p = unmarshal res.body
    attrs =
        paid_rewards: p.user.paid_rewards
    EntitySetState ENTITY_ID,
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
			scriptService scripts.ScriptService,
			supervisor plugins.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "sensor")

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor")
			pluginsCh := WaitPlugins(eventBus, time.Second*5, "sensor")
			supervisor.Start(context.Background())
			defer supervisor.Shutdown(context.Background())
			So(<-serviceCh, ShouldBeTrue)
			So(<-pluginsCh, ShouldBeTrue)

			// bind convey
			RegisterConvey(scriptService, ctx)

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
			sensorEnt.Actions = []*models.EntityAction{
				{
					Name:        "CHECK",
					Description: "condition check",
					Script:      plugScript,
				},
			}
			sensorEnt.States = []*models.EntityState{
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
			sensorEnt.Attributes = models.Attributes{
				"paid_rewards": {
					Name: "paid_rewards",
					Type: commonPkg.AttributeFloat,
				},
			}
			err = adaptors.Entity.Add(context.Background(), sensorEnt)
			ctx.So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(context.Background(), &models.EntityStorage{
				EntityId:   sensorEnt.Id,
				Attributes: sensorEnt.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			eventBus.Publish("system/models/entities/"+sensorEnt.Id.String(), events.EventCreatedEntityModel{
				EntityId: sensorEnt.Id,
			})

			time.Sleep(time.Second)

			// ------------------------------------------------

			t.Run("sensor stats", func(t *testing.T) {
				Convey("stats", t, func(ctx C) {

					ctx2, cancel := context.WithCancel(context.Background())
					go func() { _ = MockHttpServer(ctx2, host, port, []byte(response1)) }()
					time.Sleep(time.Millisecond * 500)
					supervisor.CallAction(sensorEnt.Id, "CHECK", nil)
					time.Sleep(time.Second)
					cancel()

					lastState, err := adaptors.EntityStorage.GetLastByEntityId(context.Background(), sensorEnt.Id)
					ctx.So(err, ShouldBeNil)

					ctx.So(lastState.Attributes["paid_rewards"], ShouldEqual, 9.08)
				})
			})

		})
	})
}
