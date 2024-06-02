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
	"sync"
	"testing"
	"time"

	"github.com/e154/bus"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	super "github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCommon(t *testing.T) {

	Convey("sensor", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor super.Supervisor,
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

			// add entity
			// ------------------------------------------------

			sensorEnt := GetNewSensor("device1")
			sensorEnt.Actions = []*m.EntityAction{
				{
					Name:        "ACTION1",
					Description: "action description",
				},
			}
			sensorEnt.States = []*m.EntityState{
				{
					Name:        "STATE1",
					Description: "state description",
				},
				{
					Name:        "STATE2",
					Description: "state description",
				},
			}
			sensorEnt.Attributes = m.Attributes{
				"v": {
					Name: "v",
					Type: common.AttributeString,
				},
			}
			sensorEnt.Settings = m.Attributes{
				"v": {
					Name: "v",
					Type: common.AttributeString,
				},
			}

			err := adaptors.Entity.Add(context.Background(), sensorEnt)
			ctx.So(err, ShouldBeNil)

			eventBus.Publish("system/models/entities/"+sensorEnt.Id.String(), events.EventCreatedEntityModel{
				EntityId: sensorEnt.Id,
			})

			//--

			sensor2Ent := GetNewSensor("device2")
			sensor2Ent.Attributes = m.Attributes{
				"v": {
					Name: "v",
					Type: common.AttributeString,
				},
			}

			err = adaptors.Entity.Add(context.Background(), sensor2Ent)
			ctx.So(err, ShouldBeNil)

			eventBus.Publish("system/models/entities/"+sensor2Ent.Id.String(), events.EventCreatedEntityModel{
				EntityId: sensor2Ent.Id,
			})
			//--

			sensor3Ent := GetNewSensor("device3")
			sensor3Ent.Attributes = m.Attributes{
				"v": {
					Name: "v",
					Type: common.AttributeString,
				},
			}

			err = adaptors.Entity.Add(context.Background(), sensor3Ent)
			ctx.So(err, ShouldBeNil)

			eventBus.Publish("system/models/entities/"+sensor3Ent.Id.String(), events.EventCreatedEntityModel{
				EntityId: sensor3Ent.Id,
			})
			//--

			time.Sleep(time.Second)

			t.Run("sensor states", func(t *testing.T) {
				Convey("states", t, func(ctx C) {

					// common
					// ------------------------------------------------
					ch := make(chan events.EventStateChanged, 99)
					defer close(ch)
					fn := func(_ string, msg interface{}) {
						switch v := msg.(type) {
						case events.EventStateChanged:
							ch <- v
						}
					}
					eventBus.Subscribe("system/entities/+", fn, false)
					defer eventBus.Unsubscribe("system/entities/+", fn)

					// 1
					// ------------------------------------------------
					ctx.Println("v1")
					supervisor.SetState(sensorEnt.Id, super.EntityStateParams{
						NewState: nil,
						AttributeValues: m.AttributeValue{
							"v": "V1",
						},
						StorageSave: false,
					})

					// wait message
					msg, ok := WaitT[events.EventStateChanged](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.StorageSave, ShouldBeFalse)
					ctx.So(msg.DoNotSaveMetric, ShouldBeFalse)
					ctx.So(msg.PluginName, ShouldEqual, "sensor")
					ctx.So(msg.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.OldState.EntityId, ShouldBeZeroValue)
					ctx.So(msg.OldState.Value, ShouldBeNil)
					ctx.So(msg.OldState.State, ShouldBeNil)
					ctx.So(msg.OldState.Attributes, ShouldBeNil)
					ctx.So(msg.OldState.Settings, ShouldBeNil)
					ctx.So(msg.OldState.LastChanged, ShouldBeNil)
					ctx.So(msg.OldState.LastUpdated, ShouldBeNil)
					ctx.So(msg.NewState.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.NewState.Value, ShouldBeNil)
					ctx.So(msg.NewState.State, ShouldBeNil)
					ctx.So(msg.NewState.Attributes, ShouldNotBeNil)
					ctx.So(msg.NewState.Attributes["v"].String(), ShouldEqual, "V1")
					ctx.So(msg.NewState.Settings, ShouldNotBeNil)
					ctx.So(msg.NewState.Settings["v"].String(), ShouldBeZeroValue)
					ctx.So(msg.NewState.LastChanged, ShouldBeNil)
					ctx.So(msg.NewState.LastUpdated, ShouldNotBeNil)

					v1NewStateLastChange := msg.NewState.LastChanged // null
					v1NewStateLastUpdate := msg.NewState.LastUpdated

					time.Sleep(time.Millisecond * 500)

					// 2
					// ------------------------------------------------
					ctx.Println("\nv2")
					supervisor.SetState(sensorEnt.Id, super.EntityStateParams{
						NewState: nil,
						AttributeValues: m.AttributeValue{
							"v": "V2",
						},
						StorageSave: true,
					})

					// wait message
					msg, ok = WaitT[events.EventStateChanged](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.StorageSave, ShouldBeTrue)
					ctx.So(msg.DoNotSaveMetric, ShouldBeFalse)
					ctx.So(msg.PluginName, ShouldEqual, "sensor")
					ctx.So(msg.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.OldState.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.OldState.Value, ShouldBeNil)
					ctx.So(msg.OldState.State, ShouldBeNil)
					ctx.So(msg.OldState.Attributes, ShouldNotBeNil)
					ctx.So(msg.OldState.Attributes["v"].String(), ShouldEqual, "V1")
					ctx.So(msg.OldState.Settings, ShouldNotBeNil)
					ctx.So(msg.OldState.Settings["v"].String(), ShouldBeZeroValue)
					ctx.So(msg.OldState.LastChanged, ShouldBeNil)
					ctx.So(msg.OldState.LastUpdated, ShouldNotBeNil)
					ctx.So(msg.NewState.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.NewState.Value, ShouldBeNil)
					ctx.So(msg.NewState.State, ShouldBeNil)
					ctx.So(msg.NewState.Attributes, ShouldNotBeNil)
					ctx.So(msg.NewState.Attributes["v"].String(), ShouldEqual, "V2")
					ctx.So(msg.NewState.Settings, ShouldNotBeNil)
					ctx.So(msg.NewState.Settings["v"].String(), ShouldBeZeroValue)
					ctx.So(msg.NewState.LastChanged, ShouldNotBeNil)
					ctx.So(msg.NewState.LastUpdated, ShouldNotBeNil)

					v2NewStateLastChange := msg.NewState.LastChanged
					v2NewStateLastUpdate := msg.NewState.LastUpdated

					v2NewStateLastChanged := msg.NewState.LastChanged
					ctx.So(v2NewStateLastChanged.Sub(*msg.OldState.LastUpdated), ShouldEqual, 0)
					ctx.So(v2NewStateLastChanged.Sub(*v1NewStateLastUpdate), ShouldEqual, 0)

					v2OldStateLastChanged := msg.OldState.LastChanged // null
					v2OldStateLastUpdate := msg.OldState.LastUpdated
					ctx.So(v2OldStateLastChanged, ShouldBeNil)
					ctx.So(v1NewStateLastChange, ShouldBeNil)
					ctx.So(fmt.Sprintf("%v", v2OldStateLastUpdate), ShouldEqual, fmt.Sprintf("%v", v1NewStateLastUpdate))

					time.Sleep(time.Millisecond * 500)

					// 3
					// ------------------------------------------------
					ctx.Println("\nv3")
					supervisor.SetState(sensorEnt.Id, super.EntityStateParams{
						NewState: common.String("STATE1"),
						AttributeValues: m.AttributeValue{
							"v": "V3",
						},
						StorageSave: true,
					})

					// wait message
					msg, ok = WaitT[events.EventStateChanged](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.StorageSave, ShouldBeTrue)
					ctx.So(msg.DoNotSaveMetric, ShouldBeFalse)
					ctx.So(msg.PluginName, ShouldEqual, "sensor")
					ctx.So(msg.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.OldState.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.OldState.Value, ShouldBeNil)
					ctx.So(msg.OldState.State, ShouldBeNil)
					ctx.So(msg.OldState.Attributes, ShouldNotBeNil)
					ctx.So(msg.OldState.Attributes["v"].String(), ShouldEqual, "V2")
					ctx.So(msg.OldState.Settings, ShouldNotBeNil)
					ctx.So(msg.OldState.Settings["v"].String(), ShouldBeZeroValue)
					ctx.So(msg.OldState.LastChanged, ShouldNotBeNil)
					ctx.So(msg.OldState.LastUpdated, ShouldNotBeNil)
					ctx.So(msg.NewState.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.NewState.Value, ShouldBeNil)
					ctx.So(msg.NewState.State, ShouldNotBeNil)
					ctx.So(msg.NewState.State.Name, ShouldEqual, "STATE1")
					ctx.So(msg.NewState.State.Description, ShouldEqual, "state description")
					ctx.So(msg.NewState.State.ImageUrl, ShouldBeNil)
					ctx.So(msg.NewState.State.Icon, ShouldBeNil)
					ctx.So(msg.NewState.Attributes, ShouldNotBeNil)
					ctx.So(msg.NewState.Attributes["v"].String(), ShouldEqual, "V3")
					ctx.So(msg.NewState.Settings, ShouldNotBeNil)
					ctx.So(msg.NewState.Settings["v"].String(), ShouldBeZeroValue)
					ctx.So(msg.NewState.LastChanged, ShouldNotBeNil)
					ctx.So(msg.NewState.LastUpdated, ShouldNotBeNil)

					v3NewStateLastChange := msg.NewState.LastChanged
					v3NewStateLastUpdate := msg.NewState.LastUpdated

					v3NewStateLastChanged := msg.NewState.LastChanged
					ctx.So(v3NewStateLastChanged.Sub(*msg.OldState.LastUpdated), ShouldEqual, 0)
					ctx.So(v3NewStateLastChanged.Sub(*v2NewStateLastUpdate), ShouldEqual, 0)

					v3OldStateLastChanged := msg.OldState.LastChanged
					v3OldStateLastUpdate := msg.OldState.LastUpdated
					ctx.So(v3OldStateLastChanged.Sub(*v2NewStateLastChange), ShouldEqual, 0)
					ctx.So(v3OldStateLastUpdate.Sub(*v2NewStateLastUpdate), ShouldEqual, 0)

					time.Sleep(time.Millisecond * 500)

					// 4 (skip settings)
					// ------------------------------------------------
					ctx.Println("\nv4")
					supervisor.SetState(sensorEnt.Id, super.EntityStateParams{
						SettingsValue: m.AttributeValue{
							"v": "V4",
						},
						StorageSave: true,
					})

					// wait message
					msg, ok = WaitT[events.EventStateChanged](time.Second*2, ch)
					ctx.So(ok, ShouldBeFalse)

					// 5 (skip unknown fields)
					// ------------------------------------------------
					ctx.Println("\nv5")
					supervisor.SetState(sensorEnt.Id, super.EntityStateParams{
						SettingsValue: m.AttributeValue{
							"foo": "bar",
						},
						StorageSave: true,
					})

					// wait message
					msg, ok = WaitT[events.EventStateChanged](time.Second*2, ch)
					ctx.So(ok, ShouldBeFalse)

					time.Sleep(time.Millisecond * 500)

					// 6
					// ------------------------------------------------
					ctx.Println("\nv6")
					supervisor.SetState(sensorEnt.Id, super.EntityStateParams{
						NewState: common.String("FOO"),
						AttributeValues: m.AttributeValue{
							"v": "V6",
						},
						StorageSave: true,
					})

					// wait message
					msg, ok = WaitT[events.EventStateChanged](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.StorageSave, ShouldBeTrue)
					ctx.So(msg.DoNotSaveMetric, ShouldBeFalse)
					ctx.So(msg.PluginName, ShouldEqual, "sensor")
					ctx.So(msg.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.OldState.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.OldState.Value, ShouldBeNil)
					ctx.So(msg.OldState.State, ShouldNotBeNil)
					ctx.So(msg.OldState.State.Name, ShouldEqual, "STATE1")
					ctx.So(msg.OldState.Attributes, ShouldNotBeNil)
					ctx.So(msg.OldState.Attributes["v"].String(), ShouldEqual, "V3")
					ctx.So(msg.OldState.Settings, ShouldNotBeNil)
					ctx.So(msg.OldState.Settings["v"].String(), ShouldBeZeroValue)
					ctx.So(msg.OldState.LastChanged, ShouldNotBeNil)
					ctx.So(msg.OldState.LastUpdated, ShouldNotBeNil)
					ctx.So(msg.NewState.EntityId, ShouldEqual, "sensor.device1")
					ctx.So(msg.NewState.Value, ShouldBeNil)
					ctx.So(msg.NewState.State, ShouldNotBeNil)
					ctx.So(msg.NewState.State.Name, ShouldEqual, "STATE1")
					ctx.So(msg.NewState.State.Description, ShouldEqual, "state description")
					ctx.So(msg.NewState.State.Icon, ShouldBeNil)
					ctx.So(msg.NewState.Attributes, ShouldNotBeNil)
					ctx.So(msg.NewState.Attributes["v"].String(), ShouldEqual, "V6")
					ctx.So(msg.NewState.Settings, ShouldNotBeNil)
					ctx.So(msg.NewState.Settings["v"].String(), ShouldBeZeroValue)
					ctx.So(msg.NewState.LastChanged, ShouldNotBeNil)
					ctx.So(msg.NewState.LastUpdated, ShouldNotBeNil)

					v6NewStateLastChanged := msg.NewState.LastChanged
					ctx.So(v6NewStateLastChanged.Sub(*msg.OldState.LastUpdated), ShouldEqual, 0)
					ctx.So(v6NewStateLastChanged.Sub(*v3NewStateLastUpdate), ShouldEqual, 0)

					v6OldStateLastChanged := msg.OldState.LastChanged
					v6OldStateLastUpdate := msg.OldState.LastUpdated
					ctx.So(v6OldStateLastChanged.Sub(*v3NewStateLastChange), ShouldEqual, 0)
					ctx.So(v6OldStateLastUpdate.Sub(*v3NewStateLastUpdate), ShouldEqual, 0)

					v6NewStateLastUpdate := msg.NewState.LastUpdated

					time.Sleep(time.Second)

					list, total, err := adaptors.EntityStorage.List(context.Background(), 10, 0, "asc", "id", []common.EntityId{}, nil, nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 3)

					ctx.So(list[0].EntityId, ShouldEqual, "sensor.device1")
					ctx.So(list[0].Attributes["v"], ShouldEqual, "V2")
					ctx.So(list[0].State, ShouldBeZeroValue)
					ctx.So(v2NewStateLastUpdate.Sub(list[0].CreatedAt), ShouldEqual, 0)

					ctx.So(list[1].EntityId, ShouldEqual, "sensor.device1")
					ctx.So(list[1].Attributes["v"], ShouldEqual, "V3")
					ctx.So(list[1].State, ShouldEqual, "STATE1")
					ctx.So(v3NewStateLastUpdate.Sub(list[1].CreatedAt), ShouldEqual, 0)

					ctx.So(list[2].EntityId, ShouldEqual, "sensor.device1")
					ctx.So(list[2].Attributes["v"], ShouldEqual, "V6")
					ctx.So(list[2].State, ShouldEqual, "STATE1")
					ctx.So(v6NewStateLastUpdate.Sub(list[2].CreatedAt), ShouldEqual, 0)
				})
			})

			t.Run("sensor states2", func(t *testing.T) {
				Convey("states", t, func(ctx C) {

					// common
					// ------------------------------------------------
					wg := &sync.WaitGroup{}
					wg.Add(100)
					fn := func(_ string, msg interface{}) {
						switch msg.(type) {
						case events.EventStateChanged:
							wg.Done()
						}
					}
					eventBus.Subscribe("system/entities/sensor.device2", fn, false)
					defer eventBus.Unsubscribe("system/entities/sensor.device2", fn)

					// 7
					// ------------------------------------------------
					ctx.Println("\nv7")
					for i := 0; i < 100; i++ {
						supervisor.SetState(sensor2Ent.Id, super.EntityStateParams{
							AttributeValues: m.AttributeValue{
								"v": fmt.Sprintf("V%d", i),
							},
							StorageSave: true,
						})
					}

					ok := WaitGroupTimeout(wg, time.Second*3)
					ctx.So(ok, ShouldBeTrue)

					time.Sleep(time.Second)

					_, total, err := adaptors.EntityStorage.List(context.Background(), 500, 0, "asc", "id", []common.EntityId{sensor2Ent.Id}, nil, nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 100)
				})
			})

			t.Run("sensor states3", func(t *testing.T) {
				Convey("states", t, func(ctx C) {

					// common
					// ------------------------------------------------
					wg := &sync.WaitGroup{}
					wg.Add(1000000)
					fn := func(_ string, msg interface{}) {
						switch msg.(type) {
						case events.EventStateChanged:
							wg.Done()
						}
					}
					eventBus.Subscribe("system/entities/sensor.device3", fn, false)
					defer eventBus.Unsubscribe("system/entities/sensor.device3", fn)

					// 8
					// ------------------------------------------------
					ctx.Println("\nv8")
					for i := 0; i < 1000000; i++ {
						supervisor.SetState(sensor3Ent.Id, super.EntityStateParams{
							AttributeValues: m.AttributeValue{
								"v": fmt.Sprintf("V%d", i),
							},
						})
					}

					ok := WaitGroupTimeout(wg, time.Second*3)
					ctx.So(ok, ShouldBeTrue)

					time.Sleep(time.Second)

					_, total, err := adaptors.EntityStorage.List(context.Background(), 25, 0, "asc", "id", []common.EntityId{sensor3Ent.Id}, nil, nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 0)
				})
			})
		})
	})
}
