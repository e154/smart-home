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
