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

package api

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestDeviceState(t *testing.T) {

	type newDeviceStateRequest struct {
		ResponseCode int
		State        models.NewDeviceState
		Id           int64
	}

	var statuses []newDeviceStateRequest
	var client *Client
	var positiveDeviceStateId int64

	Convey("POST /device_state", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core,
			accessList *access_list.AccessListService, ) {

			// clear database
			err := migrations.Purge()
			ctx.So(err, ShouldBeNil)

			// add roles
			AddRoles(adaptors, accessList, ctx)

			go server.Start()

			time.Sleep(time.Second * 1)

			client = NewClient(server.GetEngine())

			deviceId, err := AddDevice(adaptors, nil)
			ctx.So(err, ShouldBeNil)
			ctx.So(deviceId, ShouldNotBeEmpty)

			state1 := models.NewDeviceState{
				Description: "",
				SystemName:  "",
				Device: &models.DeviceStateDevice{
					Id: deviceId,
				},
			}

			// negative
			client.SetToken("qqweqwe")
			res := client.NewDeviceState(state1)
			ctx.So(res.Code, ShouldEqual, 401)

			statuses = []newDeviceStateRequest{
				{
					ResponseCode: 400,
					State: models.NewDeviceState{
						Description: "",
						SystemName:  "",
						Device: &models.DeviceStateDevice{
							Id: 0,
						},
					},
				},
				{
					ResponseCode: 400,
					State: models.NewDeviceState{
						Description: "",
						SystemName:  "",
						Device: &models.DeviceStateDevice{
							Id: deviceId,
						},
					},
				},
				{
					ResponseCode: 200,
					State: models.NewDeviceState{
						Description: "",
						SystemName:  "state1",
						Device: &models.DeviceStateDevice{
							Id: deviceId,
						},
					},
				},
				{
					ResponseCode: 500,
					State: models.NewDeviceState{
						Description: "",
						SystemName:  "state1",
						Device: &models.DeviceStateDevice{
							Id: deviceId,
						},
					},
				},
				{
					ResponseCode: 200,
					State: models.NewDeviceState{
						Description: "description",
						SystemName:  "state2",
						Device: &models.DeviceStateDevice{
							Id: deviceId,
						},
					},
				},
			}

			// login
			err = client.LoginAsAdmin()
			ctx.So(err, ShouldBeNil)

			for i, deviceStatus := range statuses {

				//fmt.Println("---")
				//debug.Println(deviceStatus.State)

				// positive
				res = client.NewDeviceState(deviceStatus.State)
				ctx.So(res.Code, ShouldEqual, deviceStatus.ResponseCode)

				if deviceStatus.ResponseCode != 200 {
					continue
				}

				device := &models.DeviceState{}
				err = json.Unmarshal(res.Body.Bytes(), device)
				ctx.So(err, ShouldBeNil)
				statuses[i].Id = device.Id
				positiveDeviceStateId = device.Id
			}
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /device_state/{id}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			// negative
			client.SetToken("qqweqwe")

			res := client.GetDeviceState(positiveDeviceStateId)
			ctx.So(res.Code, ShouldEqual, 401)

			// login
			err := client.LoginAsAdmin()
			ctx.So(err, ShouldBeNil)

			// positive
			res = client.GetDeviceState(positiveDeviceStateId)
			ctx.So(res.Code, ShouldEqual, 200)

			for _, deviceStatus := range statuses {

				if deviceStatus.ResponseCode != 200 {
					continue
				}

				res = client.GetDeviceState(deviceStatus.Id)
				ctx.So(res.Code, ShouldEqual, 200)

				state := &models.DeviceState{}
				err = json.Unmarshal(res.Body.Bytes(), state)
				ctx.So(err, ShouldBeNil)

				ctx.So(state.Id, ShouldEqual, deviceStatus.Id)
				//ctx.So(state.Device, ShouldNotBeNil)
				//ctx.So(state.Device.Id, ShouldEqual, deviceStatus.State.Device.Id)
				ctx.So(state.Description, ShouldEqual, deviceStatus.State.Description)
				ctx.So(state.SystemName, ShouldEqual, deviceStatus.State.SystemName)
			}
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("PUT /device_state/{id}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("DELETE /device_state/{id}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /device_states", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})
}
