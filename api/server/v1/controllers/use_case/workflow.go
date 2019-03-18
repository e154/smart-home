package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"errors"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
)

func AddWorkflow(workflow *m.Workflow,
	adaptors *adaptors.Adaptors,
	core *core.Core) (result *models.Workflow, errs []*validation.Error, err error) {

	// validation
	_, errs = workflow.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = adaptors.Workflow.Add(workflow); err != nil {
		return
	}

	workflow.Id = id

	// add workflow
	if workflow, err = adaptors.Workflow.GetById(workflow.Id); err != nil {
		return
	}

	result = &models.Workflow{}
	if err = common.Copy(&result, &workflow); err != nil {
		return
	}

	err = core.AddWorkflow(workflow)

	return
}

func GetWorkflowById(workflowId int64, adaptors *adaptors.Adaptors) (result *models.Workflow, err error) {

	var workflow *m.Workflow
	if workflow, err = adaptors.Workflow.GetById(workflowId); err != nil {
		return
	}

	result = &models.Workflow{}
	err = common.Copy(&result, &workflow)

	return
}

func UpdateWorkflow(params *models.UpdateWorkflow,
	adaptors *adaptors.Adaptors,
	core *core.Core) (result *models.Workflow, errs []*validation.Error, err error) {

	var workflow *m.Workflow
	if workflow, err = adaptors.Workflow.GetById(params.Id); err != nil {
		return
	}

	if err = common.Copy(&workflow, &params); err != nil {
		return
	}

	// validation
	_, errs = workflow.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.Workflow.Update(workflow); err != nil {
		return
	}

	// reload workflow
	if workflow, err = adaptors.Workflow.GetById(workflow.Id); err != nil {
		return
	}

	result = &models.Workflow{}
	if err = common.Copy(&result, &workflow); err != nil {
		return
	}

	err = core.UpdateWorkflow(workflow)

	return
}

func GetWorkflowList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (result []*models.Workflow, total int64, err error) {

	var items []*m.Workflow
	if items, total, err = adaptors.Workflow.List(limit, offset, order, sortBy); err != nil {
		return
	}

	result = make([]*models.Workflow, 0)
	err = common.Copy(&result, &items)

	return
}

func DeleteWorkflowById(workflowId int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	if workflowId == 0 {
		err = errors.New("workflow id is null")
		return
	}

	var workflow *m.Workflow
	if workflow, err = adaptors.Workflow.GetById(workflowId); err != nil {
		return
	}

	// update core
	if err = core.DeleteWorkflow(workflow); err != nil {
		return
	}

	err = adaptors.Workflow.Delete(workflow.Id)

	return
}

func SearchWorkflow(query string, limit, offset int, adaptors *adaptors.Adaptors) (result []*models.Workflow, total int64, err error) {

	var items []*m.Workflow
	if items, total, err = adaptors.Workflow.Search(query, limit, offset); err != nil {
		return
	}

	result = make([]*models.Workflow, 0)
	err = common.Copy(&result, &items)

	return
}

func UpdateWorkflowScenario(workflowId int64, workflowScenarioId int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	var workflow *m.Workflow
	workflow, err = adaptors.Workflow.GetById(workflowId)
	if err != nil {
		return
	}

	if err = adaptors.Workflow.SetScenario(workflow, workflowScenarioId); err != nil {
		return
	}

	// update core
	if workflow, err = adaptors.Workflow.GetById(workflow.Id); err != nil {
		return
	}
	core.UpdateWorkflowScenario(workflow)

	return
}
