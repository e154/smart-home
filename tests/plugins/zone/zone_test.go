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

package zone

import (
	"context"
	"github.com/e154/smart-home/common/events"
	"sync"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestZone(t *testing.T) {

	Convey("zone", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			visor supervisor.Supervisor,
			automation automation.Automation,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "zone")

			// add entity
			// ------------------------------------------------
			zoneEnt := GetNewZone()
			err := adaptors.Entity.Add(zoneEnt)
			So(err, ShouldBeNil)

			eventBus.Publish("system/entities/"+zoneEnt.Id.String(), events.EventCreatedEntity{
				EntityId: zoneEnt.Id,
			})

			time.Sleep(time.Millisecond * 500)

			// ------------------------------------------------
			wgAdd := sync.WaitGroup{}
			wgAdd.Add(1)
			wgUpdate := sync.WaitGroup{}
			wgUpdate.Add(1)
			fn := func(_ string, msg interface{}) {

				switch v := msg.(type) {
				case events.EventStateChanged:
					if v.PluginName != "zone" {
						return
					}

					//settings := v.NewState.Settings
					//ctx.So(settings["elevation"].Value, ShouldEqual, 10)
					//ctx.So(settings["lat"].Value, ShouldEqual, 10.881)
					//ctx.So(settings["lon"].Value, ShouldEqual, 107.570)
					//ctx.So(settings["timezone"].Value, ShouldEqual, 7)
					wgUpdate.Done()

				case events.EventAddedActor:
					if v.PluginName != "zone" {
						return
					}

					//settings := v.Settings
					//ctx.So(settings["elevation"].Value, ShouldEqual, 150)
					//ctx.So(settings["lat"].Value, ShouldEqual, 54.9022)
					//ctx.So(settings["lon"].Value, ShouldEqual, 83.0335)
					//ctx.So(settings["timezone"].Value, ShouldEqual, 7)
					wgAdd.Done()

				default:
					//fmt.Printf("new event: %v\n", reflect.TypeOf(v).String())
				}
			}
			_ = eventBus.Subscribe("system/entities/+", fn)
			defer func() {
				_ = eventBus.Unsubscribe("system/entities/+", fn)
			}()

			// ------------------------------------------------

			visor.Start(context.Background())
			automation.Start()

			//...
			wgAdd.Wait()
			_ = visor.SetState(zoneEnt.Id, supervisor.EntityStateParams{
				SettingsValue: m.AttributeValue{
					"elevation": 10,
					"lat":       10.881,
					"lon":       107.570,
					"timezone":  7,
				},
			})

			wgUpdate.Wait()
		})
	})
}
