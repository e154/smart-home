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

package scripts

import (
	"fmt"
	"testing"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

func Test3(t *testing.T) {

	var state string
	store = func(i interface{}) {
		state = fmt.Sprintf("%v", i)
	}

	Convey("eval script", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService) {

			storeRegisterCallback(scriptService)

			// engine
			// ------------------------------------------------
			script := &m.Script{
				Lang: common.ScriptLangJavascript,
			}
			engine, err := scriptService.NewEngine(script)
			So(err, ShouldBeNil)

			// scripts
			// ------------------------------------------------
			scripts := GetScripts(ctx, scriptService, adaptors, 3)

			// execute script
			// ------------------------------------------------
			_, err = engine.EvalString(scripts["script3"].Compiled)
			So(err, ShouldBeNil)

			_, err = engine.AssertFunction("on_enter")
			So(err, ShouldBeNil)
			So(state, ShouldEqual, "on_enter")

			_, err = engine.AssertFunction("on_exit")
			So(err, ShouldBeNil)
			So(state, ShouldEqual, "on_exit")
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
