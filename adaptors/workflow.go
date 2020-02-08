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
	"errors"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
	"go/types"
)

type Workflow struct {
	table *db.Workflows
	db    *gorm.DB
}

func GetWorkflowAdaptor(d *gorm.DB) *Workflow {
	return &Workflow{
		table: &db.Workflows{Db: d},
		db:    d,
	}
}

func (n *Workflow) Add(workflow *m.Workflow) (id int64, err error) {

	dbWorkflow := n.toDb(workflow)
	if id, err = n.table.Add(dbWorkflow); err != nil {
		return
	}

	return
}

func (n *Workflow) GetAllEnabled() (list []*m.Workflow, err error) {

	var dbList []*db.Workflow
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.Workflow, 0)
	for _, dbWorkflow := range dbList {
		workflow := n.fromDb(dbWorkflow)
		list = append(list, workflow)
	}

	return
}

func (n *Workflow) GetById(workflowId int64) (workflow *m.Workflow, err error) {

	var dbWorkflow *db.Workflow
	if dbWorkflow, err = n.table.GetById(workflowId); err != nil {
		return
	}

	workflow = n.fromDb(dbWorkflow)

	return
}

func (n *Workflow) GetByWorkflowScenarioId(workflowScenarioId int64) (workflow *m.Workflow, err error) {

	var dbWorkflow *db.Workflow
	if dbWorkflow, err = n.table.GetByWorkflowScenarioId(workflowScenarioId); err != nil {
		return
	}

	workflow = n.fromDb(dbWorkflow)

	return
}

func (n *Workflow) Update(workflow *m.Workflow) (err error) {
	dbWorkflow := n.toDb(workflow)
	if err = n.table.Update(dbWorkflow); err != nil {
		return
	}

	//err = n.UpdateScripts(workflow)

	return
}

func (n *Workflow) Delete(workflowId int64) (err error) {
	err = n.table.Delete(workflowId)
	return
}

func (n *Workflow) List(limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*m.Workflow, total int64, err error) {
	var dbList []*db.Workflow
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort, onlyEnabled); err != nil {
		return
	}

	list = make([]*m.Workflow, 0)
	for _, dbWorkflow := range dbList {
		workflow := n.fromDb(dbWorkflow)
		list = append(list, workflow)
	}

	return
}

func (n *Workflow) DependencyLoading(workflow *m.Workflow) (err error) {
	dbWorkflow := n.toDb(workflow)
	err = n.table.DependencyLoading(dbWorkflow)
	return
}

func (n *Workflow) AddScript(workflow *m.Workflow, script *m.Script) (err error) {
	err = n.table.AddScript(workflow.Id, script.Id)
	return
}

func (n *Workflow) RemoveScript(workflow *m.Workflow, script *m.Script) (err error) {
	err = n.table.RemoveScript(workflow.Id, script.Id)
	return
}

func (n *Workflow) UpdateScripts(wf *m.Workflow) (err error) {

	var dbWf *db.Workflow
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

func (n *Workflow) SetScenario(workflow *m.Workflow, s interface{}) (err error) {
	var scenarioId *int64
	switch x := s.(type) {
	case int64:
		scenarioId = &x
	case *m.WorkflowScenario:
		scenarioId = &x.Id
	case types.Nil:
		scenarioId = nil
	default:
		err = errors.New("unknown scenario type")
		return
	}
	err = n.table.SetScenario(workflow.Id, scenarioId)
	return
}

func (n *Workflow) Search(query string, limit, offset int) (list []*m.Workflow, total int64, err error) {
	var dbList []*db.Workflow
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Workflow, 0)
	for _, dbWorkflow := range dbList {
		ver := n.fromDb(dbWorkflow)
		list = append(list, ver)
	}

	return
}

func (n *Workflow) fromDb(dbWorkflow *db.Workflow) (workflow *m.Workflow) {
	workflow = &m.Workflow{
		Id:          dbWorkflow.Id,
		Name:        dbWorkflow.Name,
		Description: dbWorkflow.Description,
		Status:      dbWorkflow.Status,
		CreatedAt:   dbWorkflow.CreatedAt,
		UpdatedAt:   dbWorkflow.UpdatedAt,
	}

	// scripts
	workflow.Scripts = make([]*m.Script, 0)
	scriptAdaptor := GetScriptAdaptor(n.db)
	for _, dbScript := range dbWorkflow.Scripts {
		script, _ := scriptAdaptor.fromDb(dbScript)
		workflow.Scripts = append(workflow.Scripts, script)
	}

	// scenario
	scenarioAdaptor := GetWorkflowScenarioAdaptor(n.db)
	if dbWorkflow.WorkflowScenarioId != nil {
		if dbWorkflow.WorkflowScenario == nil {
			for _, dbScenario := range dbWorkflow.Scenarios {
				if dbScenario.Id == *dbWorkflow.WorkflowScenarioId {
					workflow.Scenario = scenarioAdaptor.fromDb(dbScenario)
					break
				}
			}
		} else {
			workflow.Scenario = scenarioAdaptor.fromDb(dbWorkflow.WorkflowScenario)
		}
	}

	// scenarios
	workflow.Scenarios = make([]*m.WorkflowScenario, 0)
	for _, dbScenario := range dbWorkflow.Scenarios {
		scenario := scenarioAdaptor.fromDb(dbScenario)
		workflow.Scenarios = append(workflow.Scenarios, scenario)
	}

	return
}

func (n *Workflow) toDb(workflow *m.Workflow) (dbWorkflow *db.Workflow) {
	dbWorkflow = &db.Workflow{
		Id:          workflow.Id,
		Name:        workflow.Name,
		Description: workflow.Description,
		Status:      workflow.Status,
	}

	if workflow.Scenario != nil && workflow.Scenario.Id != 0 {
		dbWorkflow.WorkflowScenarioId = &workflow.Scenario.Id
	}

	return
}
