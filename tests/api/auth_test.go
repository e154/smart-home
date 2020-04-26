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

func TestAuth(t *testing.T) {

	Convey("POST /signin", t, func(ctx C) {
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

			type signinReqParams struct {
				Login    string
				Pass     string
				RespCode int
			}

			client := NewClient(server.GetEngine())

			reqParams := []signinReqParams{
				{"guest@e154.ru", "guest", 401},
				{"admin@e154.ru", "admin", 200},
				{"admin@e154.ru", "admin1", 403},
				{"admin1@e154.ru", "admin", 401},
				{"user@e154.ru", "user", 200},
				{"user1@e154.ru", "user", 401},
				{"user@e154.ru", "user1", 403},
				{"demo@e154.ru", "demo", 200},
			}

			for _, req := range reqParams {
				//fmt.Println(req.Login, req.Pass)
				client.BasicAuth(req.Login, req.Pass)
				res := client.Signin()
				ctx.So(res.Code, ShouldEqual, req.RespCode)
			}

			client.BasicAuth("admin@e154.ru", "admin")
			res := client.Signin()
			ctx.So(res.Code, ShouldEqual, 200)

			currentUser := &models.AuthSignInResponse{}
			err = json.Unmarshal(res.Body.Bytes(), currentUser)
			ctx.So(err, ShouldBeNil)

			ctx.So(currentUser.CurrentUser.Id, ShouldEqual, 1)
			ctx.So(currentUser.CurrentUser.Nickname, ShouldEqual, "admin")
			ctx.So(currentUser.CurrentUser.FirstName, ShouldEqual, "")
			ctx.So(currentUser.CurrentUser.Email, ShouldEqual, "admin@e154.ru")
			ctx.So(currentUser.CurrentUser.LastName, ShouldEqual, "")
			ctx.So(len(currentUser.CurrentUser.History), ShouldEqual, 1)
			ctx.So(currentUser.CurrentUser.Role, ShouldNotBeNil)
			ctx.So(currentUser.CurrentUser.Role.Name, ShouldEqual, "admin")
			ctx.So(currentUser.AccessToken, ShouldNotBeNil)

			accessToken = currentUser.AccessToken
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /access_list", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			client := NewClient(server.GetEngine())

			client.SetToken(invalidToken1)
			res := client.GetAccessList()
			ctx.So(res.Code, ShouldEqual, 401)

			client.SetToken(accessToken)
			res = client.GetAccessList()
			ctx.So(res.Code, ShouldEqual, 200)

			type AccessList struct {
				AccessList models.AccessList `json:"access_list"`
			}

			accessList := &AccessList{}
			err := json.Unmarshal(res.Body.Bytes(), &accessList)
			ctx.So(err, ShouldBeNil)

			ctx.So(len(accessList.AccessList), ShouldEqual, 23)

			countrer := 0
			for item, _ := range accessList.AccessList {
				switch item {
				case "dashboard",
					"device",
					"flow",
					"device_action",
					"device_state",
					"gate",
					"log",
					"script",
					"template",
					"user",
					"ws",
					"image",
					"map",
					"map_zone",
					"mqtt",
					"node",
					"notifr",
					"scenarios",
					"worker",
					"workflow",
					"alexa",
					"zigbee2mqtt",
					"metric":
					countrer++
				default:
					countrer--
				}
			}
			ctx.So(countrer, ShouldEqual, 23)
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("POST /recovery", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("POST /reset", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("POST /signout", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			client := NewClient(server.GetEngine())
			client.SetToken(accessToken)
			res := client.Signout()
			ctx.So(res.Code, ShouldEqual, 200)

			err := core.Stop()
			So(err, ShouldBeNil)

			server.Shutdown()
		})
		if err != nil {
			panic(err.Error())
		}
	})
}
