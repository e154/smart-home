package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type Workflow struct {
	Id          int64               `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
	//Flows       []*Flow             `orm:"-" json:"flows"`
	Scripts     []*Script           `json:"scripts"`
	Scenario    *WorkflowScenario   `json:"scenario"`
	Scenarios   []*WorkflowScenario `json:"scenarios"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

func (d *Workflow) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
