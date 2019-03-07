package models

import "time"

type NewWorkflow struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateWorkflow struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type WorkflowModel struct {
	Id          int64               `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
	Scripts     []*Script           `json:"scripts"`
	Scenario    *WorkflowScenario   `json:"scenario"`
	Scenarios   []*WorkflowScenario `json:"scenarios"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type Workflows []*WorkflowModel

type ResponseWorkflowModel struct {
	Code ResponseType `json:"code"`
	Data struct {
		WorkflowModel *WorkflowModel `json:"workflow"`
	} `json:"data"`
}

type ResponseWorkflowList struct {
	Code ResponseType `json:"code"`
	Data struct {
		Items  []*WorkflowModel `json:"items"`
		Limit  int64            `json:"limit"`
		Offset int64            `json:"offset"`
		Total  int64            `json:"total"`
	} `json:"data"`
}

type SearchWorkflowResponse struct {
	Workflows []WorkflowModel `json:"workflows"`
}

type WorkflowUpdateWorkflowScenario struct {
	WorkflowScenarioId int64 `json:"workflow_scenario_id"`
}
