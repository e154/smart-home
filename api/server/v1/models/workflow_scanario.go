package models

import "time"

// swagger:model
type WorkflowScenario struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	SystemName string    `json:"system_name"`
	WorkflowId int64     `json:"workflow_id"`
	Scripts    []*Script `json:"scripts"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// swagger:model
type NewWorkflowScenario struct {
	Name       string `json:"name"`
	SystemName string `json:"system_name"`
	WorkflowId int64  `json:"workflow_id"`
}

// swagger:model
type UpdateWorkflowScenario struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	SystemName string    `json:"system_name"`
	WorkflowId int64     `json:"workflow_id"`
	Scripts    []*Script `json:"scripts"`
}

type WorkflowScenarioListModel struct {
	Scenarios []*WorkflowScenario `json:"scenarios"`
}

type SearchWorkflowScenarioResponse struct {
	WorkflowScenarios []WorkflowScenario `json:"scenarios"`
}
