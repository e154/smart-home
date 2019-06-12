package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type WorkflowScenario struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name" valid:"MaxSize(254);Required"`
	SystemName string    `json:"system_name" valid:"MaxSize(254);Required"`
	WorkflowId int64     `json:"workflow_id" valid:"Required"`
	Scripts    []*Script `json:"scripts"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (d *WorkflowScenario) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}