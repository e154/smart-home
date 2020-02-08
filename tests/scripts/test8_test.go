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
// 				(script18, script19)
//
// add workflow scenarios
// 				(wfScenario1 + script20)
// 				(wfScenario2 + script21)
//
// add flow (flow1)
// +----------+
// | handler  |
// | script22 |
// |          |
// +----------+
//
// add flow (flow2)
// +----------+
// | handler  |
// | script23 |
// |          |
// +----------+
//
// select scenario from:
//				- flow script
//
func Test8(t *testing.T) {

	counter := 0

	pool := []string{
		"enter script18",
		"enter script19",
		"enter script20",
		"enter script22",
		"exit script20",
		"enter script18",
		"enter script19",
		"enter script21",
		"exit script21",
	}

	initCallback := func(ctx C) {
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
	}

	Convey("workflow bind", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			c *core.Core) {

			initCallback(ctx)

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

			scripts := GetScripts(ctx, scriptService, adaptors, 18,19,20,21,22,23)

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

			err = adaptors.Workflow.AddScript(workflow, scripts["script18"])
			So(err, ShouldBeNil)

			err = adaptors.Workflow.AddScript(workflow, scripts["script19"])
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

			err = adaptors.WorkflowScenario.AddScript(wfScenario1, scripts["script20"])
			So(err, ShouldBeNil)

			wfScenario2 := &m.WorkflowScenario{
				Name:       "wf scenario 2",
				SystemName: "wf_scenario_2",
				WorkflowId: workflow.Id,
			}

			wfScenario2.Id, err = adaptors.WorkflowScenario.Add(wfScenario2)
			So(err, ShouldBeNil)

			err = adaptors.WorkflowScenario.AddScript(wfScenario2, scripts["script21"])
			So(err, ShouldBeNil)

			err = adaptors.Workflow.SetScenario(workflow, wfScenario1)
			So(err, ShouldBeNil)

			// add flow1
			// +----------+
			// | handler  |
			// | script22 |
			// |          |
			// +----------+
			flow1 := &m.Flow{
				Name:               "flow1",
				Status:             Enabled,
				WorkflowId:         workflow.Id,
				WorkflowScenarioId: wfScenario1.Id,
			}

			ok, _ := flow1.Valid()
			So(ok, ShouldEqual, true)

			flow1.Id, err = adaptors.Flow.Add(flow1)
			So(err, ShouldBeNil)

			feHandler := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &scripts["script22"].Id,
			}

			ok, _ = feHandler.Valid()
			So(ok, ShouldEqual, true)

			feHandler.Uuid, err = adaptors.FlowElement.Add(feHandler)
			So(err, ShouldBeNil)

			// add flow2
			// +----------+
			// | handler  |
			// | script23 |
			// |          |
			// +----------+
			flow2 := &m.Flow{
				Name:               "flow2",
				Status:             Enabled,
				WorkflowId:         workflow.Id,
				WorkflowScenarioId: wfScenario2.Id,
			}

			ok, _ = flow2.Valid()
			So(ok, ShouldEqual, true)

			flow2.Id, err = adaptors.Flow.Add(flow2)
			So(err, ShouldBeNil)

			feHandler2 := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow2.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &scripts["script23"].Id,
			}

			ok, _ = feHandler2.Valid()
			So(ok, ShouldEqual, true)

			feHandler2.Uuid, err = adaptors.FlowElement.Add(feHandler2)
			So(err, ShouldBeNil)

			// run
			// ------------------------------------------------
			err = c.Run()
			So(err, ShouldBeNil)

			workflowCore, err := c.GetWorkflow(workflow.Id)
			So(err, ShouldBeNil)

			flowCore, err := workflowCore.GetFLow(flow1.Id)
			So(err, ShouldBeNil)

			message := core.NewMessage()
			message.SetVar("val", 1)

			// create context
			var ctx1 context.Context
			ctx1, _ = context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
			ctx1 = context.WithValue(ctx1, "msg", message)

			// send message ...
			err = flowCore.NewMessage(ctx1)
			So(err, ShouldBeNil)

			time.Sleep(time.Second * 5)

			err = c.Stop()
			So(err, ShouldBeNil)
		})
	})
}
