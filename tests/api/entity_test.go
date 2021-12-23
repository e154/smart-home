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
	"github.com/e154/smart-home/common/debug"
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
			//ctx.So(area.Id, ShouldBeZeroValue)

			// script
			script, err := AddScript("script1", sourceScript, adaptors, scriptService)
			ctx.So(err, ShouldBeNil)

			t.Run("success", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					request := &gw.NewEntityRequest{
						Name:        "light",
						PluginName:  "sensor",
						Description: "light toggle",
						Area: &gw.NewEntityRequest_Area{
							Id: area.Id,
						},
						Hidden:   false,
						AutoLoad: true,
						Actions: []*gw.NewEntityRequest_Action{
							{
								Name:        "ON",
								Description: "toggle on",
								Script: &gw.NewEntityRequest_Action_Script{
									Id: script.Id,
								},
							},
							{
								Name:        "OFF",
								Description: "toggle off",
								Script: &gw.NewEntityRequest_Action_Script{
									Id: script.Id,
								},
							},
							{
								Name:        "CHECK",
								Description: "check status",
								Script: &gw.NewEntityRequest_Action_Script{
									Id: script.Id,
								},
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
						},
						Settings: nil,
					}

					entity, err := client.AddEntity(c, request)
					ctx.So(err, ShouldBeNil)
					debug.Println(entity)
					//ctx.So(bridge, ShouldNotBeNil)

				})
			})

		})

		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
