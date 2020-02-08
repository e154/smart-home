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

package env1

import (
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

func addWorkflow(adaptors *adaptors.Adaptors,
	deviceActions map[string]*m.DeviceAction,
	scripts map[string]*m.Script) (workflow1 *m.Workflow) {

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

	err = adaptors.Workflow.AddScript(workflow1, scripts["wflow_script_v1"])
	So(err, ShouldBeNil)

	// add workflow scenario
	// ------------------------------------------------
	wfScenario1 := &m.WorkflowScenario{
		Name:       "Будний день(weekday)",
		SystemName: "weekday",
		WorkflowId: workflow1.Id,
	}
	wfScenario2 := &m.WorkflowScenario{
		Name:       "Выходные (weekend)",
		SystemName: "weekend",
		WorkflowId: workflow1.Id,
	}
	ok, _ = wfScenario1.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = wfScenario2.Valid()
	So(ok, ShouldEqual, true)

	wfScenario1.Id, err = adaptors.WorkflowScenario.Add(wfScenario1)
	So(err, ShouldBeNil)
	err = adaptors.WorkflowScenario.AddScript(wfScenario1, scripts["wflow_scenario_weekday_v1"])
	So(err, ShouldBeNil)

	wfScenario2.Id, err = adaptors.WorkflowScenario.Add(wfScenario2)
	So(err, ShouldBeNil)
	err = adaptors.WorkflowScenario.AddScript(wfScenario2, scripts["wflow_scenario_weekend_v1"])
	So(err, ShouldBeNil)

	err = adaptors.Workflow.SetScenario(workflow1, wfScenario1)
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
		ScriptId:      &scripts["base_script"].Id,
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

	// add worker
	worker := &m.Worker{
		Name:           "worker",
		Time:           "* * * * * *",
		Status:         "enabled",
		WorkflowId:     workflow1.Id,
		FlowId:         flow1.Id,
		DeviceActionId: deviceActions["mb_dev1_condition_check_v1"].Id,
	}

	ok, _ = worker.Valid()
	So(ok, ShouldEqual, true)

	worker.Id, err = adaptors.Worker.Add(worker)
	So(err, ShouldBeNil)

	// add command flow
	// ------------------------------------------------
	flow2 := &m.Flow{
		Name:               "flow2",
		Status:             Enabled,
		WorkflowId:         workflow1.Id,
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
		ScriptId:      &scripts["base_script"].Id,
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
		Name:        "con1",
		ElementFrom: feHandler2.Uuid,
		ElementTo:   feTask2.Uuid,
		FlowId:      flow2.Id,
		PointFrom:   1,
		PointTo:     10,
	}
	connect4 := &m.Connection{
		Name:        "con2",
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

	// add worker
	worker2 := &m.Worker{
		Name:           "worker2",
		Time:           "5,10,15,20,25,30,35,40,45,50,55 * * * * *",
		Status:         "enabled",
		WorkflowId:     workflow1.Id,
		FlowId:         flow2.Id,
		DeviceActionId: deviceActions["cmd_condition_check_v1"].Id,
	}

	ok, _ = worker2.Valid()
	So(ok, ShouldEqual, true)

	worker2.Id, err = adaptors.Worker.Add(worker2)
	So(err, ShouldBeNil)
	return
}
