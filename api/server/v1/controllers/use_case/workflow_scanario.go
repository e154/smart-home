package use_case

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"errors"
)

func GetWorkflowScenarioById(workflowId, scenarioId int64, adaptors *adaptors.Adaptors) (workflowScenario *m.WorkflowScenario, err error) {

	if workflowScenario, err = adaptors.WorkflowScenario.GetById(scenarioId); err != nil {
		return
	}

	if workflowScenario.WorkflowId != workflowId {
		err = fmt.Errorf("record not found")
	}

	return
}

func AddWorkflowScenario(workflowScenario *m.WorkflowScenario, adaptors *adaptors.Adaptors) (ok bool, id int64, errs []*validation.Error, err error) {

	// validation
	ok, errs = workflowScenario.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if id, err = adaptors.WorkflowScenario.Add(workflowScenario); err != nil {
		return
	}

	workflowScenario.Id = id

	return
}

func GetWorkflowScenarioList(workflowId int64, adaptors *adaptors.Adaptors) (items []*m.WorkflowScenario, total int64, err error) {

	items, total, err = adaptors.WorkflowScenario.ListByWorkflow(workflowId)

	return
}

func DeleteWorkflowScenarioById(workflowScenarioId int64, adaptors *adaptors.Adaptors) (err error) {

	if workflowScenarioId == 0 {
		err = errors.New("scenario id is null")
		return
	}

	err = adaptors.WorkflowScenario.Delete(workflowScenarioId)

	return
}

func SearchWorkflowScenario(query string, limit, offset int, adaptors *adaptors.Adaptors) (scenarios []*m.WorkflowScenario, total int64, err error) {

	scenarios, total, err = adaptors.WorkflowScenario.Search(query, limit, offset)

	return
}

func WorkflowUpdateWorkflowScenario(workflowScenario *m.WorkflowScenario, adaptors *adaptors.Adaptors)(ok bool, errs []*validation.Error, err error) {

	// validation
	ok, errs = workflowScenario.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if err = adaptors.WorkflowScenario.Update(workflowScenario); err != nil {
		return
	}

	return
}