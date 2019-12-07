package models

import (
	"time"
)

type RedactorWorkflowModel struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RedactorConnector struct {
	Id    string `json:"id"`
	Start struct {
		Object string `json:"object"`
		Point  int64  `json:"point"`
	} `json:"start"`
	End struct {
		Object string `json:"object"`
		Point  int64  `json:"point"`
	} `json:"end"`
	FlowType  string `json:"flow_type"`
	Title     string `json:"title"`
	Direction string `json:"direction"`
}

type RedactorObject struct {
	Id   string `json:"id"`
	Type struct {
		Name   string      `json:"name"`
		Start  interface{} `json:"start"`
		End    interface{} `json:"end"`
		Status string      `json:"status"`
		Action string      `json:"action"`
	} `json:"type"`
	Position struct {
		Top  int64 `json:"top"`
		Left int64 `json:"left"`
	} `json:"position"`
	Status        string  `json:"status"`
	Error         string  `json:"error"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	PrototypeType string  `json:"prototype_type"`
	Script        *Script `json:"script"`
	FlowLink      *Flow   `json:"flow_link"`
}

// swagger:model
type RedactorFlow struct {
	Id            int64                  `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Status        string                 `json:"status"`
	Objects       []*RedactorObject      `json:"objects"`
	Connectors    []*RedactorConnector   `json:"connectors"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"update_at"`
	Workflow      *RedactorWorkflowModel `json:"workflow"`
	Subscriptions []*FlowSubscription    `json:"subscriptions"`
	Scenario      *WorkflowScenario      `json:"scenario"`
	Workers       []*Worker              `json:"workers"`
}
