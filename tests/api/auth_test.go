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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {

	var accessToken string

	Convey("/signin", t, func(ctx C) {
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
				res := client.Signin(req.Login, req.Pass)
				ctx.So(res.Code, ShouldEqual, req.RespCode)
			}

			res := client.Signin("admin@e154.ru", "admin")
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

	Convey("/access_list", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			// request params
			req := func(token string) (w *httptest.ResponseRecorder) {
				request, _ := http.NewRequest("GET", "/api/v1/access_list", nil)
				request.Header.Add("accept", "application/json")
				request.Header.Set("Authorization", token)
				w = httptest.NewRecorder()
				server.GetEngine().ServeHTTP(w, request)
				return
			}

			res := req("qweqweasd1")
			ctx.So(res.Code, ShouldEqual, 401)

			res = req(accessToken)
			ctx.So(res.Code, ShouldEqual, 200)

			type AccessList struct {
				AccessList models.AccessList `json:"access_list"`
			}

			accessList := &AccessList{}
			err := json.Unmarshal(res.Body.Bytes(), &accessList)
			ctx.So(err, ShouldBeNil)

			ctx.So(len(accessList.AccessList), ShouldEqual, 20)

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
					"workflow":
					countrer++
				default:
					countrer--
				}
			}
			ctx.So(countrer, ShouldEqual, 20)
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("/recovery", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("/reset", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("/signout", t, func(ctx C) {
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
		})
		if err != nil {
			panic(err.Error())
		}
	})
}
