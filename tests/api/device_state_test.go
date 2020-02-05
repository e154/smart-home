package api

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDeviceState(t *testing.T) {

	Convey("POST /device_state", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core,
			accessList *access_list.AccessListService, ) {


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
