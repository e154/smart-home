package responses

import "github.com/e154/smart-home/api/server/v1/models"

// swagger:response WorkflowScenarioSearch
type WorkflowScenarioSearch struct {
	// in:body
	Body struct {
		Scenarios []*models.WorkflowScenario `json:"scenarios"`
	}
}

// swagger:response WorkflowScenarios
type WorkflowScenarios struct {
	// in:body
	Body struct {
		Scenarios []*models.WorkflowScenario `json:"scenarios"`
	}
}
