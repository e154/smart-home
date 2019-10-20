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
// 				(script4, script5)
//
// add workflow scenarios
// 				(wfScenario1 + script6)
// 				(wfScenario2 + script7)
//
// add flow (flow1)
// +----------+
// | handler  |
// | script8  |
// |          |
// +----------+
//
// run core
//
func Test5(t *testing.T) {

	counter := 0

	pool := []string{
		"main workflow",
		"main workflow desc",
		"foo",
		"wf_scenario_1",
		"wf scenario 1",
		"foo",
		"main workflow",
		"main workflow desc",
		"bar",
		"main workflow",
		"main workflow desc",
		"fee",
		"exit from wf_scenario_1",
		"main workflow",
		"main workflow desc",
		"foo",
		"wf_scenario_2",
		"wf scenario 2",
		"foo",
		"enter to wf_scenario_2",
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
				So(v, ShouldEqual, pool[counter])
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

			//
			script4 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test4",
				Source:      coffeeScript4,
				Description: "test4",
			}

			engine4, err := scriptService.NewEngine(script4)
			So(err, ShouldBeNil)
			err = engine4.Compile()
			So(err, ShouldBeNil)
			script4Id, err := adaptors.Script.Add(script4)
			So(err, ShouldBeNil)
			script4, err = adaptors.Script.GetById(script4Id)
			So(err, ShouldBeNil)

			script5 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test5",
				Source:      coffeeScript5,
				Description: "test5",
			}

			engine5, err := scriptService.NewEngine(script5)
			So(err, ShouldBeNil)
			err = engine5.Compile()
			So(err, ShouldBeNil)
			script5Id, err := adaptors.Script.Add(script5)
			So(err, ShouldBeNil)
			script5, err = adaptors.Script.GetById(script5Id)
			So(err, ShouldBeNil)

			script6 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test6",
				Source:      coffeeScript6,
				Description: "test6",
			}

			engine6, err := scriptService.NewEngine(script6)
			So(err, ShouldBeNil)
			err = engine6.Compile()
			So(err, ShouldBeNil)
			script6Id, err := adaptors.Script.Add(script6)
			So(err, ShouldBeNil)
			script6, err = adaptors.Script.GetById(script6Id)
			So(err, ShouldBeNil)

			script7 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test7",
				Source:      coffeeScript7,
				Description: "test7",
			}

			engine7, err := scriptService.NewEngine(script7)
			So(err, ShouldBeNil)
			err = engine7.Compile()
			So(err, ShouldBeNil)
			script7Id, err := adaptors.Script.Add(script7)
			So(err, ShouldBeNil)
			script7, err = adaptors.Script.GetById(script7Id)
			So(err, ShouldBeNil)

			script8 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test8",
				Source:      coffeeScript8,
				Description: "test8",
			}

			engine8, err := scriptService.NewEngine(script8)
			So(err, ShouldBeNil)
			err = engine8.Compile()
			So(err, ShouldBeNil)
			script8Id, err := adaptors.Script.Add(script8)
			So(err, ShouldBeNil)
			script8, err = adaptors.Script.GetById(script8Id)
			So(err, ShouldBeNil)

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

			err = adaptors.Workflow.AddScript(workflow, script4)
			So(err, ShouldBeNil)

			err = adaptors.Workflow.AddScript(workflow, script5)
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

			err = adaptors.WorkflowScenario.AddScript(wfScenario1, script6)
			So(err, ShouldBeNil)

			wfScenario2 := &m.WorkflowScenario{
				Name:       "wf scenario 2",
				SystemName: "wf_scenario_2",
				WorkflowId: workflow.Id,
			}

			wfScenario2.Id, err = adaptors.WorkflowScenario.Add(wfScenario2)
			So(err, ShouldBeNil)

			err = adaptors.WorkflowScenario.AddScript(wfScenario2, script7)
			So(err, ShouldBeNil)

			err = adaptors.Workflow.SetScenario(workflow, wfScenario1)
			So(err, ShouldBeNil)

			// add flow1
			// +----------+
			// | handler  |
			// | script7  |
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
				ScriptId:      &script8.Id,
			}

			ok, _ = feHandler.Valid()
			So(ok, ShouldEqual, true)

			feHandler.Uuid, err = adaptors.FlowElement.Add(feHandler)
			So(err, ShouldBeNil)

			// run
			// ------------------------------------------------
			err = c.Run()
			So(err, ShouldBeNil)

			time.Sleep(time.Second * 1)

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

			Println("send message ...")
			err = flowCore.NewMessage(ctx1)
			So(err, ShouldBeNil)
		})
	})
}
