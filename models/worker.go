package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type Worker struct {
	Id             int64         `json:"id"`
	Name           string        `json:"name" valid:"MaxSize(254);Required"`
	Time           string        `json:"time" valid:"Required"`
	Status         string        `json:"status" valid:"Required"`
	Workflow       *Workflow     `json:"workflow"`
	WorkflowId     int64         `json:"workflow_id" valid:"Required"`
	Flow           *Flow         `json:"flow"`
	FlowId         int64         `json:"flow_id" valid:"Required"`
	DeviceAction   *DeviceAction `json:"device_action"`
	DeviceActionId int64         `json:"device_action_id" valid:"Required"`
	Workers        []*Worker     `json:"workers"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

func (d *Worker) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}