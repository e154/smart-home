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

package api

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/controllers"
	gw "github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	container2 "github.com/e154/smart-home/tests/api/container"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEntity(t *testing.T) {

	const sourceScript = `
`
	var createRequest = &gw.NewEntityRequest{
		Name:        "light",
		PluginName:  "sensor",
		Description: "light toggle",
		Area:        &gw.NewEntityRequest_Area{},
		AutoLoad:    true,
		Actions: []*gw.NewEntityRequest_Action{
			{
				Name:        "ON",
				Description: "toggle on",
				Script:      &gw.NewEntityRequest_Action_Script{},
			},
			{
				Name:        "OFF",
				Description: "toggle off",
				Script:      &gw.NewEntityRequest_Action_Script{},
			},
			{
				Name:        "CHECK",
				Description: "check status",
				Script:      &gw.NewEntityRequest_Action_Script{},
			},
		},
		States: []*gw.NewEntityRequest_State{
			{
				Name:        "OK",
				Description: "status ok",
			},
			{
				Name:        "ERROR",
				Description: "status error",
			},
		},
		Attributes: map[string]*gw.Attribute{
			"value1": {
				Name:    "value",
				Type:    gw.Types_STRING,
				String_: common.String("some text"),
			},
			"value2": {
				Name: "value2",
				Type: gw.Types_INT,
				Int:  common.Int64(123),
			},
			"value3": {
				Name: "value3",
				Type: gw.Types_BOOL,
				Bool: common.Bool(true),
			},
			"value4": {
				Name:  "value4",
				Type:  gw.Types_FLOAT,
				Float: common.Float32(0.123),
			},
		},
		Settings: map[string]*gw.Attribute{
			"value1": {
				Name:    "value",
				Type:    gw.Types_STRING,
				String_: common.String("some text"),
			},
			"value2": {
				Name: "value2",
				Type: gw.Types_INT,
				Int:  common.Int64(123),
			},
			"value3": {
				Name: "value3",
				Type: gw.Types_BOOL,
				Bool: common.Bool(true),
			},
			"value4": {
				Name:  "value4",
				Type:  gw.Types_FLOAT,
				Float: common.Float32(0.123),
			},
		},
	}

	var updateRequest = &gw.UpdateEntityRequest{
		Id:          "sensor.light",
		Name:        "light",
		PluginName:  "sensor",
		Description: "light toggle FX",
		Area:        &gw.UpdateEntityRequest_Area{},
		AutoLoad:    false,
		Actions: []*gw.UpdateEntityRequest_Action{
			{
				Name:        "ON FX",
				Description: "toggle on FX",
				Script:      &gw.UpdateEntityRequest_Action_Script{},
			},
			{
				Name:        "OFF FX",
				Description: "toggle off FX",
				Script:      &gw.UpdateEntityRequest_Action_Script{},
			},
			{
				Name:        "CHECK FX",
				Description: "check status FX",
				Script:      &gw.UpdateEntityRequest_Action_Script{},
			},
			{
				Name:        "FX",
				Description: "status FX",
				Script:      &gw.UpdateEntityRequest_Action_Script{},
			},
		},
		States: []*gw.UpdateEntityRequest_State{
			{
				Name:        "FX",
				Description: "status FX",
			},
		},
		Attributes: map[string]*gw.Attribute{
			"value1": {
				Name:    "value",
				Type:    gw.Types_STRING,
				String_: common.String("some text  FX"),
			},
			"value2": {
				Name: "value2",
				Type: gw.Types_INT,
				Int:  common.Int64(456),
			},
			"value3": {
				Name: "value3",
				Type: gw.Types_BOOL,
				Bool: common.Bool(false),
			},
			"value4  FX": {
				Name:  "value4  FX",
				Type:  gw.Types_FLOAT,
				Float: common.Float32(0.456),
			},
		},
		Settings: map[string]*gw.Attribute{
			"value1": {
				Name:    "value",
				Type:    gw.Types_STRING,
				String_: common.String("some text FX"),
			},
			"value2": {
				Name: "value2",
				Type: gw.Types_INT,
				Int:  common.Int64(456),
			},
			"value3": {
				Name: "value3",
				Type: gw.Types_BOOL,
				Bool: common.Bool(false),
			},
			"value4 FX": {
				Name:  "value4 FX",
				Type:  gw.Types_FLOAT,
				Float: common.Float32(0.456),
			},
		},
	}

	Convey("entity", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager,
			controllers *controllers.Controllers,
			dialer *container2.Dialer) {

			eventBus.Purge()
			scriptService.Purge()

			err := migrations.Purge()
			ctx.So(err, ShouldBeNil)

			c := context.Background()
			conn, err := grpc.DialContext(c, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer.Call()))
			ctx.So(err, ShouldBeNil)
			defer conn.Close()

			client := gw.NewEntityServiceClient(conn)

			// plugins
			err = AddPlugin(adaptors, "sensor")
			ctx.So(err, ShouldBeNil)

			// area
			area, err := AddArea(adaptors, "zone51")
			ctx.So(err, ShouldBeNil)
			ctx.So(area, ShouldNotBeNil)
			area2, err := AddArea(adaptors, "zone52")
			ctx.So(err, ShouldBeNil)
			ctx.So(area2, ShouldNotBeNil)
			//ctx.So(area.Id, ShouldBeZeroValue)

			// script
			script, err := AddScript("script1", sourceScript, adaptors, scriptService)
			ctx.So(err, ShouldBeNil)
			script2, err := AddScript("script2", sourceScript, adaptors, scriptService)
			ctx.So(err, ShouldBeNil)

			createRequest.Area.Id = area.Id
			for i, _ := range createRequest.Actions {
				createRequest.Actions[i].Script.Id = script.Id
			}
			updateRequest.Area.Id = area2.Id
			for i, _ := range updateRequest.Actions {
				updateRequest.Actions[i].Script.Id = script2.Id
			}

			t.Run("create", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					entity, err := client.AddEntity(c, createRequest)
					ctx.So(err, ShouldBeNil)
					//debug.Println(entity)

					ctx.So(entity.Id, ShouldEqual, "sensor.light")
					ctx.So(entity.PluginName, ShouldEqual, "sensor")
					ctx.So(entity.Description, ShouldEqual, "light toggle")
					ctx.So(entity.Area, ShouldNotBeNil)
					ctx.So(entity.Area.Id, ShouldEqual, area.Id)
					ctx.So(entity.Area.Name, ShouldEqual, area.Name)
					ctx.So(entity.Area.Description, ShouldEqual, area.Description)
					ctx.So(entity.AutoLoad, ShouldEqual, createRequest.AutoLoad)
					ctx.So(len(entity.Actions), ShouldEqual, 3)
					ctx.So(len(entity.States), ShouldEqual, 2)
					ctx.So(entity.CreatedAt, ShouldNotBeNil)
					ctx.So(entity.UpdatedAt, ShouldNotBeNil)

					// actions
					for _, action := range entity.Actions {
						switch action.Name {
						case "ON":
							ctx.So(action.Description, ShouldEqual, "toggle on")
						case "OFF":
							ctx.So(action.Description, ShouldEqual, "toggle off")
						case "CHECK":
							ctx.So(action.Description, ShouldEqual, "check status")
						}
						ctx.So(action.Script, ShouldNotBeNil)
						ctx.So(action.Script.Id, ShouldEqual, script.Id)
						ctx.So(action.Script.Lang, ShouldEqual, script.Lang)
						ctx.So(action.Script.Name, ShouldEqual, script.Name)
						ctx.So(action.Script.Description, ShouldEqual, script.Description)
						ctx.So(action.Script.CreatedAt, ShouldNotBeNil)
						ctx.So(action.Script.UpdatedAt, ShouldNotBeNil)
					}

					// states
					for _, state := range entity.States {
						switch state.Name {
						case "OK":
							ctx.So(state.Description, ShouldEqual, "status ok")
						case "ERROR":
							ctx.So(state.Description, ShouldEqual, "status error")
						}
					}

					// attributes
					for key, item := range entity.Attributes {
						switch key {
						case "value1":
							ctx.So(item.Name, ShouldEqual, "value")
							ctx.So(item.Type, ShouldEqual, gw.Types_STRING)
							ctx.So(*item.String_, ShouldEqual, "")
						case "value2":
							ctx.So(item.Name, ShouldEqual, "value2")
							ctx.So(item.Type, ShouldEqual, gw.Types_INT)
							ctx.So(*item.Int, ShouldEqual, 0)
						case "value3":
							ctx.So(item.Name, ShouldEqual, "value3")
							ctx.So(item.Type, ShouldEqual, gw.Types_BOOL)
							ctx.So(*item.Bool, ShouldEqual, false)
						case "value4":
							ctx.So(item.Name, ShouldEqual, "value4")
							ctx.So(item.Type, ShouldEqual, gw.Types_FLOAT)
							ctx.So(*item.Float, ShouldEqual, 0)
						}
					}

					// settings
					for key, item := range entity.Settings {
						switch key {
						case "value1":
							ctx.So(item.Name, ShouldEqual, "value")
							ctx.So(item.Type, ShouldEqual, gw.Types_STRING)
							ctx.So(*item.String_, ShouldEqual, "some text")
						case "value2":
							ctx.So(item.Name, ShouldEqual, "value2")
							ctx.So(item.Type, ShouldEqual, gw.Types_INT)
							ctx.So(*item.Int, ShouldEqual, 123)
						case "value3":
							ctx.So(item.Name, ShouldEqual, "value3")
							ctx.So(item.Type, ShouldEqual, gw.Types_BOOL)
							ctx.So(*item.Bool, ShouldEqual, true)
						case "value4":
							ctx.So(item.Name, ShouldEqual, "value4")
							ctx.So(item.Type, ShouldEqual, gw.Types_FLOAT)
							ctx.So(*item.Float, ShouldEqual, 0.123)
						}
					}

					t.Run("update", func(t *testing.T) {
						Convey("", t, func(ctx C) {
							entity, err := client.UpdateEntity(c, updateRequest)
							ctx.So(err, ShouldBeNil)
							//debug.Println(entity)

							ctx.So(entity.Id, ShouldEqual, "sensor.light")
							ctx.So(entity.PluginName, ShouldEqual, "sensor")
							ctx.So(entity.Description, ShouldEqual, "light toggle FX")
							ctx.So(entity.Area, ShouldNotBeNil)
							ctx.So(entity.Area.Id, ShouldEqual, area2.Id)
							ctx.So(entity.Area.Name, ShouldEqual, area2.Name)
							ctx.So(entity.Area.Description, ShouldEqual, area2.Description)
							ctx.So(entity.AutoLoad, ShouldEqual, updateRequest.AutoLoad)
							ctx.So(len(entity.Actions), ShouldEqual, 4)
							ctx.So(len(entity.States), ShouldEqual, 1)
							ctx.So(entity.CreatedAt, ShouldNotBeNil)
							ctx.So(entity.UpdatedAt, ShouldNotBeNil)

							// actions
							for _, action := range entity.Actions {
								switch action.Name {
								case "ON FX":
									ctx.So(action.Description, ShouldEqual, "toggle on FX")
								case "OFF FX":
									ctx.So(action.Description, ShouldEqual, "toggle off FX")
								case "CHECK FX":
									ctx.So(action.Description, ShouldEqual, "check status FX")
								case "FX":
									ctx.So(action.Description, ShouldEqual, "status FX")
								}
								ctx.So(action.Script, ShouldNotBeNil)
								ctx.So(action.Script.Id, ShouldEqual, script2.Id)
								ctx.So(action.Script.Lang, ShouldEqual, script2.Lang)
								ctx.So(action.Script.Name, ShouldEqual, script2.Name)
								ctx.So(action.Script.Description, ShouldEqual, script2.Description)
								ctx.So(action.Script.CreatedAt, ShouldNotBeNil)
								ctx.So(action.Script.UpdatedAt, ShouldNotBeNil)
							}

							// states
							for _, state := range entity.States {
								switch state.Name {
								case "FX":
									ctx.So(state.Description, ShouldEqual, "status FX")
								}
							}

							// attributes
							for key, item := range entity.Attributes {
								switch key {
								case "value1":
									ctx.So(item.Name, ShouldEqual, "value")
									ctx.So(item.Type, ShouldEqual, gw.Types_STRING)
									ctx.So(*item.String_, ShouldEqual, "")
								case "value2":
									ctx.So(item.Name, ShouldEqual, "value2")
									ctx.So(item.Type, ShouldEqual, gw.Types_INT)
									ctx.So(*item.Int, ShouldEqual, 0)
								case "value3":
									ctx.So(item.Name, ShouldEqual, "value3")
									ctx.So(item.Type, ShouldEqual, gw.Types_BOOL)
									ctx.So(*item.Bool, ShouldEqual, false)
								case "value4  FX":
									ctx.So(item.Name, ShouldEqual, "value4  FX")
									ctx.So(item.Type, ShouldEqual, gw.Types_FLOAT)
									ctx.So(*item.Float, ShouldEqual, 0)
								}
							}

							// settings
							for key, item := range entity.Settings {
								switch key {
								case "value1":
									ctx.So(item.Name, ShouldEqual, "value")
									ctx.So(item.Type, ShouldEqual, gw.Types_STRING)
									ctx.So(*item.String_, ShouldEqual, "some text FX")
								case "value2":
									ctx.So(item.Name, ShouldEqual, "value2")
									ctx.So(item.Type, ShouldEqual, gw.Types_INT)
									ctx.So(*item.Int, ShouldEqual, 456)
								case "value3":
									ctx.So(item.Name, ShouldEqual, "value3")
									ctx.So(item.Type, ShouldEqual, gw.Types_BOOL)
									ctx.So(*item.Bool, ShouldEqual, false)
								case "value4  FX":
									ctx.So(item.Name, ShouldEqual, "value4  FX")
									ctx.So(item.Type, ShouldEqual, gw.Types_FLOAT)
									ctx.So(*item.Float, ShouldEqual, 0.456)
								}
							}

						})
					})

					t.Run("update failed", func(t *testing.T) {
						Convey("", t, func(ctx C) {
							updateRequest.Id = "qwe"
							_, err := client.UpdateEntity(c, updateRequest)
							ctx.So(err, ShouldNotBeNil)
							//debug.Println(entity)

						})
					})

					t.Run("list", func(t *testing.T) {
						Convey("", t, func(ctx C) {
							createRequest.Name += "2"
							entity, err := client.AddEntity(c, createRequest)
							ctx.So(err, ShouldBeNil)
							//debug.Println(entity)
							ctx.So(entity.Id, ShouldEqual, "sensor.light2")

							// list
							listRequest := &gw.GetEntityListRequest{}
							result, err := client.GetEntityList(c, listRequest)
							ctx.So(err, ShouldBeNil)

							//debug.Println(result)

							ctx.So(len(result.Items), ShouldEqual, 2)
							ctx.So(result.Meta.ObjectsCount, ShouldEqual, 2)
							ctx.So(result.Meta.Limit, ShouldEqual, 200)
							ctx.So(result.Meta.Offset, ShouldEqual, 0)
						})
					})

					t.Run("get by id", func(t *testing.T) {
						Convey("", t, func(ctx C) {

							entity, err := client.GetEntity(c, &gw.GetEntityRequest{Id: "sensor.light234"})
							ctx.So(err, ShouldNotBeNil)

							entity, err = client.GetEntity(c, &gw.GetEntityRequest{Id: "sensor.light2"})
							ctx.So(err, ShouldBeNil)
							//debug.Println(entity)
							ctx.So(entity.Id, ShouldEqual, "sensor.light2")
							ctx.So(entity.PluginName, ShouldEqual, "sensor")
							ctx.So(entity.Description, ShouldEqual, "light toggle")
							ctx.So(entity.Area, ShouldNotBeNil)
							ctx.So(entity.Area.Id, ShouldEqual, area.Id)
							ctx.So(entity.Area.Name, ShouldEqual, area.Name)
							ctx.So(entity.Area.Description, ShouldEqual, area.Description)
							ctx.So(entity.AutoLoad, ShouldEqual, createRequest.AutoLoad)
							ctx.So(len(entity.Actions), ShouldEqual, 3)
							ctx.So(len(entity.States), ShouldEqual, 2)
							ctx.So(entity.CreatedAt, ShouldNotBeNil)
							ctx.So(entity.UpdatedAt, ShouldNotBeNil)

							// actions
							for _, action := range entity.Actions {
								switch action.Name {
								case "ON":
									ctx.So(action.Description, ShouldEqual, "toggle on")
								case "OFF":
									ctx.So(action.Description, ShouldEqual, "toggle off")
								case "CHECK":
									ctx.So(action.Description, ShouldEqual, "check status")
								}
								ctx.So(action.Script, ShouldNotBeNil)
								ctx.So(action.Script.Id, ShouldEqual, script.Id)
								ctx.So(action.Script.Lang, ShouldEqual, script.Lang)
								ctx.So(action.Script.Name, ShouldEqual, script.Name)
								ctx.So(action.Script.Description, ShouldEqual, script.Description)
								ctx.So(action.Script.CreatedAt, ShouldNotBeNil)
								ctx.So(action.Script.UpdatedAt, ShouldNotBeNil)
							}

							// states
							for _, state := range entity.States {
								switch state.Name {
								case "OK":
									ctx.So(state.Description, ShouldEqual, "status ok")
								case "ERROR":
									ctx.So(state.Description, ShouldEqual, "status error")
								}
							}

							// attributes
							for key, item := range entity.Attributes {
								switch key {
								case "value1":
									ctx.So(item.Name, ShouldEqual, "value")
									ctx.So(item.Type, ShouldEqual, gw.Types_STRING)
									ctx.So(*item.String_, ShouldEqual, "")
								case "value2":
									ctx.So(item.Name, ShouldEqual, "value2")
									ctx.So(item.Type, ShouldEqual, gw.Types_INT)
									ctx.So(*item.Int, ShouldEqual, 0)
								case "value3":
									ctx.So(item.Name, ShouldEqual, "value3")
									ctx.So(item.Type, ShouldEqual, gw.Types_BOOL)
									ctx.So(*item.Bool, ShouldEqual, false)
								case "value4":
									ctx.So(item.Name, ShouldEqual, "value4")
									ctx.So(item.Type, ShouldEqual, gw.Types_FLOAT)
									ctx.So(*item.Float, ShouldEqual, 0)
								}
							}

							// settings
							for key, item := range entity.Settings {
								switch key {
								case "value1":
									ctx.So(item.Name, ShouldEqual, "value")
									ctx.So(item.Type, ShouldEqual, gw.Types_STRING)
									ctx.So(*item.String_, ShouldEqual, "some text")
								case "value2":
									ctx.So(item.Name, ShouldEqual, "value2")
									ctx.So(item.Type, ShouldEqual, gw.Types_INT)
									ctx.So(*item.Int, ShouldEqual, 123)
								case "value3":
									ctx.So(item.Name, ShouldEqual, "value3")
									ctx.So(item.Type, ShouldEqual, gw.Types_BOOL)
									ctx.So(*item.Bool, ShouldEqual, true)
								case "value4":
									ctx.So(item.Name, ShouldEqual, "value4")
									ctx.So(item.Type, ShouldEqual, gw.Types_FLOAT)
									ctx.So(*item.Float, ShouldEqual, 0.123)
								}
							}
						})
					})

					t.Run("search", func(t *testing.T) {
						Convey("", t, func(ctx C) {
							searchRequest := &gw.SearchEntityRequest{
								Query:  "light2",
								Limit:  10,
								Offset: 0,
							}
							result, err := client.SearchEntity(c, searchRequest)
							ctx.So(err, ShouldBeNil)
							ctx.So(len(result.Items), ShouldEqual, 1)

							searchRequest.Query = "light"
							result, err = client.SearchEntity(c, searchRequest)
							ctx.So(err, ShouldBeNil)
							ctx.So(len(result.Items), ShouldEqual, 2)
							//debug.Println(result)
						})
					})

					t.Run("delete", func(t *testing.T) {
						Convey("", t, func(ctx C) {
							// delete
							deleteRequest := &gw.DeleteEntityRequest{
								Id: entity.Id,
							}
							_, err = client.DeleteEntity(c, deleteRequest)
							ctx.So(err, ShouldBeNil)

							// list
							listRequest := &gw.GetEntityListRequest{}
							result, err := client.GetEntityList(c, listRequest)
							ctx.So(err, ShouldBeNil)
							ctx.So(len(result.Items), ShouldEqual, 1)
							ctx.So(result.Meta.ObjectsCount, ShouldEqual, 1)
							ctx.So(result.Meta.Limit, ShouldEqual, 200)
							ctx.So(result.Meta.Offset, ShouldEqual, 0)

						})
					})
				})
			})
		})

		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
