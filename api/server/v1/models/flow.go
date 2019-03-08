package models

import (
	"time"
)

type FlowListModel struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Workflow    struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"workflow"`
	Workers []struct {
		Id int64 `json:"id"`
	} `json:"workers"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FlowWorkflowModel struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FlowWorkerModel struct {
	Id             int64              `json:"id"`
	Name           string             `json:"name" valid:"MaxSize(254);Required"`
	Time           string             `json:"time" valid:"Required"`
	Status         string             `json:"status" valid:"Required"`
	Workflow       *FlowWorkflowModel `json:"workflow"`
	WorkflowId     int64              `json:"workflow_id" valid:"Required"`
	FlowId         int64              `json:"flow_id" valid:"Required"`
	DeviceAction   *DeviceAction      `json:"device_action"`
	DeviceActionId int64              `json:"device_action_id" valid:"Required"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type NewFlowModel struct {
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

type UpdateFlowModel struct {
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

type FlowModel struct {
	Id                 int64                  `json:"id"`
	Name               string                 `json:"name" valid:"MaxSize(254);Required"`
	Description        string                 `json:"description" valid:"MaxSize(254)"`
	Status             string                 `json:"status" valid:"Required"`
	Workflow           *FlowWorkflowModel     `json:"workflow"`
	WorkflowId         int64                  `json:"workflow_id" valid:"Required"`
	WorkflowScenarioId int64                  `json:"workflow_scenario_id" valid:"Required"`
	Connections        []*FlowConnectionModel `json:"connections"`
	FlowElements       []*FlowElementModel    `json:"flow_elements"`
	Workers            []*FlowWorkerModel     `json:"workers"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

type ResponseFlow struct {
	Flow *FlowModel `json:"flow"`
}

type ResponseFlowList struct {
	Items []FlowListModel `json:"items"`
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}

type ResponseSearchFlow struct {
	Flows []FlowModel `json:"flows"`
}
