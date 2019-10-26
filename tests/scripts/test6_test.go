package scripts

import (
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
// add flow (flow2)
// +----------+
// | handler  |
// | script17 |
// |          |
// +----------+
//
// select scenario from:
//				- workflow scenario script
//
func Test6(t *testing.T) {

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
		"foo",
		"wf_scenario_2",
		"wf scenario 2",
		"foo",
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

			scripts := GetScripts(ctx, scriptService, adaptors, 4, 5, 6, 7, 8, 17)

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

			err = adaptors.Workflow.AddScript(workflow, scripts["script4"])
			So(err, ShouldBeNil)

			err = adaptors.Workflow.AddScript(workflow, scripts["script5"])
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

			err = adaptors.WorkflowScenario.AddScript(wfScenario1, scripts["script6"])
			So(err, ShouldBeNil)

			wfScenario2 := &m.WorkflowScenario{
				Name:       "wf scenario 2",
				SystemName: "wf_scenario_2",
				WorkflowId: workflow.Id,
			}

			wfScenario2.Id, err = adaptors.WorkflowScenario.Add(wfScenario2)
			So(err, ShouldBeNil)

			err = adaptors.WorkflowScenario.AddScript(wfScenario2, scripts["script7"])
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
				ScriptId:      &scripts["script8"].Id,
			}

			ok, _ = feHandler.Valid()
			So(ok, ShouldEqual, true)

			feHandler.Uuid, err = adaptors.FlowElement.Add(feHandler)
			So(err, ShouldBeNil)

			// add flow2
			// +----------+
			// | handler  |
			// | script8  |
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
				ScriptId:      &scripts["script17"].Id,
			}

			ok, _ = feHandler2.Valid()
			So(ok, ShouldEqual, true)

			feHandler2.Uuid, err = adaptors.FlowElement.Add(feHandler2)
			So(err, ShouldBeNil)

			// run
			// ------------------------------------------------
			err = c.Run()
			So(err, ShouldBeNil)

			//time.Sleep(time.Second * 2)

			workflowCore, err := c.GetWorkflow(workflow.Id)
			So(err, ShouldBeNil)

			//time.Sleep(time.Second * 2)

			_, err = workflowCore.GetFLow(flow1.Id)
			So(err, ShouldNotBeNil)

			//message := core.NewMessage()
			//message.SetVar("val", 1)
			//
			//// create context
			//var ctx1 context.Context
			//ctx1, _ = context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
			//ctx1 = context.WithValue(ctx1, "msg", message)
			//
			//Println("send message ...")
			//err = flowCore.NewMessage(ctx1)
			//So(err, ShouldBeNil)

			time.Sleep(time.Second * 2)

			err = c.Stop()
			So(err, ShouldBeNil)
		})
	})
}
