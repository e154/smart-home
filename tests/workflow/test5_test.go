// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package workflow

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/models/devices"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// create node
//
// create device
// 			add device property
//
//
//
func Test5(t *testing.T) {

	Convey("add scripts", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			c *core.Core) {

			// stop core
			// ------------------------------------------------
			err := c.Stop()
			So(err, ShouldBeNil)

			// clear database
			migrations.Purge()

			// add node
			node := &m.Node{
				Name:     "node",
				Login:    "node",
				Password: "node",
				Status:   "enabled",
			}
			ok, _ := node.Valid()
			So(ok, ShouldEqual, true)

			node.Id, err = adaptors.Node.Add(node)
			So(err, ShouldBeNil)

			// add parent device
			parentDevice := &m.Device{
				Name:       "device",
				Status:     "enabled",
				Type:       "default",
				Node:       node,
				Properties: []byte("{}"),
			}
			ok, _ = parentDevice.Valid()
			So(ok, ShouldEqual, true)

			smartBusConfig1 := &DevSmartBusConfig{
				Baud:     19200,
				Timeout:  457,
				StopBits: 2,
				Sleep:    0,
			}

			ok, _ = parentDevice.SetProperties(smartBusConfig1)
			So(ok, ShouldEqual, true)
			parentDevice.Id, err = adaptors.Device.Add(parentDevice)
			So(err, ShouldBeNil)
			parentDevice, err = adaptors.Device.GetById(parentDevice.Id)
			So(err, ShouldBeNil)

			err = c.Stop()
			So(err, ShouldBeNil)
		})
	})
}
