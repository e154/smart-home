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

// create workflow
//
// add workflow scenarios (wf_scenario_1)
//
// add flow (flow1)
//
// +----------+     +----------+    +----------+
// | handler  |     |  task    |    |  emitter |
// | script12 +-----> script13 +----> script14 |
// |          |     |          |    |          |
// +----------+     +----------+    +----------+
//
func Test6(t *testing.T) {

	var story = make([]string, 0)
	var scriptCounter string

	store = func(i interface{}) {
		cmd := fmt.Sprintf("%v", i)

		story = append(story, cmd)
	}

	store2 = func(i interface{}) {
		scriptCounter = fmt.Sprintf("%v", i)
	}

	Convey("break flow process", t, func(ctx C) {
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

			// add device
			// ------------------------------------------------
			script12 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test12",
				Source:      coffeeScript12,
				Description: "test12",
			}

			engine12, err := scriptService.NewEngine(script12)
			So(err, ShouldBeNil)
			err = engine12.Compile()
			So(err, ShouldBeNil)
			script12Id, err := adaptors.Script.Add(script12)
			So(err, ShouldBeNil)
			script12, err = adaptors.Script.GetById(script12Id)
			So(err, ShouldBeNil)

			storeRegisterCallback(scriptService)

			script13 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test13",
				Source:      coffeeScript13,
				Description: "test13",
			}

			engine13, err := scriptService.NewEngine(script13)
			So(err, ShouldBeNil)
			err = engine13.Compile()
			So(err, ShouldBeNil)
			script13Id, err := adaptors.Script.Add(script13)
			So(err, ShouldBeNil)
			script13, err = adaptors.Script.GetById(script13Id)
			So(err, ShouldBeNil)

			script14 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test14",
				Source:      coffeeScript14,
				Description: "test14",
			}

			engine14, err := scriptService.NewEngine(script14)
			So(err, ShouldBeNil)
			err = engine14.Compile()
			So(err, ShouldBeNil)
			script14Id, err := adaptors.Script.Add(script14)
			So(err, ShouldBeNil)
			script14, err = adaptors.Script.GetById(script14Id)
			So(err, ShouldBeNil)

			// add workflow
			// ------------------------------------------------
			workflow := &m.Workflow{
				Name:        "main workflow",
				Description: "main workflow desc",
				Status:      "enabled",
			}

			ok, _ := workflow.Valid()
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

			wfScenario1.Id, err = adaptors.WorkflowScenario.Add(wfScenario1)
			So(err, ShouldBeNil)

			err = adaptors.Workflow.SetScenario(workflow, wfScenario1)
			So(err, ShouldBeNil)

			// +----------+     +----------+    +----------+
			// | handler  |     |  task    |    |  emitter |
			// | script12 +-----> script13 +----> script14 |
			// |          |     |          |    |          |
			// +----------+     +----------+    +----------+
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

			// add handler
			feHandler := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &script12Id,
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
				ScriptId:      &script14Id,
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
				ScriptId:      &script13Id,
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

			// get flow
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
			var ctx context.Context
			ctx, _ = context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
			ctx = context.WithValue(ctx, "msg", message)

			// send message ...
			err = flowCore.NewMessage(ctx)
			So(err, ShouldBeNil)

			So(len(story), ShouldEqual, 3)
			So(scriptCounter, ShouldEqual, "3")

			err = c.Stop()
			So(err, ShouldBeNil)
		})
	})
}
