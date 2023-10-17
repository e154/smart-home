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
	m "github.com/e154/smart-home/models"
	"testing"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

func fooBind(num int64) func(j string) string {
	return func(j string) string {
		return fmt.Sprintf("%d_%s", num, j)
	}
}

func Test14(t *testing.T) {

	Convey("scripts run syn command", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService) {

			script1 := &m.Script{
				Lang:   "coffeescript",
				Name:   "test28",
				Source: "foo 'bar'",
			}

			scriptService.PushFunctions("foo", fooBind(12))

			engine1, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)
			err = engine1.Compile()
			So(err, ShouldBeNil)

			result, err := engine1.Do()
			So(err, ShouldBeNil)

			So(result, ShouldEqual, "12_bar")
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
