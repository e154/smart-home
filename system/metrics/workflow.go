// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package metrics

import (
	"sync"

	"github.com/e154/smart-home/common"

	"github.com/e154/smart-home/adaptors"
	"github.com/rcrowley/go-metrics"
)

// WorkflowStatus ...
type WorkflowStatus struct {
	Id         int64 `json:"id"`
	ScenarioId int64 `json:"scenario_id"`
}

// Workflow ...
type Workflow struct {
	Total    int64                    `json:"total"`
	Disabled int64                    `json:"disabled"`
	Status   map[int64]WorkflowStatus `json:"status"`
}

// WorkflowManager ...
type WorkflowManager struct {
	publisher  IPublisher
	enabled    metrics.Counter
	total      metrics.Counter
	updateLock sync.Mutex
	status     map[int64]int64
}

// NewWorkflowManager ...
func NewWorkflowManager(publisher IPublisher,
	adaptors *adaptors.Adaptors) (wf *WorkflowManager) {

	wf = &WorkflowManager{
		total:     metrics.NewCounter(),
		enabled:   metrics.NewCounter(),
		publisher: publisher,
		status:    make(map[int64]int64),
	}

	//if _, total, err := adaptors.Workflow.List(999, 0, "", "", false); err == nil {
	//	wf.total.Inc(total)
	//}

	return
}

// GetStatus ...
func (d *WorkflowManager) GetStatus(workflowId int64) (status WorkflowStatus, err error) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	if scenarioId, ok := d.status[workflowId]; ok {
		status.ScenarioId = scenarioId
		status.Id = workflowId
		return
	}

	err = common.ErrNotFound

	return
}

func (d *WorkflowManager) update(t interface{}) {
	switch v := t.(type) {
	case WorkflowUpdateScenario:
		d.updateLock.Lock()
		d.status[v.Id] = v.ScenarioId
		d.updateLock.Unlock()
	case WorkflowAdd:
		d.total.Inc(v.TotalNum)
		d.enabled.Inc(v.EnabledNum)
	case WorkflowDelete:
		d.total.Dec(v.TotalNum)
		d.enabled.Dec(v.EnabledNum)
	default:
		return
	}

	d.broadcast()
}

// Snapshot ...
func (d *WorkflowManager) Snapshot() Workflow {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	var status = make(map[int64]WorkflowStatus)
	for id, scenId := range d.status {
		status[id] = WorkflowStatus{
			Id:         id,
			ScenarioId: scenId,
		}
	}

	return Workflow{
		Total:    d.enabled.Count(),
		Disabled: d.total.Count() - d.enabled.Count(),
		Status:   status,
	}
}

func (d *WorkflowManager) broadcast() {
	go d.publisher.Broadcast("workflow")
}

// WorkflowUpdateScenario ...
type WorkflowUpdateScenario struct {
	Id         int64
	ScenarioId int64
}

// WorkflowAdd ...
type WorkflowAdd struct {
	TotalNum   int64
	EnabledNum int64
}

// WorkflowDelete ...
type WorkflowDelete struct {
	TotalNum   int64
	EnabledNum int64
}
