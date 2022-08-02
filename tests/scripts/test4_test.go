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

func Test4(t *testing.T) {

	var state string
	store = func(i interface{}) {
		state = fmt.Sprintf("%v", i)
	}

	Convey("javascript PushStruct", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migratxions *migrations.Migrations,
			scriptService scripts.ScriptService) {

			storeRegisterCallback(scriptService)

			script1 := &m.Script{
				Lang: common.ScriptLangCoffee,
			}
			engine, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)

			for i := 0; i < 2; i++ {
				engine.PushStruct("test", &MyStruct{
					Int:     42,
					Float64: 21.0,
					Nested:  &MyStruct{Int: 21},
				})
			}

			_, err = engine.EvalString(`store([
				test.int,
				test.multiply(2),
				test.nested.int,
				test.nested.multiply(3)
			])`)
			So(err, ShouldBeNil)

			So(state, ShouldEqual, "[42 84 21 63]")
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
