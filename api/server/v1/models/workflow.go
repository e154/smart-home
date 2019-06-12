package models

import "time"

// swagger:model
type NewWorkflow struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// swagger:model
type UpdateWorkflow struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// swagger:model
type Workflow struct {
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

// swagger:model
type WorkflowUpdateWorkflowScenario struct {
	WorkflowScenarioId int64 `json:"workflow_scenario_id"`
}
