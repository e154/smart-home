package models

import (
	"time"
)

type FlowShort struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Workflow    struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Scenario struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"scenario,omitempty"`
	} `json:"workflow"`
	Workers []struct {
		Id int64 `json:"id"`
	} `json:"workers"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FlowWorkflow struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FlowWorker struct {
	Id             int64         `json:"id"`
	Name           string        `json:"name" valid:"MaxSize(254);Required"`
	Time           string        `json:"time" valid:"Required"`
	Status         string        `json:"status" valid:"Required"`
	Workflow       *FlowWorkflow `json:"workflow"`
	WorkflowId     int64         `json:"workflow_id" valid:"Required"`
	FlowId         int64         `json:"flow_id" valid:"Required"`
	DeviceAction   *DeviceAction `json:"device_action"`
	DeviceActionId int64         `json:"device_action_id" valid:"Required"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// swagger:model
type NewFlow struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Workflow    struct {
		Id int64 `json:"id"`
	} `json:"workflow"`
	Scenario struct {
		Id int64 `json:"id"`
	} `json:"scenario"`
}

// swagger:model
type UpdateFlow struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Workflow    struct {
		Id int64 `json:"id"`
	} `json:"workflow"`
	Scenario struct {
		Id int64 `json:"id"`
	} `json:"scenario"`
}

// swagger:model
type Flow struct {
	Id                 int64               `json:"id"`
	Name               string              `json:"name" valid:"MaxSize(254);Required"`
	Description        string              `json:"description" valid:"MaxSize(254)"`
	Status             string              `json:"status" valid:"Required"`
	Workflow           *FlowWorkflow       `json:"workflow"`
	WorkflowId         int64               `json:"workflow_id" valid:"Required"`
	WorkflowScenarioId int64               `json:"workflow_scenario_id" valid:"Required"`
	Connections        []*FlowConnection   `json:"connections"`
	FlowElements       []*FlowElement      `json:"flow_elements"`
	Workers            []*FlowWorker       `json:"workers"`
	Subscriptions      []*FlowSubscription `json:"subscriptions"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}
