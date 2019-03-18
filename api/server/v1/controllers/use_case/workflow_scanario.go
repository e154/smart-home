package use_case

import (
	"fmt"
	"errors"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/api/server/v1/models"
)

func GetWorkflowScenarioById(workflowId, scenarioId int64, adaptors *adaptors.Adaptors) (result *models.WorkflowScenario, err error) {

	var workflowScenario *m.WorkflowScenario
	if workflowScenario, err = adaptors.WorkflowScenario.GetById(scenarioId); err != nil {
		return
	}

	if workflowScenario.WorkflowId != workflowId {
		err = fmt.Errorf("record not found")
	}

	result = &models.WorkflowScenario{}
	err = common.Copy(&result, &workflowScenario)

	return
}

func AddWorkflowScenario(workflowScenario *m.WorkflowScenario, adaptors *adaptors.Adaptors) (result *models.WorkflowScenario, errs []*validation.Error, err error) {

	// validation
	_, errs = workflowScenario.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = adaptors.WorkflowScenario.Add(workflowScenario); err != nil {
		return
	}

	if workflowScenario, err = adaptors.WorkflowScenario.GetById(id); err != nil {
		return
	}

	result = &models.WorkflowScenario{}
	err = common.Copy(&result, &workflowScenario)

	return
}

func GetWorkflowScenarioList(workflowId int64, adaptors *adaptors.Adaptors) (result []*models.WorkflowScenario, total int64, err error) {

	var items []*m.WorkflowScenario
	if items, total, err = adaptors.WorkflowScenario.ListByWorkflow(workflowId); err != nil {
		return
	}

	result = make([]*models.WorkflowScenario, 0)
	err = common.Copy(&result, &items)

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

func SearchWorkflowScenario(query string, limit, offset int, adaptors *adaptors.Adaptors) (result []*models.WorkflowScenario, total int64, err error) {

	var items []*m.WorkflowScenario
	if items, total, err = adaptors.WorkflowScenario.Search(query, limit, offset); err != nil {
		return
	}

	result = make([]*models.WorkflowScenario, 0)
	err = common.Copy(&result, &items)

	return
}

func WorkflowUpdateWorkflowScenario(params *models.UpdateWorkflowScenario,
	adaptors *adaptors.Adaptors) (result *models.WorkflowScenario,
	errs []*validation.Error, err error) {

	var workflowScenario *m.WorkflowScenario
	if workflowScenario, err = adaptors.WorkflowScenario.GetById(params.Id); err != nil {
		return
	}

	if workflowScenario.WorkflowId != params.Id {
		err = fmt.Errorf("record not found")
	}

	workflowScenario.Id = params.Id
	workflowScenario.Name = params.Name
	workflowScenario.WorkflowId = int64(params.WorkflowId)
	workflowScenario.SystemName = params.SystemName
	workflowScenario.Scripts = make([]*m.Script, 0)

	for _, s := range params.Scripts {
		script := &m.Script{}
		if err = common.Copy(&script, &s); err != nil {
			log.Error(err.Error())
		}
		workflowScenario.Scripts = append(workflowScenario.Scripts, script)
	}

	// validation
	_, errs = workflowScenario.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.WorkflowScenario.Update(workflowScenario); err != nil {
		return
	}

	if workflowScenario, err = adaptors.WorkflowScenario.GetById(workflowScenario.Id); err != nil {
		return
	}

	result = &models.WorkflowScenario{}
	err = common.Copy(&result, &workflowScenario)

	return
}
