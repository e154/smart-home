// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package dto

import (
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

// Automation ...
type Automation struct{}

// NewAutomationDto ...
func NewAutomationDto() Automation {
	return Automation{}
}

// AddTask ...
func (r Automation) AddTask(obj *stub.ApiNewTaskRequest) (task *m.NewTask) {
	if obj == nil {
		return
	}
	task = &m.NewTask{
		Name:         obj.Name,
		Description:  obj.Description,
		Enabled:      obj.Enabled,
		Condition:    common.ConditionType(obj.Condition),
		TriggerIds:   obj.TriggerIds,
		ConditionIds: obj.ConditionIds,
		ActionIds:    obj.ActionIds,
		AreaId:       obj.AreaId,
	}
	return
}

// UpdateTask ...
func (r Automation) UpdateTask(obj *stub.AutomationServiceUpdateTaskJSONBody, id int64) (task *m.UpdateTask) {
	task = &m.UpdateTask{
		Id:           id,
		Name:         obj.Name,
		Description:  obj.Description,
		Enabled:      obj.Enabled,
		Condition:    common.ConditionType(obj.Condition),
		TriggerIds:   obj.TriggerIds,
		ConditionIds: obj.ConditionIds,
		ActionIds:    obj.ActionIds,
		AreaId:       obj.AreaId,
	}

	return
}

// ToListResult ...
func (r Automation) ToListResult(list []*m.Task) []*stub.ApiTask {

	items := make([]*stub.ApiTask, 0, len(list))

	for _, i := range list {
		items = append(items, r.GetTask(i))
	}

	return items
}

// GetTask ...
func (r Automation) GetTask(task *m.Task) (obj *stub.ApiTask) {

	obj = &stub.ApiTask{
		Id:          task.Id,
		Name:        task.Name,
		Description: task.Description,
		Enabled:     task.Enabled,
		IsLoaded:    common.Bool(task.IsLoaded),
		Area:        GetStubArea(task.Area),
		AreaId:      task.AreaId,
		Condition:   string(task.Condition),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	// triggers
	for _, tr := range task.Triggers {
		obj.Triggers = append(obj.Triggers, stub.ApiTrigger{
			Id:   tr.Id,
			Name: tr.Name,
		})
		obj.TriggerIds = append(obj.TriggerIds, tr.Id)
	}

	// conditions
	for _, con := range task.Conditions {
		obj.Conditions = append(obj.Conditions, stub.ApiCondition{
			Id:   con.Id,
			Name: con.Name,
		})
		obj.ConditionIds = append(obj.ConditionIds, con.Id)
	}

	// actions
	for _, action := range task.Actions {
		obj.Actions = append(obj.Actions, stub.ApiAction{
			Id:   action.Id,
			Name: action.Name,
		})
		obj.ActionIds = append(obj.ActionIds, action.Id)
	}

	// telemetry
	if task.Telemetry != nil {
		obj.Telemetry = make([]stub.ApiTelemetryItem, 0)
	}

	for _, item := range task.Telemetry {
		stateItem := stub.ApiTelemetryItem{
			Name:         item.Name,
			Num:          int32(item.Num),
			Start:        item.Start.UnixNano(),
			End:          nil,
			TimeEstimate: int64(item.TimeEstimate),
			Attributes:   item.Attributes,
			Status:       string(item.Status),
			Level:        int32(item.Level),
		}
		if item.End != nil {
			stateItem.End = common.Int64(item.End.UnixNano())
		}
		obj.Telemetry = append(obj.Telemetry, stateItem)
	}

	return
}
