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
func (r Automation) AddTask(obj *api.NewTaskRequest) (task *m.NewTask) {
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

// ImportTask ...
func (r Automation) ImportTask(obj *api.Task) (task *m.Task) {
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
		AreaId:      obj.AreaId,
	}

	// triggers
	for _, t := range obj.Triggers {
		trigger := &m.Trigger{
			Name:       t.Name,
			PluginName: t.PluginName,
			ScriptId:   t.ScriptId,
			EntityId:   common.NewEntityIdFromPtr(t.EntityId),
			Payload:    AttributeFromApi(t.Attributes),
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
			Name:             a.Name,
			EntityId:         common.NewEntityIdFromPtr(a.EntityId),
			EntityActionName: a.EntityActionName,
			ScriptId:         a.ScriptId,
		}
		task.Actions = append(task.Actions, action)
	}
	return
}

// UpdateTask ...
func (r Automation) UpdateTask(obj *api.UpdateTaskRequest) (task *m.UpdateTask) {
	task = &m.UpdateTask{
		Id:           obj.Id,
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

// GetTaskList ...
func (r Automation) GetTaskList(list []*m.Task, total uint64, pagination common.PageParams) *api.GetTaskListResult {

	items := make([]*api.Task, 0, len(list))

	for _, i := range list {
		items = append(items, r.GetTask(i))
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

// GetTask ...
func (r Automation) GetTask(task *m.Task) (obj *api.Task) {

	obj = &api.Task{
		Id:          task.Id,
		Name:        task.Name,
		Description: task.Description,
		Enabled:     task.Enabled,
		IsLoaded:    common.Bool(task.IsLoaded),
		Area:        ToArea(task.Area),
		AreaId:      task.AreaId,
		Condition:   string(task.Condition),
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}

	// triggers
	for _, tr := range task.Triggers {
		obj.Triggers = append(obj.Triggers, &api.Trigger{
			Id:   tr.Id,
			Name: tr.Name,
		})
		obj.TriggerIds = append(obj.TriggerIds, tr.Id)
	}

	// conditions
	for _, con := range task.Conditions {
		obj.Conditions = append(obj.Conditions, &api.Condition{
			Id:   con.Id,
			Name: con.Name,
		})
		obj.ConditionIds = append(obj.ConditionIds, con.Id)
	}

	// actions
	for _, action := range task.Actions {
		obj.Actions = append(obj.Actions, &api.Action{
			Id:   action.Id,
			Name: action.Name,
		})
		obj.ActionIds = append(obj.ActionIds, action.Id)
	}

	// telemetry
	if task.Telemetry != nil {
		obj.Telemetry = make([]*api.TelemetryItem, 0)
	}

	for _, item := range task.Telemetry {
		stateItem := &api.TelemetryItem{
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
