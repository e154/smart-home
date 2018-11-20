package workflow

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/core"
)

//
// create workflow
// 				(script1, script2)
//
// add workflow scenarios
// 				(wf_scenario_1 + script3, wf_scenario_2 + script4)
//
// run workflow
//
func Test1(t *testing.T) {

	var script1, script2, script3, script4 *m.Script
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
			script2 = &m.Script{
				Lang:        "coffeescript",
				Name:        "test2",
				Source:      coffeeScript2,
				Description: "test2",
			}
			script3 = &m.Script{
				Lang:        "coffeescript",
				Name:        "test3",
				Source:      coffeeScript3,
				Description: "test3",
			}
			script4 = &m.Script{
				Lang:        "coffeescript",
				Name:        "test4",
				Source:      coffeeScript4,
				Description: "test4",
			}

			engine1, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)
			err = engine1.Compile()
			So(err, ShouldBeNil)
			script1Id, err := adaptors.Script.Add(script1)
			So(err, ShouldBeNil)
			script1, err = adaptors.Script.GetById(script1Id)
			So(err, ShouldBeNil)

			engine2, err := scriptService.NewEngine(script2)
			So(err, ShouldBeNil)
			err = engine2.Compile()
			So(err, ShouldBeNil)
			script2Id, err := adaptors.Script.Add(script2)
			So(err, ShouldBeNil)
			script2, err = adaptors.Script.GetById(script2Id)
			So(err, ShouldBeNil)

			engine3, err := scriptService.NewEngine(script3)
			So(err, ShouldBeNil)
			err = engine3.Compile()
			So(err, ShouldBeNil)
			script3Id, err := adaptors.Script.Add(script3)
			So(err, ShouldBeNil)
			script3, err = adaptors.Script.GetById(script3Id)
			So(err, ShouldBeNil)

			engine4, err := scriptService.NewEngine(script4)
			So(err, ShouldBeNil)
			err = engine4.Compile()
			So(err, ShouldBeNil)
			script4Id, err := adaptors.Script.Add(script4)
			So(err, ShouldBeNil)
			script4, err = adaptors.Script.GetById(script4Id)
			So(err, ShouldBeNil)
		})
	})

	var workflow *m.Workflow
	Convey("add workflow", t, func(ctx C) {

		container.Invoke(func(adaptors *adaptors.Adaptors) {


			// create workflow
			workflow = &m.Workflow{
				Name: "main workflow",
				Description: "main workflow desc",
				Status: "enabled",
			}

			wfId, err := adaptors.Workflow.Add(workflow)
			So(err, ShouldBeNil)
			workflow.Id = wfId

			err = adaptors.Workflow.AddScript(workflow, script1)
			So(err, ShouldBeNil)

			err = adaptors.Workflow.AddScript(workflow, script2)
			So(err, ShouldBeNil)
		})
	})

	var wfScenario1, wfScenario2 *m.WorkflowScenario
	Convey("add workflow scenarios", t, func(ctx C) {

		container.Invoke(func(adaptors *adaptors.Adaptors) {

			wfScenario1 = &m.WorkflowScenario{
				Name: "wf scenario 1",
				SystemName: "wf_scenario_1",
				WorkflowId: workflow.Id,
			}
			wfScenario2 = &m.WorkflowScenario{
				Name: "wf scenario 2",
				SystemName: "wf_scenario_2",
				WorkflowId: workflow.Id,
			}

			wfScenarioId1, err := adaptors.WorkflowScenario.Add(wfScenario1)
			So(err, ShouldBeNil)
			wfScenario1.Id = wfScenarioId1
			err = adaptors.WorkflowScenario.AddScript(wfScenario1, script3)
			So(err, ShouldBeNil)

			wfScenarioId2, err := adaptors.WorkflowScenario.Add(wfScenario2)
			So(err, ShouldBeNil)
			wfScenario2.Id = wfScenarioId2
			err = adaptors.WorkflowScenario.AddScript(wfScenario2, script4)
			So(err, ShouldBeNil)

			workflow.Scenario = wfScenario1
			err = adaptors.Workflow.Update(workflow)
			So(err, ShouldBeNil)
		})
	})

	Convey("check workflow", t, func(ctx C) {
		container.Invoke(func(adaptors *adaptors.Adaptors) {
			wf, err := adaptors.Workflow.GetById(workflow.Id)
			So(err, ShouldBeNil)

			So(wf.Id, ShouldEqual, workflow.Id)
			So(len(wf.Scenarios), ShouldEqual, 2)
			So(wf.Scenario, ShouldNotBeNil)
			So(wf.Scenario.Id, ShouldEqual, wfScenario1.Id)
			So(wf.Name, ShouldEqual, workflow.Name)
			So(wf.Description, ShouldEqual, workflow.Description)
			So(wf.Status, ShouldEqual, "enabled")
		})
	})

	Convey("run workflow", t, func(ctx C) {

		container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService *scripts.ScriptService) {

			var err error
			workflow, err = adaptors.Workflow.GetById(workflow.Id)
			So(err, ShouldBeNil)

			wf := core.NewWorkflow(workflow, adaptors, scriptService)
			err = wf.Run()
			So(err, ShouldBeNil)

			err = wf.Stop()
			So(err, ShouldBeNil)
		})
	})
}
