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
	"context"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

//
// create workflow
//
// add workflow scenarios (wf_scenario_1 + script27)
//
// add flow (flow1)
// +----------+        +----------+
// | handler  |  con1  |  emitter |
// | script25 +--------> script26 |
// |          |        |          |
// +----------+        +-----+----+
//      ^                    |
//      |        con2        |
//      +--------------------+
//
//
// run core
//
func Test12(t *testing.T) {

	var story = make([]string, 0)
	//var scriptCounter string

	store = func(i interface{}) {
		cmd := fmt.Sprintf("%v", i)

		story = append(story, cmd)
	}

	//store2 = func(i interface{}) {
	//	scriptCounter = fmt.Sprintf("%v", i)
	//}

	Convey("detect circle flow", t, func(ctx C) {
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
			err = migrations.Purge()
			So(err, ShouldBeNil)

			// create scripts
			// ------------------------------------------------
			script25 := &m.Script{
				Lang:        "coffeescript",
				Name:        "script25",
				Source:      coffeeScript25,
				Description: "script25",
			}
			script26 := &m.Script{
				Lang:        "coffeescript",
				Name:        "script26",
				Source:      coffeeScript26,
				Description: "script26",
			}
			script27 := &m.Script{
				Lang:        "coffeescript",
				Name:        "script27",
				Source:      coffeeScript7,
				Description: "script27",
			}

			ok, _ := script25.Valid()
			So(ok, ShouldEqual, true)

			ok, _ = script26.Valid()
			So(ok, ShouldEqual, true)

			ok, _ = script27.Valid()
			So(ok, ShouldEqual, true)

			engine25, err := scriptService.NewEngine(script25)
			So(err, ShouldBeNil)
			err = engine25.Compile()
			So(err, ShouldBeNil)
			script25Id, err := adaptors.Script.Add(script25)
			So(err, ShouldBeNil)
			script25, err = adaptors.Script.GetById(script25Id)
			So(err, ShouldBeNil)

			engine26, err := scriptService.NewEngine(script26)
			So(err, ShouldBeNil)
			err = engine26.Compile()
			So(err, ShouldBeNil)
			script26Id, err := adaptors.Script.Add(script26)
			So(err, ShouldBeNil)
			script26, err = adaptors.Script.GetById(script26Id)
			So(err, ShouldBeNil)

			engine27, err := scriptService.NewEngine(script27)
			So(err, ShouldBeNil)
			err = engine27.Compile()
			So(err, ShouldBeNil)
			script27Id, err := adaptors.Script.Add(script27)
			So(err, ShouldBeNil)
			script27, err = adaptors.Script.GetById(script27Id)
			So(err, ShouldBeNil)

			// create workflow
			// ------------------------------------------------
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
			// ------------------------------------------------
			wfScenario1 := &m.WorkflowScenario{
				Name:       "wf scenario 1",
				SystemName: "wf_scenario_1",
				WorkflowId: workflow.Id,
			}

			ok, _ = wfScenario1.Valid()
			So(ok, ShouldEqual, true)

			wfScenarioId1, err := adaptors.WorkflowScenario.Add(wfScenario1)
			So(err, ShouldBeNil)
			wfScenario1.Id = wfScenarioId1

			err = adaptors.WorkflowScenario.AddScript(wfScenario1, script27)
			So(err, ShouldBeNil)

			workflow.Scenario = wfScenario1
			err = adaptors.Workflow.Update(workflow)
			So(err, ShouldBeNil)

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

			feEmitter := &m.FlowElement{
				Name:          "emitter",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageEmitter,
				ScriptId:      &script25.Id,
			}

			feHandler := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &script26.Id,
			}

			ok, _ = feEmitter.Valid()
			So(ok, ShouldEqual, true)

			ok, _ = feHandler.Valid()
			So(ok, ShouldEqual, true)

			feEmitter.Uuid, err = adaptors.FlowElement.Add(feEmitter)
			So(err, ShouldBeNil)

			feHandler.Uuid, err = adaptors.FlowElement.Add(feHandler)
			So(err, ShouldBeNil)

			con1 := &m.Connection{
				Name:        "con1",
				ElementFrom: feHandler.Uuid,
				ElementTo:   feEmitter.Uuid,
				FlowId:      flow1.Id,
				PointFrom:   1,
				PointTo:     1,
			}

			ok, _ = con1.Valid()
			So(ok, ShouldEqual, true)

			con1.Uuid, err = adaptors.Connection.Add(con1)
			So(err, ShouldBeNil)

			con2 := &m.Connection{
				Name:        "con2",
				ElementFrom: feEmitter.Uuid,
				ElementTo:   feHandler.Uuid,
				FlowId:      flow1.Id,
				PointFrom:   2,
				PointTo:     2,
			}

			ok, _ = con2.Valid()
			So(ok, ShouldEqual, true)

			con2.Uuid, err = adaptors.Connection.Add(con2)
			So(err, ShouldBeNil)

			err = c.Run()
			So(err, ShouldBeNil)

			workflowCore, err := c.GetWorkflow(workflow.Id)
			So(err, ShouldBeNil)

			flowCore, err := workflowCore.GetFLow(flow1.Id)
			So(err, ShouldBeNil)

			message := core.NewMessage()

			// create context
			ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
			ctx = context.WithValue(ctx, "msg", message)

			circularErr := flowCore.NewMessage(ctx)
			So(circularErr, ShouldNotBeNil)

			err = c.Stop()
			So(err, ShouldBeNil)
		})
	})
}
