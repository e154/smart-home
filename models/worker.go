// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

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
