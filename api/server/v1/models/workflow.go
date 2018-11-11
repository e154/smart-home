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

type Workflows []*Workflow

type ResponseWorkflow struct {
	Code ResponseType `json:"code"`
	Data struct {
		Workflow *Workflow `json:"workflow"`
	} `json:"data"`
}

type ResponseWorkflowList struct {
	Code ResponseType `json:"code"`
	Data struct {
		Items  []*Workflow `json:"items"`
		Limit  int64       `json:"limit"`
		Offset int64       `json:"offset"`
		Total  int64       `json:"total"`
	} `json:"data"`
}
