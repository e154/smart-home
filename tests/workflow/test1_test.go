package workflow

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/core"
	cr "github.com/e154/smart-home/system/cron"
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

	var _scripts map[string]*m.Script
	Convey("add scripts", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			// clear database
			migrations.Purge()

			// create scripts
			_scripts = GetScripts(ctx, scriptService, adaptors, 1,2,3,4)
		})
	})

	var workflow *m.Workflow
	Convey("add workflow", t, func(ctx C) {

		_ = container.Invoke(func(adaptors *adaptors.Adaptors) {

			// create workflow
			workflow = &m.Workflow{
				Name:        "main workflow",
				Description: "main workflow desc",
				Status:      "enabled",
			}

			wfId, err := adaptors.Workflow.Add(workflow)
			So(err, ShouldBeNil)
			workflow.Id = wfId

			err = adaptors.Workflow.AddScript(workflow, _scripts["script1"])
			So(err, ShouldBeNil)

			err = adaptors.Workflow.AddScript(workflow, _scripts["script2"])
			So(err, ShouldBeNil)
		})
	})

	var wfScenario1, wfScenario2 *m.WorkflowScenario
	Convey("add workflow scenarios", t, func(ctx C) {

		_ = container.Invoke(func(adaptors *adaptors.Adaptors) {

			wfScenario1 = &m.WorkflowScenario{
				Name:       "wf scenario 1",
				SystemName: "wf_scenario_1",
				WorkflowId: workflow.Id,
			}
			wfScenario2 = &m.WorkflowScenario{
				Name:       "wf scenario 2",
				SystemName: "wf_scenario_2",
				WorkflowId: workflow.Id,
			}

			wfScenarioId1, err := adaptors.WorkflowScenario.Add(wfScenario1)
			So(err, ShouldBeNil)
			wfScenario1.Id = wfScenarioId1
			err = adaptors.WorkflowScenario.AddScript(wfScenario1, _scripts["script3"])
			So(err, ShouldBeNil)

			wfScenarioId2, err := adaptors.WorkflowScenario.Add(wfScenario2)
			So(err, ShouldBeNil)
			wfScenario2.Id = wfScenarioId2
			err = adaptors.WorkflowScenario.AddScript(wfScenario2, _scripts["script4"])
			So(err, ShouldBeNil)

			workflow.Scenario = wfScenario1
			err = adaptors.Workflow.Update(workflow)
			So(err, ShouldBeNil)
		})
	})

	Convey("check workflow", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors) {
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

		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService *scripts.ScriptService,
			cron *cr.Cron,
			c *core.Core) {

			var err error
			workflow, err = adaptors.Workflow.GetById(workflow.Id)
			So(err, ShouldBeNil)

			wf := core.NewWorkflow(workflow, adaptors, scriptService, cron, c)
			err = wf.Run()
			So(err, ShouldBeNil)

			err = wf.Stop()
			So(err, ShouldBeNil)
		})
	})
}
