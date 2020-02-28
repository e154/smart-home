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
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
	"time"
)

type Flow struct {
	Id                 int64                `json:"id"`
	Name               string               `json:"name" valid:"MaxSize(254);Required"`
	Description        string               `json:"description" valid:"MaxSize(254)"`
	Status             StatusType           `json:"status" valid:"Required"`
	Workflow           *Workflow            `json:"workflow"`
	WorkflowId         int64                `json:"workflow_id" valid:"Required"`
	WorkflowScenarioId int64                `json:"workflow_scenario_id" valid:"Required"`
	Connections        []*Connection        `json:"connections"`
	FlowElements       []*FlowElement       `json:"flow_elements"`
	Workers            []*Worker            `json:"workers"`
	Subscriptions      []*FlowSubscription  `json:"subscriptions"`
	Zigbee2mqttDevices []*Zigbee2mqttDevice `json:"zigbee2mqtt_devices"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

func (d *Flow) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
