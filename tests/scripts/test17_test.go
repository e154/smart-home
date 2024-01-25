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

func Test17(t *testing.T) {

	Convey("concat scripts", t, func(ctx C) {
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

			result, err := engine.EvalString(scripts["script29"].Compiled + scripts["script30"].Compiled + scripts["script31"].Compiled)
			So(err, ShouldBeNil)
			fmt.Println(result)
			//So(result, ShouldEqual, "1")

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
