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

package endpoint

import (
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

type WorkflowEndpoint struct {
	*CommonEndpoint
}

func NewWorkflowEndpoint(common *CommonEndpoint) *WorkflowEndpoint {
	return &WorkflowEndpoint{
		CommonEndpoint: common,
	}
}

func (n *WorkflowEndpoint) Add(params *m.Workflow) (result *m.Workflow, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.Workflow.Add(params); err != nil {
		return
	}

	if result, err = n.adaptors.Workflow.GetById(id); err != nil {
		return
	}

	err = n.core.AddWorkflow(result)

	return
}

func (n *WorkflowEndpoint) GetById(workflowId int64) (result *m.Workflow, err error) {

	result, err = n.adaptors.Workflow.GetById(workflowId)

	return
}

func (n *WorkflowEndpoint) Update(params *m.Workflow,
) (result *m.Workflow, errs []*validation.Error, err error) {

	var workflow *m.Workflow
	if workflow, err = n.adaptors.Workflow.GetById(params.Id); err != nil {
		return
	}

	if err = common.Copy(&workflow, &params, common.JsonEngine); err != nil {
		return
	}

	// validation
	_, errs = workflow.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Workflow.Update(workflow); err != nil {
		return
	}

	if err = n.adaptors.Workflow.UpdateScripts(workflow); err != nil {
		return
	}

	// reload workflow
	if result, err = n.adaptors.Workflow.GetById(workflow.Id); err != nil {
		return
	}

	err = n.core.UpdateWorkflow(result)

	return
}

func (n *WorkflowEndpoint) GetList(limit, offset int64, order, sortBy string, onlyEnabled bool) (result []*m.Workflow, total int64, err error) {

	result, total, err = n.adaptors.Workflow.List(limit, offset, order, sortBy, onlyEnabled)

	return
}

func (n *WorkflowEndpoint) Delete(workflowId int64) (err error) {

	if workflowId == 0 {
		err = errors.New("workflow id is null")
		return
	}

	var workflow *m.Workflow
	if workflow, err = n.adaptors.Workflow.GetById(workflowId); err != nil {
		return
	}

	// update core
	if err = n.core.DeleteWorkflow(workflow); err != nil {
		return
	}

	err = n.adaptors.Workflow.Delete(workflow.Id)

	return
}

func (n *WorkflowEndpoint) Search(query string, limit, offset int) (result []*m.Workflow, total int64, err error) {

	result, total, err = n.adaptors.Workflow.Search(query, limit, offset)

	return
}

func (n *WorkflowEndpoint) UpdateScenario(workflowId int64, workflowScenarioId int64) (err error) {

	var workflow *m.Workflow
	workflow, err = n.adaptors.Workflow.GetById(workflowId)
	if err != nil {
		return
	}

	if err = n.adaptors.Workflow.SetScenario(workflow, workflowScenarioId); err != nil {
		return
	}

	err = n.core.UpdateWorkflowScenario(workflowId)

	return
}

func (n WorkflowEndpoint) AddScript(workflow *m.Workflow, script *m.Script) (err error) {
	if err = n.adaptors.Workflow.AddScript(workflow, script); err != nil {
		return
	}
	workflow.Scripts = append(workflow.Scripts, script)
	return
}
