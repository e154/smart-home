package workflow

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/common"
	. "github.com/e154/smart-home/models/devices"
	"fmt"
	"github.com/e154/smart-home/common/debug"
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
		container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			c *core.Core) {

			// clear database
			migrations.Purge()

			// add node
			node := &m.Node{
				Name:   "node",
				Ip:     "127.0.0.1",
				Port:   3001,
				Status: "enabled",
			}
			ok, _ := node.Valid()
			So(ok, ShouldEqual, true)

			var err error
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
				Baud: 19200,
				Timeout: 457,
				StopBits: 2,
				Sleep: 0,
			}

			ok, _ = parentDevice.SetProperties(smartBusConfig1)
			So(ok, ShouldEqual, true)
			parentDevice.Id, err = adaptors.Device.Add(parentDevice)
			So(err, ShouldBeNil)
			parentDevice, err = adaptors.Device.GetById(parentDevice.Id)
			So(err, ShouldBeNil)

			fmt.Println("----")
			debug.Println(parentDevice)
			fmt.Println("----")
		})
	})
}
