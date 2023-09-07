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

package plugins

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/debug"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestArea(t *testing.T) {


	Convey("area", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			automation automation.Automation,
			eventBus bus.Bus) {

			//err := migrations.Purge()
			//So(err, ShouldBeNil)

			var err error

			area := &m.Area{
				Name:        "zone 51",
				Description: "zone 51",
				Polygon:     []m.Point{
					{75.1, 29.2},
					{77.1, 29.2},
					{77.1, 29.2},
					{75.1, 29.2},
				},
			}
			area.Id, err = adaptors.Area.Add(area)
			So(err, ShouldBeNil)

			area2, err := adaptors.Area.GetById(area.Id)
			So(err, ShouldBeNil)

			debug.Println(area2)
		})
	})
}
