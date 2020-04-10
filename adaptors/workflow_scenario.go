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

// WorkflowScenario ...
type WorkflowScenario struct {
	table *db.WorkflowScenarios
	db    *gorm.DB
}

// GetWorkflowScenarioAdaptor ...
func GetWorkflowScenarioAdaptor(d *gorm.DB) *WorkflowScenario {
	return &WorkflowScenario{
		table: &db.WorkflowScenarios{Db: d},
		db:    d,
	}
}

// Add ...
func (n *WorkflowScenario) Add(workflow *m.WorkflowScenario) (id int64, err error) {

	dbWorkflowScenario := n.toDb(workflow)
	if id, err = n.table.Add(dbWorkflowScenario); err != nil {
		return
	}

	return
}

// GetById ...
func (n *WorkflowScenario) GetById(scenarioId int64) (workflow *m.WorkflowScenario, err error) {

	var dbWorkflowScenario *db.WorkflowScenario
	if dbWorkflowScenario, err = n.table.GetById(scenarioId); err != nil {
		return
	}

	workflow = n.fromDb(dbWorkflowScenario)

	return
}

// Update ...
func (n *WorkflowScenario) Update(workflow *m.WorkflowScenario) (err error) {
	dbWorkflowScenario := n.toDb(workflow)
	if err = n.table.Update(dbWorkflowScenario); err != nil {
		return
	}

	err = n.UpdateScripts(workflow)

	return
}

// Delete ...
func (n *WorkflowScenario) Delete(workflowId int64) (err error) {
	err = n.table.Delete(workflowId)
	return
}

// List ...
func (n *WorkflowScenario) List(limit, offset int64, orderBy, sort string) (list []*m.WorkflowScenario, total int64, err error) {
	var dbList []*db.WorkflowScenario
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.WorkflowScenario, 0)
	for _, dbWorkflowScenario := range dbList {
		workflow := n.fromDb(dbWorkflowScenario)
		list = append(list, workflow)
	}

	return
}

// ListByWorkflow ...
func (n *WorkflowScenario) ListByWorkflow(workflowId int64) (list []*m.WorkflowScenario, total int64, err error) {
	var dbList []*db.WorkflowScenario
	if dbList, total, err = n.table.ListByWorkflow(workflowId); err != nil {
		return
	}

	list = make([]*m.WorkflowScenario, 0)
	for _, dbWorkflowScenario := range dbList {
		workflow := n.fromDb(dbWorkflowScenario)
		list = append(list, workflow)
	}

	return
}

// AddScript ...
func (n *WorkflowScenario) AddScript(workflowScenario *m.WorkflowScenario, script *m.Script) (err error) {
	err = n.table.AddScript(workflowScenario.Id, script.Id)
	return
}

// RemoveScript ...
func (n *WorkflowScenario) RemoveScript(workflowScenario *m.WorkflowScenario, script *m.Script) (err error) {
	err = n.table.RemoveScript(workflowScenario.Id, script.Id)
	return
}

// UpdateScripts ...
func (n *WorkflowScenario) UpdateScripts(wf *m.WorkflowScenario) (err error) {

	var dbWf *db.WorkflowScenario
	dbWf, err = n.table.GetById(wf.Id)
	if err != nil {
		return
	}

	var exist bool
	for _, s1 := range wf.Scripts {
		exist = false
		for _, s2 := range dbWf.Scripts {
			if s1.Id == s2.Id {
				exist = true
				break
			}
		}
		if !exist {
			n.AddScript(wf, s1)
		}
	}

	for _, s1 := range dbWf.Scripts {
		exist = false
		for _, s2 := range wf.Scripts {
			if s1.Id == s2.Id {
				exist = true
				break
			}
		}
		if !exist {
			n.RemoveScript(wf, &m.Script{Id: s1.Id})

		}
	}

	return
}

// Search ...
func (n *WorkflowScenario) Search(query string, workflowId, limit, offset int) (list []*m.WorkflowScenario, total int64, err error) {
	var dbList []*db.WorkflowScenario
	if dbList, total, err = n.table.Search(query, workflowId, limit, offset); err != nil {
		return
	}

	list = make([]*m.WorkflowScenario, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *WorkflowScenario) fromDb(dbWorkflowScenario *db.WorkflowScenario) (workflow *m.WorkflowScenario) {
	workflow = &m.WorkflowScenario{
		Id:         dbWorkflowScenario.Id,
		Name:       dbWorkflowScenario.Name,
		WorkflowId: dbWorkflowScenario.WorkflowId,
		SystemName: dbWorkflowScenario.SystemName,
		CreatedAt:  dbWorkflowScenario.CreatedAt,
		UpdatedAt:  dbWorkflowScenario.UpdatedAt,
	}

	// scripts
	workflow.Scripts = make([]*m.Script, 0)
	scriptAdaptor := GetScriptAdaptor(n.db)
	for _, dbScript := range dbWorkflowScenario.Scripts {
		script, _ := scriptAdaptor.fromDb(dbScript)
		workflow.Scripts = append(workflow.Scripts, script)
	}

	return
}

func (n *WorkflowScenario) toDb(workflow *m.WorkflowScenario) (dbWorkflowScenario *db.WorkflowScenario) {
	dbWorkflowScenario = &db.WorkflowScenario{
		Id:         workflow.Id,
		Name:       workflow.Name,
		WorkflowId: workflow.WorkflowId,
		SystemName: workflow.SystemName,
	}
	return
}
