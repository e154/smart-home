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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Automation ...
type Automation struct{}

// NewAutomationDto ...
func NewAutomationDto() Automation {
	return Automation{}
}

// AddTask ...
func (r Automation) AddTask(obj *api.NewTaskRequest) (task *m.Task) {
	if obj == nil {
		return
	}
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
		trigger := &m.Trigger{
			Name:       t.Name,
			PluginName: t.PluginName,
			Payload:    AttributeFromApi(t.Attributes),
		}
		if t.Script != nil {
			trigger.ScriptId = common.Int64(t.Script.Id)
		}
		if t.Entity != nil {
			entityId := common.EntityId(t.Entity.Id)
			trigger.EntityId = &entityId
		}
		task.Triggers = append(task.Triggers, trigger)
	}
	// conditions
	for _, c := range obj.Conditions {
		condition := &m.Condition{
			Name: c.Name,
		}
		if c.Script != nil {
			condition.ScriptId = c.Script.Id
		}
		task.Conditions = append(task.Conditions, condition)
	}
	// actions
	for _, a := range obj.Actions {
		action := &m.Action{
			Name:           a.Name,
			EntityActionId: a.EntityActionId,
		}
		if a.Script != nil {
			action.ScriptId = common.Int64(a.Script.Id)
		}
		task.Actions = append(task.Actions, action)
	}
	return
}

// ImportTask ...
func (r Automation) ImportTask(obj *api.NewTaskRequest) (task *m.Task) {
	if obj == nil {
		return
	}
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
		_, script := ImportScript(t.Script)
		trigger := &m.Trigger{
			Name:       t.Name,
			Script:     script,
			PluginName: t.PluginName,
			Payload:    AttributeFromApi(t.Attributes),
		}
		if t.Script != nil {
			trigger.ScriptId = common.Int64(t.Script.Id)
		}
		if t.Entity != nil {
			entityId := common.EntityId(t.Entity.Id)
			trigger.EntityId = &entityId
		}
		task.Triggers = append(task.Triggers, trigger)
	}
	// conditions
	for _, c := range obj.Conditions {
		condition := &m.Condition{
			Name: c.Name,
		}
		if c.Script != nil {
			_, script := ImportScript(c.Script)
			condition.ScriptId = c.Script.Id
			condition.Script = script
		}
		task.Conditions = append(task.Conditions, condition)
	}
	// actions
	for _, a := range obj.Actions {
		action := &m.Action{
			Name:           a.Name,
			EntityActionId: a.EntityActionId,
		}
		if a.Script != nil {
			_, script := ImportScript(a.Script)
			action.ScriptId = common.Int64(a.Script.Id)
			action.Script = script
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
		trigger := &m.Trigger{
			Name:       t.Name,
			PluginName: t.PluginName,
			Payload:    AttributeFromApi(t.Attributes),
		}
		if t.Entity != nil {
			entityId := common.EntityId(t.Entity.Id)
			trigger.EntityId = &entityId
		}
		if t.Script != nil {
			trigger.ScriptId = common.Int64(t.Script.Id)
		}
		task.Triggers = append(task.Triggers, trigger)
	}
	// conditions
	for _, c := range obj.Conditions {
		condition := &m.Condition{
			Name: c.Name,
		}
		if c.Script != nil {
			condition.ScriptId = c.Script.Id
		}
		task.Conditions = append(task.Conditions, condition)
	}
	// actions
	for _, a := range obj.Actions {
		action := &m.Action{
			Name:           a.Name,
			EntityActionId: a.EntityActionId,
		}
		if a.Script != nil {
			action.ScriptId = common.Int64(a.Script.Id)
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
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToTask ...
func (r Automation) ToTask(task *m.Task) (obj *api.Task) {

	obj = &api.Task{
		Id:          task.Id,
		Name:        task.Name,
		Description: task.Description,
		Enabled:     task.Enabled,
		Area:        ToArea(task.Area),
		Condition:   string(task.Condition),
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}

	// triggers
	for _, tr := range task.Triggers {
		obj.Triggers = append(obj.Triggers, &api.Trigger{
			Name:       tr.Name,
			Script:     ToGScript(tr.Script),
			PluginName: tr.PluginName,
			Entity: &api.Trigger_Entity{
				Id: tr.EntityId.String(),
			},
			Attributes: AttributeToApi(tr.Payload),
		})
	}

	// conditions
	for _, con := range task.Conditions {
		obj.Conditions = append(obj.Conditions, &api.Condition{
			Name:   con.Name,
			Script: ToGScript(con.Script),
		})
	}

	// actions
	for _, con := range task.Actions {
		obj.Actions = append(obj.Actions, &api.Action{
			Name:   con.Name,
			Script: ToGScript(con.Script),
		})
	}

	return
}
