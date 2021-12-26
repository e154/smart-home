// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// Automation ...
type Automation struct{}

// NewAutomationDto ...
func NewAutomationDto() Automation {
	return Automation{}
}

// AddTask ...
func (r Automation) AddTask(obj *api.NewTaskRequest) (task *m.Task) {
	task = &m.Task{
		Name:        obj.Name,
		Description: obj.Description,
		Enabled:     obj.Enabled,
		Condition:   common.ConditionType(obj.Condition),
		Triggers:    make([]*m.Trigger, 0, len(obj.Triggers)),
		Conditions:  make([]*m.Condition, 0, len(obj.Conditions)),
		Actions:     make([]*m.Action, 0, len(obj.Actions)),
	}
	// area
	if obj.Area != nil {
		task.Area = &m.Area{
			Id:          obj.Area.Id,
			Name:        obj.Area.Name,
			Description: obj.Area.Description,
		}
	}
	// triggers
	for _, t := range obj.Triggers {
		entityId := common.EntityId(t.EntityId)
		trigger := &m.Trigger{
			Name:       t.Name,
			TaskId:     t.TaskId,
			EntityId:   &entityId,
			ScriptId:   t.ScriptId,
			PluginName: t.PluginName,
			Payload:    AttributeFromApi(t.Payload),
		}
		task.Triggers = append(task.Triggers, trigger)
	}
	// conditions
	for _, c := range obj.Conditions {
		condition := &m.Condition{
			Name:     c.Name,
			TaskId:   c.TaskId,
			ScriptId: c.ScriptId,
		}
		task.Conditions = append(task.Conditions, condition)
	}
	// actions
	for _, a := range obj.Actions {
		action := &m.Action{
			Name:     a.Name,
			TaskId:   a.TaskId,
			ScriptId: a.ScriptId,
		}
		task.Actions = append(task.Actions, action)
	}
	return
}

// UpdateTask ...
func (r Automation) UpdateTask(obj *api.UpdateTaskRequest) (task *m.Task) {
	task = &m.Task{
		Id:          obj.Id,
		Name:        obj.Name,
		Description: obj.Description,
		Enabled:     obj.Enabled,
		Condition:   common.ConditionType(obj.Condition),
		Triggers:    make([]*m.Trigger, 0, len(obj.Triggers)),
		Conditions:  make([]*m.Condition, 0, len(obj.Conditions)),
		Actions:     make([]*m.Action, 0, len(obj.Actions)),
	}
	// area
	if obj.Area != nil {
		task.Area = &m.Area{
			Id:          obj.Area.Id,
			Name:        obj.Area.Name,
			Description: obj.Area.Description,
		}
	}
	// triggers
	for _, t := range obj.Triggers {
		entityId := common.EntityId(t.EntityId)
		trigger := &m.Trigger{
			Name:       t.Name,
			TaskId:     t.TaskId,
			EntityId:   &entityId,
			ScriptId:   t.ScriptId,
			PluginName: t.PluginName,
			Payload:    AttributeFromApi(t.Payload),
		}
		task.Triggers = append(task.Triggers, trigger)
	}
	// conditions
	for _, c := range obj.Conditions {
		condition := &m.Condition{
			Name:     c.Name,
			TaskId:   c.TaskId,
			ScriptId: c.ScriptId,
		}
		task.Conditions = append(task.Conditions, condition)
	}
	// actions
	for _, a := range obj.Actions {
		action := &m.Action{
			Name:     a.Name,
			TaskId:   a.TaskId,
			ScriptId: a.ScriptId,
		}
		task.Actions = append(task.Actions, action)
	}
	return
}

// ToListResult ...
func (r Automation) ToListResult(list []*m.Task, total uint64, pagination common.PageParams) *api.GetTaskListResult {

	items := make([]*api.Task, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToTask(i))
	}

	return &api.GetTaskListResult{
		Items: items,
		Meta: &api.GetTaskListResult_Meta{
			Limit:        uint64(pagination.Limit),
			ObjectsCount: total,
			Offset:       uint64(pagination.Offset),
		},
	}
}

// ToTask ...
func (r Automation) ToTask(task *m.Task) (obj *api.Task) {

	return
}
