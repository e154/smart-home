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
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

func Test2(t *testing.T) {

	//var state string
	//store = func(i interface{}) {
	//	state = fmt.Sprintf("%v", i)
	//}

	//var script1 *m.Script
	Convey("require external library", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService) {

			//todo: fix
			//storeRegisterCallback(scriptService)
			//
			//script1 = &m.Script{
			//	Lang:        "coffeescript",
			//	Name:        "test2",
			//	Source:      coffeeScript2,
			//	Description: "test2",
			//}
			//
			//engine1, err := scriptService.NewEngine(script1)
			//So(err, ShouldBeNil)
			//err = engine1.Compile()
			//So(err, ShouldBeNil)
			//
			//_, err = engine1.Do()
			//So(err, ShouldBeNil)
			//
			//So(state, ShouldEqual, "123-bar-Jan")
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
