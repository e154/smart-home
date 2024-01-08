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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

func Test16(t *testing.T) {

	Convey("scripts run syn command", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService) {

			script1 := &m.Script{
				Lang:   common.ScriptLangJavascript,
				Source: javascript29,
			}
			engine, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)

			So(err, ShouldBeNil)
			err = engine.Compile()
			So(err, ShouldBeNil)

			result, err := engine.Do()
			So(err, ShouldBeNil)
			fmt.Println(result)

			result, err = engine.AssertFunction("foo")
			So(err, ShouldBeNil)

			time.Sleep(time.Second * 3)

		})
		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
