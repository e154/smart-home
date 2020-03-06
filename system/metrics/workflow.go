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

package metrics

import (
	"github.com/rcrowley/go-metrics"
	"sync"
)

type WorkflowStatus struct {
	Id         int64 `json:"id"`
	ScenarioId int64 `json:"scenario_id"`
}

type Workflow struct {
	Total  int64                    `json:"total"`
	Status map[int64]WorkflowStatus `json:"status"`
}

type WorkflowManager struct {
	publisher  IPublisher
	total      metrics.Counter
	updateLock sync.Mutex
	status     map[int64]int64
}

func NewWorkflowManager(publisher IPublisher) (wf *WorkflowManager) {
	wf = &WorkflowManager{
		total:     metrics.NewCounter(),
		publisher: publisher,
		status:    make(map[int64]int64),
	}

	return
}

func (d *WorkflowManager) GetStatus(workflowId int64) (status WorkflowStatus, err error) {

	if scenarioId, ok := d.status[workflowId]; ok {
		status.ScenarioId = scenarioId
		status.Id = workflowId
		return
	}

	err = ErrRecordNotFound

	return
}

func (d *WorkflowManager) update(t interface{}) {
	switch v := t.(type) {
	case WorkflowUpdateScenario:
		d.updateLock.Lock()
		defer d.updateLock.Unlock()
		d.status[v.Id] = v.ScenarioId
	case WorkflowAdd:
		d.total.Inc(v.Num)
	case WorkflowDelete:
		d.total.Dec(v.Num)
	default:
		return
	}

	d.broadcast()
}

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
		Total:  d.total.Count(),
		Status: status,
	}
}

func (d *WorkflowManager) broadcast() {
	go d.publisher.Broadcast("workflow")
}

type WorkflowUpdateScenario struct {
	Id         int64
	ScenarioId int64
}

type WorkflowAdd struct {
	Num int64
}

type WorkflowDelete struct {
	Num int64
}
