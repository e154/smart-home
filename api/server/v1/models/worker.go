package models

import "time"

type WorkerModel struct {
	Id             int64          `json:"id"`
	Name           string         `json:"name"`
	Time           string         `json:"time"`
	Status         string         `json:"status"`
	Workflow       *WorkflowModel `json:"workflow"`
	WorkflowId     int64          `json:"workflow_id"`
	Flow           *FlowModel     `json:"flow"`
	FlowId         int64          `json:"flow_id"`
	DeviceAction   *DeviceAction  `json:"device_action"`
	DeviceActionId int64          `json:"device_action_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
