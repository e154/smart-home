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
	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// create node
//
// create parent device
// 			child device x 2
//
// create parent device-actions (+script10)
//
// add workflow scenarios (wf_scenario_1)
//
// add flow (flow1)
// +----------+
// | handler  |
// | script11 |
// |          |
// +----------+
//
// add worker
//
//
func Test4(t *testing.T) {

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
			// ------------------------------------------------
			migrations.Purge()

			// add node
			// ------------------------------------------------
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
			// ------------------------------------------------
			parentDevice := &m.Device{
				Name:       "device",
				Status:     "enabled",
				Type:       "default",
				Node:       node,
				Properties: []byte("{}"),
			}
			ok, _ = parentDevice.Valid()
			So(ok, ShouldEqual, true)

			parentDevice.Id, err = adaptors.Device.Add(parentDevice)
			So(err, ShouldBeNil)

			// add child device
			device1 := &m.Device{
				Name:       "device",
				Status:     "enabled",
				Type:       "default",
				Device:     parentDevice,
				Properties: []byte("{}"),
			}
			device2 := &m.Device{
				Name:       "device",
				Status:     "enabled",
				Type:       "default",
				Device:     parentDevice,
				Properties: []byte("{}"),
			}
			ok, _ = device1.Valid()
			So(ok, ShouldEqual, true)
			ok, _ = device2.Valid()
			So(ok, ShouldEqual, true)

			device1.Id, err = adaptors.Device.Add(device1)
			So(err, ShouldBeNil)
			device2.Id, err = adaptors.Device.Add(device2)
			So(err, ShouldBeNil)

			// add script
			// ------------------------------------------------
			script10 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test10",
				Source:      coffeeScript10,
				Description: "test10",
			}
			script11 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test11",
				Source:      coffeeScript11,
				Description: "test11",
			}

			ok, _ = script10.Valid()
			So(ok, ShouldEqual, true)
			ok, _ = script11.Valid()
			So(ok, ShouldEqual, true)

			engine10, err := scriptService.NewEngine(script10)
			So(err, ShouldBeNil)
			err = engine10.Compile()
			So(err, ShouldBeNil)
			script10Id, err := adaptors.Script.Add(script10)
			So(err, ShouldBeNil)
			script10, err = adaptors.Script.GetById(script10Id)
			So(err, ShouldBeNil)

			engine11, err := scriptService.NewEngine(script11)
			So(err, ShouldBeNil)
			err = engine11.Compile()
			So(err, ShouldBeNil)
			script11Id, err := adaptors.Script.Add(script11)
			So(err, ShouldBeNil)
			script11, err = adaptors.Script.GetById(script11Id)
			So(err, ShouldBeNil)

			// add device action
			// ------------------------------------------------
			deviceAction := &m.DeviceAction{
				Name:     "deviceAction",
				DeviceId: parentDevice.Id,
				ScriptId: script10Id,
			}
			deviceAction.Id, err = adaptors.DeviceAction.Add(deviceAction)
			So(err, ShouldBeNil)

			// add workflow
			workflow := &m.Workflow{
				Name:        "main workflow",
				Description: "main workflow desc",
				Status:      "enabled",
			}

			ok, _ = workflow.Valid()
			So(ok, ShouldEqual, true)

			wfId, err := adaptors.Workflow.Add(workflow)
			So(err, ShouldBeNil)
			workflow.Id = wfId

			// add workflow scenario
			wfScenario1 := &m.WorkflowScenario{
				Name:       "wf scenario 1",
				SystemName: "wf_scenario_1",
				WorkflowId: workflow.Id,
			}

			ok, _ = wfScenario1.Valid()
			So(ok, ShouldEqual, true)

			wfScenario1.Id, err = adaptors.WorkflowScenario.Add(wfScenario1)
			So(err, ShouldBeNil)

			err = adaptors.Workflow.SetScenario(workflow, wfScenario1.Id)
			So(err, ShouldBeNil)

			// add flow1
			// +----------+
			// | handler  |
			// | script11 |
			// |          |
			// +----------+
			flow1 := &m.Flow{
				Name:               "flow1",
				Status:             Enabled,
				WorkflowId:         workflow.Id,
				WorkflowScenarioId: wfScenario1.Id,
			}

			ok, _ = flow1.Valid()
			So(ok, ShouldEqual, true)

			flow1.Id, err = adaptors.Flow.Add(flow1)
			So(err, ShouldBeNil)

			feHandler := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &script11.Id,
			}

			ok, _ = feHandler.Valid()
			So(ok, ShouldEqual, true)

			feHandler.Uuid, err = adaptors.FlowElement.Add(feHandler)
			So(err, ShouldBeNil)

			// add worker
			worker := &m.Worker{
				Name:           "worker",
				Time:           "* * * * * *",
				Status:         "enabled",
				WorkflowId:     workflow.Id,
				FlowId:         flow1.Id,
				DeviceActionId: deviceAction.Id,
			}

			ok, _ = worker.Valid()
			So(ok, ShouldEqual, true)

			worker.Id, err = adaptors.Worker.Add(worker)
			So(err, ShouldBeNil)

			// get flow
			// ------------------------------------------------
			err = c.Run()
			So(err, ShouldBeNil)

			err = c.Stop()
			So(err, ShouldBeNil)
		})
	})
}
