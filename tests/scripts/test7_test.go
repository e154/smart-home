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

package scripts

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
// workflow javascript bindings
//
// create workflow
// 				(script9)
//
// add workflow scenario
// 				(wfScenario1 + script10)
//
// add flow (flow1)
// +----------+     +----------+    +----------+
// | handler  |     |   task   |    |  emitter |
// | script11 +-----> script12 +----> script13 |
// |          |     |          |    |          |
// +----------+     +----------+    +----------+
//
// add flow (flow2)
// +----------+     +----------+    +----------+
// | handler  |     |   task   |    |  emitter |
// | script14 +-----> script15 +----> script16 |
// |          |     |          |    |          |
// +----------+     +----------+    +----------+
//
// scope of variables between flows and workflow and scenarios
//
func Test7(t *testing.T) {

	counter := 0

	pool := []string{
		"enter script9",
		"enter script10",
		"enter script11",
		"<nil>",
		"foo",
		"enter script12",
		"foo",
		"enter script13",
		"bar",
		"bar",
		"enter script11",
		"bar",
		"foo",
		"enter script12",
		"foo",
		"enter script13",
		"bar",
		"bar",
		"enter script14",
		"enter script15",
		"enter script16",
		"exit script10",
	}

	initCallback := func(ctx C, scriptService *scripts.ScriptService) {
		store = func(i interface{}) {
			v := fmt.Sprintf("%v", i)
			//fmt.Println("v:", v)

			if counter >= len(pool) {
				fmt.Println("========= WARNING =========")
				fmt.Printf("counter(%d), v: %v\n", counter, v)
				return
			}

			switch counter {
			default:
				ctx.So(v, ShouldEqual, pool[counter])
			}

			counter++
		}

		scriptService.PushFunctions("So", func(actual interface{}, assert string, expected interface{}) {
			switch assert {
			case "ShouldEqual":
				So(actual, ShouldEqual, expected)
			}

		})
	}

	Convey("workflow bind", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			c *core.Core) {

			initCallback(ctx, scriptService)

			// stop core
			// ------------------------------------------------
			err := c.Stop()
			So(err, ShouldBeNil)

			// clear database
			// ------------------------------------------------
			migrations.Purge()

			// scripts
			// ------------------------------------------------
			storeRegisterCallback(scriptService)

			scripts := GetScripts(ctx, scriptService, adaptors, 9, 10, 11, 12, 13, 14, 15, 16)

			// workflow
			// ------------------------------------------------
			workflow := &m.Workflow{
				Name:        "main workflow",
				Description: "main workflow desc",
				Status:      "enabled",
			}

			wfId, err := adaptors.Workflow.Add(workflow)
			So(err, ShouldBeNil)
			workflow.Id = wfId

			err = adaptors.Workflow.AddScript(workflow, scripts["script9"])
			So(err, ShouldBeNil)

			// add workflow scenario
			// ------------------------------------------------
			wfScenario1 := &m.WorkflowScenario{
				Name:       "wf scenario 1",
				SystemName: "wf_scenario_1",
				WorkflowId: workflow.Id,
			}

			wfScenario1.Id, err = adaptors.WorkflowScenario.Add(wfScenario1)
			So(err, ShouldBeNil)

			err = adaptors.WorkflowScenario.AddScript(wfScenario1, scripts["script10"])
			So(err, ShouldBeNil)

			err = adaptors.Workflow.SetScenario(workflow, wfScenario1)
			So(err, ShouldBeNil)

			// add flow (flow1)
			// +----------+     +----------+    +----------+
			// | handler  |     |  task    |    |  emitter |
			// | script11 +-----> script12 +----> script13 |
			// |          |     |          |    |          |
			// +----------+     +----------+    +----------+
			flow1 := &m.Flow{
				Name:               "flow1",
				Description:        "flow1 desc",
				Status:             Enabled,
				WorkflowId:         workflow.Id,
				WorkflowScenarioId: wfScenario1.Id,
			}
			ok, _ := flow1.Valid()
			So(ok, ShouldEqual, true)

			flow1.Id, err = adaptors.Flow.Add(flow1)
			So(err, ShouldBeNil)

			// add handler
			feHandler := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &scripts["script11"].Id,
				GraphSettings: m.FlowElementGraphSettings{
					Position: m.FlowElementGraphSettingsPosition{
						Top:  180,
						Left: 180,
					},
				},
			}
			feEmitter := &m.FlowElement{
				Name:          "emitter",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageEmitter,
				ScriptId:      &scripts["script13"].Id,
				GraphSettings: m.FlowElementGraphSettings{
					Position: m.FlowElementGraphSettingsPosition{
						Top:  180,
						Left: 560,
					},
				},
			}
			feTask1 := &m.FlowElement{
				Name:          "task",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeTask,
				ScriptId:      &scripts["script12"].Id,
				GraphSettings: m.FlowElementGraphSettings{
					Position: m.FlowElementGraphSettingsPosition{
						Top:  160,
						Left: 340,
					},
				},
			}
			ok, _ = feHandler.Valid()
			So(ok, ShouldEqual, true)
			ok, _ = feEmitter.Valid()
			So(ok, ShouldEqual, true)
			ok, _ = feTask1.Valid()
			So(ok, ShouldEqual, true)

			feHandler.Uuid, err = adaptors.FlowElement.Add(feHandler)
			So(err, ShouldBeNil)
			feEmitter.Uuid, err = adaptors.FlowElement.Add(feEmitter)
			So(err, ShouldBeNil)
			feTask1.Uuid, err = adaptors.FlowElement.Add(feTask1)
			So(err, ShouldBeNil)

			connect1 := &m.Connection{
				Name:        "con1",
				ElementFrom: feHandler.Uuid,
				ElementTo:   feTask1.Uuid,
				FlowId:      flow1.Id,
				PointFrom:   1,
				PointTo:     10,
			}
			connect2 := &m.Connection{
				Name:        "con2",
				ElementFrom: feTask1.Uuid,
				ElementTo:   feEmitter.Uuid,
				FlowId:      flow1.Id,
				PointFrom:   4,
				PointTo:     3,
			}

			ok, _ = connect1.Valid()
			So(ok, ShouldEqual, true)
			ok, _ = connect2.Valid()
			So(ok, ShouldEqual, true)

			connect1.Uuid, err = adaptors.Connection.Add(connect1)
			So(err, ShouldBeNil)
			connect2.Uuid, err = adaptors.Connection.Add(connect2)
			So(err, ShouldBeNil)

			// add flow (flow2)
			// +----------+     +----------+    +----------+
			// | handler  |     |   task   |    |  emitter |
			// | script14 +-----> script15 +----> script16 |
			// |          |     |          |    |          |
			// +----------+     +----------+    +----------+
			flow2 := &m.Flow{
				Name:               "flow2",
				Description:        "flow2 desc",
				Status:             Enabled,
				WorkflowId:         workflow.Id,
				WorkflowScenarioId: wfScenario1.Id,
			}
			ok, _ = flow2.Valid()
			So(ok, ShouldEqual, true)

			flow2.Id, err = adaptors.Flow.Add(flow2)
			So(err, ShouldBeNil)

			// add handler
			feHandler2 := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow2.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &scripts["script14"].Id,
				GraphSettings: m.FlowElementGraphSettings{
					Position: m.FlowElementGraphSettingsPosition{
						Top:  180,
						Left: 180,
					},
				},
			}
			feEmitter2 := &m.FlowElement{
				Name:          "emitter",
				FlowId:        flow2.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageEmitter,
				ScriptId:      &scripts["script16"].Id,
				GraphSettings: m.FlowElementGraphSettings{
					Position: m.FlowElementGraphSettingsPosition{
						Top:  180,
						Left: 560,
					},
				},
			}
			feTask2 := &m.FlowElement{
				Name:          "task",
				FlowId:        flow2.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeTask,
				ScriptId:      &scripts["script15"].Id,
				GraphSettings: m.FlowElementGraphSettings{
					Position: m.FlowElementGraphSettingsPosition{
						Top:  160,
						Left: 340,
					},
				},
			}
			ok, _ = feHandler2.Valid()
			So(ok, ShouldEqual, true)
			ok, _ = feEmitter2.Valid()
			So(ok, ShouldEqual, true)
			ok, _ = feTask2.Valid()
			So(ok, ShouldEqual, true)

			feHandler2.Uuid, err = adaptors.FlowElement.Add(feHandler2)
			So(err, ShouldBeNil)
			feEmitter2.Uuid, err = adaptors.FlowElement.Add(feEmitter2)
			So(err, ShouldBeNil)
			feTask2.Uuid, err = adaptors.FlowElement.Add(feTask2)
			So(err, ShouldBeNil)

			connect3 := &m.Connection{
				Name:        "con3",
				ElementFrom: feHandler2.Uuid,
				ElementTo:   feTask2.Uuid,
				FlowId:      flow2.Id,
				PointFrom:   1,
				PointTo:     10,
			}
			connect4 := &m.Connection{
				Name:        "con4",
				ElementFrom: feTask2.Uuid,
				ElementTo:   feEmitter2.Uuid,
				FlowId:      flow2.Id,
				PointFrom:   4,
				PointTo:     3,
			}

			ok, _ = connect3.Valid()
			So(ok, ShouldEqual, true)
			ok, _ = connect4.Valid()
			So(ok, ShouldEqual, true)

			connect3.Uuid, err = adaptors.Connection.Add(connect3)
			So(err, ShouldBeNil)
			connect4.Uuid, err = adaptors.Connection.Add(connect4)
			So(err, ShouldBeNil)

			// run
			// ------------------------------------------------
			err = c.Run()
			So(err, ShouldBeNil)

			workflowCore, err := c.GetWorkflow(workflow.Id)
			So(err, ShouldBeNil)

			flowCore1, err := workflowCore.GetFLow(flow1.Id)
			So(err, ShouldBeNil)

			flowCore2, err := workflowCore.GetFLow(flow2.Id)
			So(err, ShouldBeNil)

			// flow1
			for i := 0; i < 2; i++ {
				message := core.NewMessage()
				message.SetVar("val", 1)

				ctx1, _ := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
				ctx1 = context.WithValue(ctx1, "msg", message)

				// send message ...
				err = flowCore1.NewMessage(ctx1)
				So(err, ShouldBeNil)
			}

			// flow2
			message := core.NewMessage()
			message.SetVar("val", 2)

			ctx2, _ := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
			ctx2 = context.WithValue(ctx2, "msg", message)

			// send message ...
			err = flowCore2.NewMessage(ctx2)
			So(err, ShouldBeNil)

			So(flowCore2.GetMessage().GetVar("val"), ShouldEqual, 123)

			//time.Sleep(time.Second * 5)

			err = c.Stop()
			So(err, ShouldBeNil)
		})
	})
}
