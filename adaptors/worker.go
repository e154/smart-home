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

package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

// Worker ...
type Worker struct {
	table *db.Workers
	db    *gorm.DB
}

// GetWorkerAdaptor ...
func GetWorkerAdaptor(d *gorm.DB) *Worker {
	return &Worker{
		table: &db.Workers{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Worker) Add(worker *m.Worker) (id int64, err error) {

	dbWorker := n.toDb(worker)
	if id, err = n.table.Add(dbWorker); err != nil {
		return
	}

	return
}

// GetAllEnabled ...
func (n *Worker) GetAllEnabled() (list []*m.Worker, err error) {

	var dbList []*db.Worker
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.Worker, 0)
	for _, dbWorker := range dbList {
		worker := n.fromDb(dbWorker)
		list = append(list, worker)
	}

	return
}

// GetById ...
func (n *Worker) GetById(workerId int64) (worker *m.Worker, err error) {

	var dbWorker *db.Worker
	if dbWorker, err = n.table.GetById(workerId); err != nil {
		return
	}

	worker = n.fromDb(dbWorker)

	return
}

// Update ...
func (n *Worker) Update(worker *m.Worker) (err error) {
	dbWorker := n.toDb(worker)
	err = n.table.Update(dbWorker)
	return
}

// Delete ...
func (n *Worker) Delete(ids []int64) (err error) {
	err = n.table.Delete(ids)
	return
}

// List ...
func (n *Worker) List(limit, offset int64, orderBy, sort string) (list []*m.Worker, total int64, err error) {
	var dbList []*db.Worker
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Worker, 0)
	for _, dbWorker := range dbList {
		worker := n.fromDb(dbWorker)
		list = append(list, worker)
	}

	return
}

func (n *Worker) fromDb(dbWorker *db.Worker) (worker *m.Worker) {
	worker = &m.Worker{
		Id:             dbWorker.Id,
		WorkflowId:     dbWorker.WorkflowId,
		DeviceActionId: dbWorker.DeviceActionId,
		FlowId:         dbWorker.FlowId,
		Status:         dbWorker.Status,
		Name:           dbWorker.Name,
		Time:           dbWorker.Time,
		CreatedAt:      dbWorker.CreatedAt,
		UpdatedAt:      dbWorker.UpdatedAt,
	}

	// workflow
	if dbWorker.Workflow != nil {
		workflowAdaptor := GetWorkflowAdaptor(n.db)
		worker.Workflow = workflowAdaptor.fromDb(dbWorker.Workflow)
	}

	// deviceAction
	if dbWorker.DeviceAction != nil {
		deviceActionAdaptor := GetDeviceActionAdaptor(n.db)
		worker.DeviceAction = deviceActionAdaptor.fromDb(dbWorker.DeviceAction)
	}

	// flow
	if dbWorker.Flow != nil {
		flowAdaptor := GetFlowAdaptor(n.db)
		worker.Flow = flowAdaptor.fromDb(dbWorker.Flow)
	}

	return
}

func (n *Worker) toDb(worker *m.Worker) (dbWorker *db.Worker) {
	dbWorker = &db.Worker{
		Id:             worker.Id,
		WorkflowId:     worker.WorkflowId,
		DeviceActionId: worker.DeviceActionId,
		FlowId:         worker.FlowId,
		Status:         worker.Status,
		Name:           worker.Name,
		Time:           worker.Time,
	}
	return
}
