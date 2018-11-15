package workflow

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/common"
	"fmt"
)

// create flow
//
//  workflow + workflow scenario
//
//  emitter --> handler
//
func Test2(t *testing.T) {

	var script1 *m.Script
	Convey("add scripts", t, func(ctx C) {
		container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			// clear database
			migrations.Purge()

			// create scripts
			script1 = &m.Script{
				Lang:        "coffeescript",
				Name:        "test1",
				Source:      coffeeScript1,
				Description: "test1",
			}

			ok, _ := script1.Valid()
			So(ok, ShouldEqual, true)

			engine1, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)
			err = engine1.Compile()
			So(err, ShouldBeNil)
			script1Id, err := adaptors.Script.Add(script1)
			So(err, ShouldBeNil)
			script1, err = adaptors.Script.GetById(script1Id)
			So(err, ShouldBeNil)
		})
	})

	var workflow *m.Workflow
	var wfScenario1 *m.WorkflowScenario
	Convey("add workflow", t, func(ctx C) {
		container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			// create workflow
			workflow = &m.Workflow{
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
			wfScenario1 = &m.WorkflowScenario{
				Name:       "wf scenario 1",
				SystemName: "wf_scenario_1",
				WorkflowId: workflow.Id,
			}

			ok, _ = wfScenario1.Valid()
			So(ok, ShouldEqual, true)

			wfScenarioId1, err := adaptors.WorkflowScenario.Add(wfScenario1)
			So(err, ShouldBeNil)
			wfScenario1.Id = wfScenarioId1

			workflow.Scenario = wfScenario1
			err = adaptors.Workflow.Update(workflow)
			So(err, ShouldBeNil)
		})
	})

	Convey("add flow", t, func(ctx C) {
		container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			flow1 := &m.Flow{
				Name:               "flow1",
				Status:             Enabled,
				WorkflowId:         workflow.Id,
				WorkflowScenarioId: wfScenario1.Id,
			}

			ok, _ := flow1.Valid()
			So(ok, ShouldEqual, true)

			var err error
			flow1.Id, err = adaptors.Flow.Add(flow1)
			So(err, ShouldBeNil)

			feEmitter := &m.FlowElement{
				Name:          "emitter",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageEmitter,
			}

			feHandler := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
			}

			ok, _ = feEmitter.Valid()
			So(ok, ShouldEqual, true)

			ok, _ = feHandler.Valid()
			So(ok, ShouldEqual, true)

			feEmitter.Uuid, err = adaptors.FlowElement.Add(feEmitter)
			So(err, ShouldBeNil)

			feHandler.Uuid, err = adaptors.FlowElement.Add(feHandler)
			So(err, ShouldBeNil)

			connect := &m.Connection{
				Name:        "con1",
				ElementFrom: feEmitter.Uuid,
				ElementTo:   feHandler.Uuid,
				FlowId:      flow1.Id,
				PointFrom:   1,
				PointTo:     1,
			}

			ok, errs := connect.Valid()
			for _, err := range errs {
				fmt.Println(err.Name, err.String())
			}
			So(ok, ShouldEqual, true)

			connect.Uuid, err = adaptors.Connection.Add(connect)
			So(err, ShouldBeNil)
		})
	})
}
