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

package area

import (
	"testing"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestArea(t *testing.T) {

	Convey("area", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors) {

			area := &m.Area{
				Name:        "zone 51",
				Description: "zone 51",
				Polygon: []m.Point{
					{75.1, 29.2},
					{77.2, 29.3},
					{77.3, 29.4},
					{75.4, 29.5},
				},
			}

			var err error
			area.Id, err = adaptors.Area.Add(area)
			So(err, ShouldBeNil)

			area, err = adaptors.Area.GetById(area.Id)
			So(err, ShouldBeNil)

			So(area.Name, ShouldEqual, "zone 51")
			So(area.Description, ShouldEqual, "zone 51")
			So(area.Polygon[0].Lat, ShouldEqual, 29.2)
			So(area.Polygon[0].Lon, ShouldEqual, 75.1)
			So(area.Polygon[1].Lat, ShouldEqual, 29.3)
			So(area.Polygon[1].Lon, ShouldEqual, 77.2)
			So(area.Polygon[2].Lat, ShouldEqual, 29.4)
			So(area.Polygon[2].Lon, ShouldEqual, 77.3)
			So(area.Polygon[3].Lat, ShouldEqual, 29.5)
			So(area.Polygon[3].Lon, ShouldEqual, 75.4)
		})
	})
}
