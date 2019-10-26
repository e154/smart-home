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
// +----------+     +----------+    +----------+
// | handler  |     |  task    |    |  emitter |
// | script16 +-----> script17 +----> script18 |
// |          |     |          |    |          |
// +----------+     +----------+    +----------+
//
// add flow (flow2)
// +----------+     +----------+    +----------+
// | handler  |     | flowlink |    |  emitter |
// | script19 +-----> script20 +-X--> script21 |
// |          |     |          |    |          |
// +----------+     +----------+    +----------+
//
// send message flow2 to flow1
// +-----------+
// |   flow2   |
// |-----------|
// |           |
// |           +-----X----+
// |           |          |
// +-----------+          |
//                  +-----v-----+
//                  |   flow1   |
//                  |-----------|
//                  |           |
//                  |           |
//                  |           |
//                  +-----------+
//
// reset flow process after handler
//
func Test9(t *testing.T) {

	var cancelFunc context.CancelFunc
	var story = make([]string, 0)
	var scriptCounter string

	store = func(i interface{}) {
		cmd := fmt.Sprintf("%v", i)

		story = append(story, cmd)

		if cmd == "script20" {
			cancelFunc()
		}
	}

	store2 = func(i interface{}) {
		scriptCounter = fmt.Sprintf("%v", i)
	}

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

			// add device
			// ------------------------------------------------
			script16 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test16",
				Source:      coffeeScript16,
				Description: "test16",
			}

			engine16, err := scriptService.NewEngine(script16)
			So(err, ShouldBeNil)
			err = engine16.Compile()
			So(err, ShouldBeNil)
			script16Id, err := adaptors.Script.Add(script16)
			So(err, ShouldBeNil)
			script16, err = adaptors.Script.GetById(script16Id)
			So(err, ShouldBeNil)

			storeRegisterCallback(scriptService)

			script17 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test17",
				Source:      coffeeScript17,
				Description: "test17",
			}

			engine17, err := scriptService.NewEngine(script17)
			So(err, ShouldBeNil)
			err = engine17.Compile()
			So(err, ShouldBeNil)
			script17Id, err := adaptors.Script.Add(script17)
			So(err, ShouldBeNil)
			script17, err = adaptors.Script.GetById(script17Id)
			So(err, ShouldBeNil)

			script18 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test18",
				Source:      coffeeScript18,
				Description: "test18",
			}

			engine18, err := scriptService.NewEngine(script18)
			So(err, ShouldBeNil)
			err = engine18.Compile()
			So(err, ShouldBeNil)
			script18Id, err := adaptors.Script.Add(script18)
			So(err, ShouldBeNil)
			script18, err = adaptors.Script.GetById(script18Id)
			So(err, ShouldBeNil)

			script19 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test19",
				Source:      coffeeScript19,
				Description: "test19",
			}

			engine19, err := scriptService.NewEngine(script19)
			So(err, ShouldBeNil)
			err = engine19.Compile()
			So(err, ShouldBeNil)
			script19Id, err := adaptors.Script.Add(script19)
			So(err, ShouldBeNil)
			script19, err = adaptors.Script.GetById(script19Id)
			So(err, ShouldBeNil)

			script20 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test20",
				Source:      coffeeScript20,
				Description: "test20",
			}

			engine20, err := scriptService.NewEngine(script20)
			So(err, ShouldBeNil)
			err = engine20.Compile()
			So(err, ShouldBeNil)
			script20Id, err := adaptors.Script.Add(script20)
			So(err, ShouldBeNil)
			script20, err = adaptors.Script.GetById(script20Id)
			So(err, ShouldBeNil)

			script21 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test21",
				Source:      coffeeScript21,
				Description: "test21",
			}

			engine21, err := scriptService.NewEngine(script21)
			So(err, ShouldBeNil)
			err = engine21.Compile()
			So(err, ShouldBeNil)
			script21Id, err := adaptors.Script.Add(script21)
			So(err, ShouldBeNil)
			script21, err = adaptors.Script.GetById(script21Id)
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

			// add flow (flow1)
			// +----------+     +----------+    +----------+
			// | handler  |     |  task    |    |  emitter |
			// | script16 +-----> script17 +----> script18 |
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
				ScriptId:      &script16Id,
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
				ScriptId:      &script18Id,
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
				ScriptId:      &script17Id,
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
			// | handler  |     | flowlink |    |  emitter |
			// | script19 +-----> script20 +----> script21 |
			// |          |     |          |    |          |
			// +----------+     +----------+    +----------+
			flow2 := &m.Flow{
				Name:               "flow2",
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
				Name:          "handler2",
				FlowId:        flow2.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &script19Id,
				GraphSettings: m.FlowElementGraphSettings{
					Position: m.FlowElementGraphSettingsPosition{
						Top:  180,
						Left: 180,
					},
				},
			}
			feEmitter2 := &m.FlowElement{
				Name:          "emitter2",
				FlowId:        flow2.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageEmitter,
				ScriptId:      &script21Id,
				GraphSettings: m.FlowElementGraphSettings{
					Position: m.FlowElementGraphSettingsPosition{
						Top:  180,
						Left: 560,
					},
				},
			}
			flowLink2 := &m.FlowElement{
				Name:          "flow",
				FlowId:        flow2.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeFlow,
				ScriptId:      &script20Id,
				FlowLink:      &flow1.Id,
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
			ok, _ = flowLink2.Valid()
			So(ok, ShouldEqual, true)

			feHandler2.Uuid, err = adaptors.FlowElement.Add(feHandler2)
			So(err, ShouldBeNil)
			feEmitter2.Uuid, err = adaptors.FlowElement.Add(feEmitter2)
			So(err, ShouldBeNil)
			flowLink2.Uuid, err = adaptors.FlowElement.Add(flowLink2)
			So(err, ShouldBeNil)

			connect3 := &m.Connection{
				Name:        "con3",
				ElementFrom: feHandler2.Uuid,
				ElementTo:   flowLink2.Uuid,
				FlowId:      flow2.Id,
				PointFrom:   1,
				PointTo:     10,
			}
			connect4 := &m.Connection{
				Name:        "con4",
				ElementFrom: flowLink2.Uuid,
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

			// get flow
			// ------------------------------------------------
			err = c.Run()
			So(err, ShouldBeNil)

			workflowCore, err := c.GetWorkflow(workflow.Id)
			So(err, ShouldBeNil)

			flowCore, err := workflowCore.GetFLow(flow2.Id)
			So(err, ShouldBeNil)

			message := core.NewMessage()
			message.SetVar("val", 1)

			// create context
			var ctx context.Context
			ctx, cancelFunc = context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
			ctx = context.WithValue(ctx, "msg", message)

			for i := 0; i < 5; i++ {
				// send message ...
				err = flowCore.NewMessage(ctx)
				So(err, ShouldNotBeNil)
			}

			So(len(story), ShouldEqual, 2)
			So(scriptCounter, ShouldEqual, "3")

			err = c.Stop()
			So(err, ShouldBeNil)
		})
	})
}
