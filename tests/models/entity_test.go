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

package models

import (
	"context"
	"testing"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/migrations"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEntity(t *testing.T) {
	Convey("entity", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			accessList access_list.AccessListService) {

			// clear database
			_ = migrations.Purge()

			// add scripts
			script1 := &m.Script{
				Lang:   common.ScriptLangCoffee,
				Name:   "script1",
				Source: "print 'OK'",
			}
			script2 := &m.Script{
				Lang:   common.ScriptLangCoffee,
				Name:   "script2",
				Source: "print 'OK'",
			}
			var err error
			script1.Id, err = adaptors.Script.Add(context.Background(), script1)
			So(err, ShouldBeNil)
			script2.Id, err = adaptors.Script.Add(context.Background(), script2)
			So(err, ShouldBeNil)

			// plugins
			err = AddPlugin(adaptors, "sensor")
			So(err, ShouldBeNil)

			// add image
			image1 := &m.Image{
				Url:  "foo",
				Name: "foo",
			}
			image2 := &m.Image{
				Url:  "bar",
				Name: "bar",
			}
			image1.Id, err = adaptors.Image.Add(context.Background(), image1)
			So(err, ShouldBeNil)
			image2.Id, err = adaptors.Image.Add(context.Background(), image2)
			So(err, ShouldBeNil)

			t.Run("Create", func(t *testing.T) {
				Convey("", t, func(ctx C) {
					// entity
					entity := &m.Entity{
						Id:         common.EntityId("sensor.entity1"),
						PluginName: "sensor",
						AutoLoad:   true,
						ImageId:    common.Int64(image1.Id),
						Scripts: []*m.Script{
							script1,
						},
						Actions: []*m.EntityAction{
							{
								Name:     "ACTION1",
								ScriptId: common.Int64(script1.Id),
								ImageId:  common.Int64(image2.Id),
							},
						},
						States: []*m.EntityState{
							{
								Name: "STATE1",
							},
						},
						Metrics: []*m.Metric{
							{
								Name:        "bar",
								Description: "bar",
								Options: m.MetricOptions{
									Items: []m.MetricOptionsItem{
										{
											Name:        "foo",
											Description: "foo",
											Color:       "foo",
											Translate:   "foo",
											Label:       "foo",
										},
									},
								},
							},
						},
						Attributes: NetAttr(),
						Settings:   NetSettings(),
					}
					err = adaptors.Entity.Add(context.Background(), entity)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), entity.Id)
					So(err, ShouldBeNil)

					So(entity.AutoLoad, ShouldBeTrue)
					So(len(entity.Actions), ShouldEqual, 1)
					So(entity.Actions[0].Name, ShouldEqual, "ACTION1")
					So(entity.Actions[0].Image, ShouldNotBeNil)
					So(entity.Actions[0].Image.Name, ShouldEqual, "bar")
					So(entity.Image, ShouldNotBeNil)
					So(entity.Image.Name, ShouldEqual, "foo")
					So(len(entity.States), ShouldEqual, 1)
					So(entity.States[0].Name, ShouldEqual, "STATE1")
					So(len(entity.Scripts), ShouldEqual, 1)
					So(entity.Scripts[0].Name, ShouldEqual, "script1")
					So(len(entity.Metrics), ShouldEqual, 1)
					So(entity.Metrics[0].Name, ShouldEqual, "bar")
					So(entity.Metrics[0].Options.Items[0].Name, ShouldEqual, "foo")
					So(entity.Settings, ShouldNotBeEmpty)
					So(entity.Settings["s"].String(), ShouldEqual, "s")
					So(entity.Attributes, ShouldNotBeEmpty)
					So(entity.Attributes["s"].Name, ShouldEqual, "s")
					So(entity.Attributes["s"].String(), ShouldEqual, "")
				})
			})

			t.Run("Update", func(t *testing.T) {
				Convey("", t, func(ctx C) {
					// entity
					entity := &m.Entity{
						Id:         common.EntityId("sensor.entity1"),
						AutoLoad:   true,
						PluginName: "sensor",
						Scripts: []*m.Script{
							script2,
						},
						Actions: []*m.EntityAction{
							{
								Name:     "ACTION2",
								ScriptId: common.Int64(script1.Id),
							},
						},
						States: []*m.EntityState{
							{
								Name: "STATE2",
							},
						},
						Metrics: []*m.Metric{
							{
								Name:        "bar2",
								Description: "bar2",
								Options: m.MetricOptions{
									Items: []m.MetricOptionsItem{
										{
											Name:        "foo2",
											Description: "foo2",
											Color:       "foo2",
											Translate:   "foo2",
											Label:       "foo2",
										},
									},
								},
							},
						},
					}
					err = adaptors.Entity.Update(context.Background(), entity)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), entity.Id)
					So(err, ShouldBeNil)

					So(len(entity.Actions), ShouldEqual, 1)
					So(entity.Actions[0].Name, ShouldEqual, "ACTION2")
					So(len(entity.States), ShouldEqual, 1)
					So(entity.States[0].Name, ShouldEqual, "STATE2")
					So(len(entity.Scripts), ShouldEqual, 1)
					So(entity.Scripts[0].Name, ShouldEqual, "script2")
					So(len(entity.Metrics), ShouldEqual, 1)
					So(entity.Metrics[0].Name, ShouldEqual, "bar2")
					So(entity.Metrics[0].Options.Items[0].Name, ShouldEqual, "foo2")

					// v2
					entity.Actions = []*m.EntityAction{}

					err = adaptors.Entity.Update(context.Background(), entity)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), entity.Id)
					So(err, ShouldBeNil)

					So(len(entity.Actions), ShouldEqual, 0)
					So(len(entity.States), ShouldEqual, 1)
					So(entity.States[0].Name, ShouldEqual, "STATE2")
					So(len(entity.Scripts), ShouldEqual, 1)
					So(entity.Scripts[0].Name, ShouldEqual, "script2")
					So(len(entity.Metrics), ShouldEqual, 1)
					So(entity.Metrics[0].Name, ShouldEqual, "bar2")
					So(entity.Metrics[0].Options.Items[0].Name, ShouldEqual, "foo2")

					// v3
					entity.Actions = []*m.EntityAction{
						{
							Name:     "ACTION2",
							ScriptId: common.Int64(script1.Id),
						},
					}
					entity.States = []*m.EntityState{}
					err = adaptors.Entity.Update(context.Background(), entity)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), entity.Id)
					So(err, ShouldBeNil)

					So(len(entity.Actions), ShouldEqual, 1)
					So(entity.Actions[0].Name, ShouldEqual, "ACTION2")
					So(len(entity.States), ShouldEqual, 0)
					So(len(entity.Scripts), ShouldEqual, 1)
					So(entity.Scripts[0].Name, ShouldEqual, "script2")
					So(len(entity.Metrics), ShouldEqual, 1)
					So(entity.Metrics[0].Name, ShouldEqual, "bar2")
					So(entity.Metrics[0].Options.Items[0].Name, ShouldEqual, "foo2")

					// v4
					entity.Actions = []*m.EntityAction{}
					entity.States = []*m.EntityState{}
					entity.Scripts = []*m.Script{
						script1,
						script2,
					}
					err = adaptors.Entity.Update(context.Background(), entity)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), entity.Id)
					So(err, ShouldBeNil)

					So(len(entity.Actions), ShouldEqual, 0)
					So(len(entity.States), ShouldEqual, 0)
					So(len(entity.Scripts), ShouldEqual, 2)
					So(len(entity.Metrics), ShouldEqual, 1)
					So(entity.Metrics[0].Name, ShouldEqual, "bar2")
					So(entity.Metrics[0].Options.Items[0].Name, ShouldEqual, "foo2")

					// v5
					entity.Actions = []*m.EntityAction{}
					entity.States = []*m.EntityState{}
					entity.Scripts = []*m.Script{}
					entity.Metrics = []*m.Metric{}

					err = adaptors.Entity.Update(context.Background(), entity)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), entity.Id)
					So(err, ShouldBeNil)

					So(len(entity.Actions), ShouldEqual, 0)
					So(len(entity.States), ShouldEqual, 0)
					So(len(entity.Scripts), ShouldEqual, 0)
					So(len(entity.Metrics), ShouldEqual, 0)

					// v5
					entity.Actions = []*m.EntityAction{}
					entity.States = []*m.EntityState{}
					entity.Scripts = []*m.Script{}
					entity.Metrics = []*m.Metric{
						{
							Name:        "bar4",
							Description: "bar4",
							Options: m.MetricOptions{
								Items: []m.MetricOptionsItem{
									{
										Name:        "foo4",
										Description: "foo4",
										Color:       "foo4",
										Translate:   "foo4",
										Label:       "foo4",
									},
								},
							},
						},
					}

					err = adaptors.Entity.Update(context.Background(), entity)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), entity.Id)
					So(err, ShouldBeNil)

					So(len(entity.Actions), ShouldEqual, 0)
					So(len(entity.States), ShouldEqual, 0)
					So(len(entity.Scripts), ShouldEqual, 0)
					So(len(entity.Metrics), ShouldEqual, 1)
					So(entity.Metrics[0].Name, ShouldEqual, "bar4")
					So(entity.Metrics[0].Options.Items[0].Name, ShouldEqual, "foo4")
				})
			})

			t.Run("Import", func(t *testing.T) {
				Convey("", t, func(ctx C) {
					// entity
					entity := &m.Entity{
						Id:         common.EntityId("sensor.entity2"),
						PluginName: "sensor",
						Scripts: []*m.Script{
							{
								Id:     456,
								Lang:   common.ScriptLangCoffee,
								Name:   "script3",
								Source: "print 'OK'",
							},
							{
								Id:     789,
								Lang:   common.ScriptLangCoffee,
								Name:   "script2",
								Source: "print 'OK'",
							},
						},
						Actions: []*m.EntityAction{
							{
								Id:          123,
								Name:        "ACTION3",
								Description: "ACTION3",
								Script: &m.Script{
									Id:     456,
									Lang:   common.ScriptLangCoffee,
									Name:   "script3",
									Source: "print 'OK'",
								},
							},
						},
						States: []*m.EntityState{
							{
								Id:          123,
								Description: "STATE3",
								Name:        "STATE3",
							},
						},
						Metrics: []*m.Metric{
							{
								Name:        "bar3",
								Description: "bar3",
								Options: m.MetricOptions{
									Items: []m.MetricOptionsItem{
										{
											Name:        "foo3",
											Description: "foo3",
											Color:       "foo3",
											Translate:   "foo3",
											Label:       "foo3",
										},
									},
								},
							},
						},
					}
					err = adaptors.Entity.Import(context.Background(), entity)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), entity.Id)
					So(err, ShouldBeNil)

					So(len(entity.Actions), ShouldEqual, 1)
					So(entity.Actions[0].Name, ShouldEqual, "ACTION3")
					So(len(entity.States), ShouldEqual, 1)
					So(entity.States[0].Name, ShouldEqual, "STATE3")
					So(len(entity.Scripts), ShouldEqual, 2)
					So(entity.Scripts[0].Name, ShouldEqual, "script2")
					So(entity.Scripts[1].Name, ShouldEqual, "script3")
					So(len(entity.Metrics), ShouldEqual, 1)
					So(entity.Metrics[0].Name, ShouldEqual, "bar3")
					So(entity.Metrics[0].Options.Items[0].Name, ShouldEqual, "foo3")

					entity3 := &m.Entity{
						Id:         common.EntityId("sensor.entity3"),
						PluginName: "sensor",
						Scripts:    []*m.Script{},
						Actions:    []*m.EntityAction{},
						States:     []*m.EntityState{},
						Metrics:    []*m.Metric{},
						Attributes: NetAttr(),
						Settings:   NetSettings(),
					}
					err = adaptors.Entity.Import(context.Background(), entity3)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), entity3.Id)
					So(err, ShouldBeNil)

					So(len(entity.Actions), ShouldEqual, 0)
					So(len(entity.States), ShouldEqual, 0)
					So(len(entity.Scripts), ShouldEqual, 0)
					So(len(entity.Metrics), ShouldEqual, 0)
					So(entity.Settings, ShouldNotBeEmpty)
					So(entity.Settings["s"].String(), ShouldEqual, "s")
					So(entity.Attributes, ShouldNotBeEmpty)
					So(entity.Attributes["s"].Name, ShouldEqual, "s")
					So(entity.Attributes["s"].String(), ShouldEqual, "")

					err = adaptors.Entity.Import(context.Background(), entity3)
					So(err, ShouldNotBeNil)
				})
			})

			t.Run("Delete", func(t *testing.T) {
				Convey("", t, func(ctx C) {
					err = adaptors.Entity.Delete(context.Background(), "sensor.entity2")
					So(err, ShouldBeNil)

					err = adaptors.Entity.Delete(context.Background(), "sensor.entity2")
					So(err, ShouldBeNil)

				})
			})

			t.Run("List", func(t *testing.T) {
				Convey("", t, func(ctx C) {
					list, total, err := adaptors.Entity.List(context.Background(), 5, 0, "desc", "id", false, nil, nil, nil)
					So(err, ShouldBeNil)
					So(total, ShouldEqual, 2)
					So(len(list), ShouldEqual, 2)
				})
			})

			t.Run("Search", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					list, total, err := adaptors.Entity.Search(context.Background(), "entity23", 5, 0)
					So(err, ShouldBeNil)
					So(total, ShouldEqual, 0)
					So(len(list), ShouldEqual, 0)

					list, total, err = adaptors.Entity.Search(context.Background(), "entity", 5, 0)
					So(err, ShouldBeNil)
					So(total, ShouldEqual, 2)
					So(len(list), ShouldEqual, 2)

					list, total, err = adaptors.Entity.Search(context.Background(), "entity1", 5, 0)
					So(err, ShouldBeNil)
					So(total, ShouldEqual, 1)
					So(len(list), ShouldEqual, 1)
				})
			})

			t.Run("GetById", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					_, err := adaptors.Entity.GetById(context.Background(), "sensor.entity23", false)
					So(err, ShouldNotBeNil)

					entity, err := adaptors.Entity.GetById(context.Background(), "sensor.entity1", false)
					So(err, ShouldBeNil)
					So(entity.Id, ShouldEqual, "sensor.entity1")
				})
			})

			t.Run("GetByIds", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					list, err := adaptors.Entity.GetByIds(context.Background(), []common.EntityId{"sensor.entity23"}, false)
					So(err, ShouldBeNil)
					So(len(list), ShouldEqual, 0)

					list, err = adaptors.Entity.GetByIds(context.Background(), []common.EntityId{"sensor.entity1"}, false)
					So(err, ShouldBeNil)
					So(list[0].Id, ShouldEqual, "sensor.entity1")
				})
			})

			t.Run("GetByType", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					list, err := adaptors.Entity.GetByType(context.Background(), "foo", 5, 0)
					So(err, ShouldBeNil)
					So(len(list), ShouldEqual, 0)

					list, err = adaptors.Entity.GetByType(context.Background(), "sensor", 5, 0)
					So(err, ShouldBeNil)
					So(len(list), ShouldEqual, 1)
				})
			})

			t.Run("UpdateAutoload", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					err := adaptors.Entity.UpdateAutoload(context.Background(), "sensor.entity1", false)
					So(err, ShouldBeNil)

					entity, err := adaptors.Entity.GetById(context.Background(), "sensor.entity1", false)
					So(err, ShouldBeNil)
					So(entity.AutoLoad, ShouldEqual, false)

					err = adaptors.Entity.UpdateAutoload(context.Background(), "sensor.entity1", true)
					So(err, ShouldBeNil)

					entity, err = adaptors.Entity.GetById(context.Background(), "sensor.entity1", false)
					So(err, ShouldBeNil)
					So(entity.AutoLoad, ShouldEqual, true)
				})
			})

		})
		So(err, ShouldBeNil)
	})
}
