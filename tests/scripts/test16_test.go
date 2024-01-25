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
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test16(t *testing.T) {

	Convey("merge scripts", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService) {

			// clear database
			// ------------------------------------------------
			//_ = migrations.Purge()

			scripts := GetScripts(ctx, scriptService, adaptors, []int{29, 30, 31}...)

			engine, err := scriptService.NewEngine(&m.Script{
				Lang: "javascript",
			})
			So(err, ShouldBeNil)

			// a = 1
			result, err := engine.EvalScript(scripts["script29"])
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "1")

			// b = 2
			result, err = engine.EvalScript(scripts["script30"])
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "2")

			// 3 + a
			result, err = engine.EvalScript(scripts["script31"])
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "4")

			_, err = engine.Do()
			So(err, ShouldBeNil)

			// a + b
			result, err = engine.AssertFunction("plus")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "3")

			// 3 + a
			result, err = engine.EvalString("m")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "4")
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
