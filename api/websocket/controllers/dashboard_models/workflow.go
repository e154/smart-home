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

package dashboard_models

import (
	metrics2 "github.com/e154/smart-home/system/metrics"
)

type Workflow struct {
	metrics *metrics2.MetricManager
}

func NewWorkflow(metrics *metrics2.MetricManager) *Workflow {
	return &Workflow{
		metrics: metrics,
	}
}

// status all workflows
func (w *Workflow) Broadcast() (map[string]interface{}, bool) {

	return map[string]interface{}{
		//"workflows": w.metrics.Workflow.Snapshot(),
	}, true
}

//func (w *Workflow) BroadcastOne(params telemetry.WorkflowScenario) (map[string]interface{}, bool) {

	//w.Status[params.WorkflowId] = m.DashboardWorkflowStatus{
	//	Id:         params.WorkflowId,
	//	ScenarioId: params.ScenarioId,
	//}
	//
	//w.Total = len(w.Status)
	//
	//return map[string]interface{}{
	//	"workflow":  map[string]interface{}{"id": params.WorkflowId, "scenario_id": params.ScenarioId},
	//	"workflows": w,
	//}, true
//}
