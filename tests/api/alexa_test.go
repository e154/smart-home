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
	"github.com/e154/smart-home/api/server/v1/responses"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestAlexa(t *testing.T) {

	type newAlexaRequest struct {
		ResponseCode int
		Alexa        models.NewAlexaSkill
		Id           int64
	}

	apps := []newAlexaRequest{
		{
			ResponseCode: 400,
			Alexa: models.NewAlexaSkill{
				SkillId:     "",
				Description: "",
				Status:      "",
			},
		},
		{
			ResponseCode: 500,
			Alexa: models.NewAlexaSkill{
				SkillId:     "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description: "",
				Status:      "",
			},
		},
		{
			ResponseCode: 500,
			Alexa: models.NewAlexaSkill{
				SkillId:     "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description: "hello world",
				Status:      "",
			},
		},
		{
			ResponseCode: 200,
			Alexa: models.NewAlexaSkill{
				SkillId:     "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description: "hello world",
				Status:      "enabled",
			},
		},
		{
			ResponseCode: 200,
			Alexa: models.NewAlexaSkill{
				SkillId:     "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description: "hello world",
				Status:      "enabled",
			},
		},
		{
			ResponseCode: 200,
			Alexa: models.NewAlexaSkill{
				SkillId:     "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description: "hello world",
				Status:      "enabled",
			},
		},
	}

	type updateAlexaRequest struct {
		ResponseCode int
		Alexa        models.UpdateAlexaSkill
	}

	updateAlexas := []updateAlexaRequest{
		{
			ResponseCode: 400,
			Alexa:        models.UpdateAlexaSkill{},
		},
		{
			ResponseCode: 200,
			Alexa: models.UpdateAlexaSkill{
				Id:          1,
				SkillId:     "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description: "hello world",
				Status:      "enabled",
			},
		},
		{
			ResponseCode: 200,
			Alexa: models.UpdateAlexaSkill{
				Id:               1,
				SkillId:          "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description:      "hello world",
				Status:           "enabled",
				OnLaunchScriptId: common.Int64(1),
			},
		},
		{
			ResponseCode: 200,
			Alexa: models.UpdateAlexaSkill{
				Id:                   1,
				SkillId:              "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description:          "hello world",
				Status:               "enabled",
				OnLaunchScriptId:     common.Int64(1),
				OnSessionEndScriptId: common.Int64(1),
			},
		},
		{
			ResponseCode: 500,
			Alexa: models.UpdateAlexaSkill{
				Id:                   1,
				SkillId:              "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description:          "hello world",
				Status:               "enabled",
				OnLaunchScriptId:     common.Int64(1),
				OnSessionEndScriptId: common.Int64(1),
				Intents: []*models.AlexaIntent{
					{
						Name:         "IntentName",
						AlexaSkillId: 1,
					},
				},
			},
		},
		{
			ResponseCode: 200,
			Alexa: models.UpdateAlexaSkill{
				Id:                   1,
				SkillId:              "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description:          "hello world",
				Status:               "enabled",
				OnLaunchScriptId:     common.Int64(1),
				OnSessionEndScriptId: common.Int64(1),
				Intents: []*models.AlexaIntent{
					{
						Name:         "IntentName",
						AlexaSkillId: 1,
						ScriptId:     1,
						Description:  "hello world",
					},
				},
			},
		},
		{
			ResponseCode: 200,
			Alexa: models.UpdateAlexaSkill{
				Id:                   1,
				SkillId:              "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description:          "hello world",
				Status:               "enabled",
				OnLaunchScriptId:     common.Int64(1),
				OnSessionEndScriptId: common.Int64(1),
				Intents: []*models.AlexaIntent{
					{
						Name:         "IntentName",
						AlexaSkillId: 1,
						ScriptId:     1,
						Description:  "hello world",
					},
					{
						Name:         "IntentName2",
						AlexaSkillId: 1,
						ScriptId:     1,
						Description:  "hello world",
					},
				},
			},
		},
		{
			ResponseCode: 200,
			Alexa: models.UpdateAlexaSkill{
				Id:                   1,
				SkillId:              "amzn1.ask.skill.2cc6856d-e79f-412e-b311-5ca7ebfa8754",
				Description:          "hello world",
				Status:               "enabled",
				OnLaunchScriptId:     common.Int64(1),
				OnSessionEndScriptId: common.Int64(1),
				Intents: []*models.AlexaIntent{
					{
						Name:         "IntentName",
						AlexaSkillId: 1,
						ScriptId:     1,
						Description:  "hello world",
					},
				},
			},
		},
	}

	var script1 *m.Script

	Convey("added scripts", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core,
			accessList *access_list.AccessListService, ) {

			// stop core
			// ------------------------------------------------
			err := core.Stop()
			So(err, ShouldBeNil)

			// clear database
			err = migrations.Purge()
			ctx.So(err, ShouldBeNil)

			// add common script
			// ------------------------------------------------
			script1 = &m.Script{
				Lang:        "coffeescript",
				Name:        "mb_dev1_condition_check_v1",
				Source:      "",
				Description: "condition check",
			}
			ok, _ := script1.Valid()
			So(ok, ShouldEqual, true)

			engine1, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)
			err = engine1.Compile()
			So(err, ShouldBeNil)
			script1Id, err := adaptors.Script.Add(script1)
			So(err, ShouldBeNil)
			script1, err = adaptors.Script.GetById(script1Id)
			So(err, ShouldBeNil)
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("POST /alexa", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core,
			accessList *access_list.AccessListService, ) {

			// add roles
			AddRoles(adaptors, accessList, ctx)

			//
			go server.Start()

			time.Sleep(time.Second * 1)

			client := NewClient(server.GetEngine())

			// login
			err := client.LoginAsAdmin()
			ctx.So(err, ShouldBeNil)

			apps[4].Alexa.OnLaunchScriptId = common.Int64(script1.Id)
			apps[4].Alexa.OnSessionEndScriptId = common.Int64(script1.Id)

			for i, req := range apps {
				client.SetToken(accessToken)
				res := client.NewAlexa(req.Alexa)
				ctx.So(res.Code, ShouldEqual, req.ResponseCode)

				if req.ResponseCode != 200 {
					continue
				}

				app := &models.AlexaSkill{}
				err = json.Unmarshal(res.Body.Bytes(), app)
				ctx.So(err, ShouldBeNil)
				apps[i].Id = app.Id
				if req.Alexa.OnSessionEndScriptId != nil {
					ctx.So(app.OnSessionEndScriptId, ShouldNotBeNil)
					ctx.So(*app.OnSessionEndScriptId, ShouldEqual, *req.Alexa.OnSessionEndScriptId)
				}
				if req.Alexa.OnLaunchScriptId != nil {
					ctx.So(app.OnLaunchScriptId, ShouldNotBeNil)
					ctx.So(*app.OnLaunchScriptId, ShouldEqual, *req.Alexa.OnLaunchScriptId)
				}

				ctx.So(app.Status, ShouldEqual, req.Alexa.Status)
				ctx.So(app.Description, ShouldEqual, req.Alexa.Description)
			}
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /alexa/{id}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			client := NewClient(server.GetEngine())

			// positive
			client.SetToken(accessToken)
			res := client.GetAlexa(404)
			ctx.So(res.Code, ShouldEqual, 404)

			for i, req := range apps {
				if req.ResponseCode != 200 {
					continue
				}

				res := client.GetAlexa(req.Id)
				ctx.So(res.Code, ShouldEqual, 200)

				if req.ResponseCode != 200 {
					continue
				}

				app := &models.AlexaSkill{}
				err := json.Unmarshal(res.Body.Bytes(), app)
				ctx.So(err, ShouldBeNil)
				apps[i].Id = app.Id
				if req.Alexa.OnSessionEndScriptId != nil {
					ctx.So(app.OnSessionEndScriptId, ShouldNotBeNil)
					ctx.So(*app.OnSessionEndScriptId, ShouldEqual, *req.Alexa.OnSessionEndScriptId)
				}
				if req.Alexa.OnLaunchScriptId != nil {
					ctx.So(app.OnLaunchScriptId, ShouldNotBeNil)
					ctx.So(*app.OnLaunchScriptId, ShouldEqual, *req.Alexa.OnLaunchScriptId)
				}

				ctx.So(app.Status, ShouldEqual, req.Alexa.Status)
				ctx.So(app.Description, ShouldEqual, req.Alexa.Description)
			}
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("PUT /alexa/{id}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			client := NewClient(server.GetEngine())

			// positive
			client.SetToken(accessToken)
			res := client.GetAlexa(404)
			ctx.So(res.Code, ShouldEqual, 404)

			for _, req := range updateAlexas {

				res := client.UpdateAlexa(1, req.Alexa)
				ctx.So(res.Code, ShouldEqual, req.ResponseCode)

				if req.ResponseCode != 200 {
					continue
				}

				res = client.GetAlexa(1)
				ctx.So(res.Code, ShouldEqual, 200)

				app := &models.AlexaSkill{}
				err := json.Unmarshal(res.Body.Bytes(), app)
				ctx.So(err, ShouldBeNil)

				if req.Alexa.OnSessionEndScriptId != nil {
					ctx.So(app.OnSessionEndScriptId, ShouldNotBeNil)
					ctx.So(*app.OnSessionEndScriptId, ShouldEqual, *req.Alexa.OnSessionEndScriptId)
				}
				if req.Alexa.OnLaunchScriptId != nil {
					ctx.So(app.OnLaunchScriptId, ShouldNotBeNil)
					ctx.So(*app.OnLaunchScriptId, ShouldEqual, *req.Alexa.OnLaunchScriptId)
				}
				if len(req.Alexa.Intents) > 0 {
					ctx.So(len(app.Intents), ShouldEqual, len(req.Alexa.Intents))
				}

				ctx.So(app.Status, ShouldEqual, req.Alexa.Status)
				ctx.So(app.Description, ShouldEqual, req.Alexa.Description)
			}
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("DELETE /alexa/{id}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			client := NewClient(server.GetEngine())

			// negative
			client.SetToken(invalidToken1)
			res := client.DeleteAlexa(404)
			ctx.So(res.Code, ShouldEqual, 401)
			client.SetToken(invalidToken2)
			res = client.DeleteAlexa(404)
			ctx.So(res.Code, ShouldEqual, 403)

			// positive
			client.SetToken(accessToken)
			res = client.DeleteAlexa(404)
			ctx.So(res.Code, ShouldEqual, 404)

			res = client.DeleteAlexa(1)
			ctx.So(res.Code, ShouldEqual, 200)

			res = client.DeleteAlexa(1)
			ctx.So(res.Code, ShouldEqual, 404)
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /alexas", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			client := NewClient(server.GetEngine())

			// negative
			client.SetToken(invalidToken1)
			res := client.GetAlexaList(5, 0, "DESC", "id")
			ctx.So(res.Code, ShouldEqual, 401)
			client.SetToken(invalidToken2)
			res = client.GetAlexaList(5, 0, "DESC", "id")
			ctx.So(res.Code, ShouldEqual, 403)

			// positive
			client.SetToken(accessToken)

			listGetter := func(limit, offset, realLimit, realOffset int) {
				res = client.GetAlexaList(limit, offset, "DESC", "id")
				ctx.So(res.Code, ShouldEqual, 200)

				appList := responses.AlexaSkillList{}
				err := json.Unmarshal(res.Body.Bytes(), &appList.Body)
				ctx.So(err, ShouldBeNil)

				ctx.So(appList.Body.Meta.ObjectCount, ShouldEqual, 2)
				ctx.So(len(appList.Body.Items), ShouldEqual, realLimit)
				ctx.So(appList.Body.Meta.Limit, ShouldEqual, limit)
				ctx.So(appList.Body.Meta.Offset, ShouldEqual, realOffset)
			}

			listGetter(1, 0, 1, 0)
			listGetter(1, 3, 0, 3)
			listGetter(5, 0, 2, 0)

		})
		if err != nil {
			panic(err.Error())
		}
	})

}
