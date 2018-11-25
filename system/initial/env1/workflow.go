package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	. "github.com/e154/smart-home/common"
)

func addWorkflow(adaptors *adaptors.Adaptors,
	deviceAction1 *m.DeviceAction,
	script4 *m.Script) (workflow1 *m.Workflow) {

	workflow1 = &m.Workflow{
		Name:        "workflow1",
		Description: "workflow1 desc",
		Status:      "enabled",
	}
	ok, _ := workflow1.Valid()
	So(ok, ShouldEqual, true)

	var err error
	workflow1.Id, err = adaptors.Workflow.Add(workflow1)
	So(err, ShouldBeNil)

	// add workflow scenario
	// ------------------------------------------------
	wfScenario1 := &m.WorkflowScenario{
		Name:       "workflow1 scenario",
		SystemName: "wf_scenario_1",
		WorkflowId: workflow1.Id,
	}
	ok, _ = wfScenario1.Valid()
	So(ok, ShouldEqual, true)

	wfScenario1.Id, err = adaptors.WorkflowScenario.Add(wfScenario1)
	So(err, ShouldBeNil)

	// add flow1
	// ------------------------------------------------
	flow1 := &m.Flow{
		Name:               "flow1",
		Status:             Enabled,
		WorkflowId:         workflow1.Id,
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
	}
	feEmitter := &m.FlowElement{
		Name:          "emitter",
		FlowId:        flow1.Id,
		Status:        Enabled,
		PrototypeType: FlowElementsPrototypeMessageEmitter,
	}
	feTask1 := &m.FlowElement{
		Name:          "task",
		FlowId:        flow1.Id,
		Status:        Enabled,
		PrototypeType: FlowElementsPrototypeTask,
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
		PointTo:     1,
	}
	connect2 := &m.Connection{
		Name:        "con2",
		ElementFrom: feTask1.Uuid,
		ElementTo:   feEmitter.Uuid,
		FlowId:      flow1.Id,
		PointFrom:   1,
		PointTo:     1,
	}

	ok, _ = connect1.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = connect2.Valid()
	So(ok, ShouldEqual, true)

	connect1.Uuid, err = adaptors.Connection.Add(connect1)
	So(err, ShouldBeNil)
	connect2.Uuid, err = adaptors.Connection.Add(connect2)
	So(err, ShouldBeNil)

	// add worker
	worker := &m.Worker{
		Name:           "worker",
		Time:           "* * * * * *",
		Status:         "enabled",
		WorkflowId:     workflow1.Id,
		FlowId:         flow1.Id,
		DeviceActionId: deviceAction1.Id,
	}

	ok, _ = worker.Valid()
	So(ok, ShouldEqual, true)

	worker.Id, err = adaptors.Worker.Add(worker)
	So(err, ShouldBeNil)

	return
}
