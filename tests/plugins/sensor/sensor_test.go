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

package sensor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSensor(t *testing.T) {

	const sensorSourceScript = `
checkStatus =->
    res = HTTP.get("http://%s:%d/?t=12345678")
    if res.error 
        Actor.setState
            'new_state': 'ERROR'
        return
    p = unmarshal res.body
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
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "sensor")

			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

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
			err = adaptors.Entity.Add(context.Background(), sensorEnt)
			ctx.So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(context.Background(), &m.EntityStorage{
				EntityId:   sensorEnt.Id,
				Attributes: sensorEnt.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			eventBus.Publish("system/entities/"+sensorEnt.Id.String(), events.EventCreatedEntity{
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