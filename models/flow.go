package models

import (
	"time"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
)

type Flow struct {
	Id                 int64          `json:"id"`
	Name               string         `json:"name" valid:"MaxSize(254);Required"`
	Description        string         `json:"description" valid:"MaxSize(254)"`
	Status             StatusType     `json:"status" valid:"Required"`
	Workflow           *Workflow      `json:"workflow"`
	WorkflowId         int64          `json:"workflow_id" valid:"Required"`
	WorkflowScenarioId int64          `json:"workflow_scenario_id" valid:"Required"`
	Connections        []*Connection  `json:"connections"`
	FlowElements       []*FlowElement `json:"flow_elements"`
	Workers            []*Worker      `json:"workers"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

func (d *Flow) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
