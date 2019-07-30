package endpoint

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
	m "github.com/e154/smart-home/models"
	"errors"
	"fmt"
)

type WorkflowScenarioEndpoint struct {
	*CommonEndpoint
}

func NewWorkflowScenarioEndpoint(common *CommonEndpoint) *WorkflowScenarioEndpoint {
	return &WorkflowScenarioEndpoint{
		CommonEndpoint: common,
	}
}

func (n *WorkflowScenarioEndpoint) GetById(workflowId, scenarioId int64) (result *m.WorkflowScenario, err error) {

	result, err = n.adaptors.WorkflowScenario.GetById(scenarioId)
	
	return
}

func (n *WorkflowScenarioEndpoint) Add(params *m.WorkflowScenario) (result *m.WorkflowScenario, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.WorkflowScenario.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.WorkflowScenario.GetById(id)

	return
}

func (n *WorkflowScenarioEndpoint) GetList(workflowId int64) (result []*m.WorkflowScenario, total int64, err error) {

	result, total, err = n.adaptors.WorkflowScenario.ListByWorkflow(workflowId)

	return
}

func (n *WorkflowScenarioEndpoint) Delete(workflowScenarioId int64) (err error) {

	if workflowScenarioId == 0 {
		err = errors.New("scenario id is null")
		return
	}

	err = n.adaptors.WorkflowScenario.Delete(workflowScenarioId)

	return
}

func (n *WorkflowScenarioEndpoint) Search(query string, limit, offset int) (result []*m.WorkflowScenario, total int64, err error) {

	result, total, err = n.adaptors.WorkflowScenario.Search(query, limit, offset)

	return
}

func (n *WorkflowScenarioEndpoint) Update(params *m.WorkflowScenario) (result *m.WorkflowScenario,
	errs []*validation.Error, err error) {

	var workflowScenario *m.WorkflowScenario
	if workflowScenario, err = n.adaptors.WorkflowScenario.GetById(params.Id); err != nil {
		return
	}

	if workflowScenario.WorkflowId != params.Id {
		err = fmt.Errorf("record not found")
	}

	common.Copy(&workflowScenario, &params, common.JsonEngine)

	// validation
	_, errs = workflowScenario.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.WorkflowScenario.Update(workflowScenario); err != nil {
		return
	}

	result, err = n.adaptors.WorkflowScenario.GetById(workflowScenario.Id)

	return
}