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

package example1

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
)

// TriggerManager ...
type TriggerManager struct {
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
}

// NewTriggerManager ...
func NewTriggerManager(adaptors *adaptors.Adaptors) *TriggerManager {
	return &TriggerManager{
		adaptors: adaptors,
	}
}

func (t *TriggerManager) Create(scripts []*m.Script,
	entities []*m.Entity) []*m.Task {

	var script *m.Script
	if len(scripts) > 0 {
		script = scripts[0]
	}

	t.addTimerTask("l3n1_timer1", script, entities[0])
	t.addTimerTask("l3n2_timer2", script, entities[1])
	t.addTimerTask("l3n3_timer3", script, entities[2])
	t.addTimerTask("l3n4_timer4", script, entities[3])

	t.addCheckTask("l3n1_check1", script, entities[0])
	t.addCheckTask("l3n2_check2", script, entities[0])
	t.addCheckTask("l3n3_check3", script, entities[0])
	t.addCheckTask("l3n4_check4", script, entities[0])
	return []*m.Task{}
}

func (t *TriggerManager) addTimerTask(name string,
	script *m.Script,
	entity *m.Entity) (task *m.Task) {

	task = &m.Task{
		Name:      fmt.Sprintf("task_%s", name),
		Enabled:   true,
		Condition: common.ConditionAnd,
	}
	task.AddTrigger(&m.Trigger{
		Name:       fmt.Sprintf("trigger_%s", name),
		Script:     script,
		PluginName: "time",
		EntityId:   &entity.Id,
		Payload: m.Attributes{
			triggers.CronOptionTrigger: {
				Name:  triggers.CronOptionTrigger,
				Type:  common.AttributeString,
				Value: "0,5,10,15,20,25,30,35,40,45,50,55 * * * * *", //every 5 seconds
			},
		},
	})
	err := t.adaptors.Task.Add(task)
	So(err, ShouldBeNil)

	return
}

func (t *TriggerManager) addCheckTask(name string,
	script *m.Script,
	entity *m.Entity) (task *m.Task) {

	task = &m.Task{
		Name:      fmt.Sprintf("task_%s", name),
		Enabled:   true,
		Condition: common.ConditionAnd,
	}
	task.AddTrigger(&m.Trigger{
		Name:       fmt.Sprintf("trigger_%s", name),
		Script:     script,
		PluginName: "state_change",
		EntityId:   &entity.Id,
	})
	task.AddCondition(&m.Condition{
		Name:   fmt.Sprintf("condition_%s", name),
		Script: script,
	})
	task.AddAction(&m.Action{
		Name:   fmt.Sprintf("action_%s", name),
		Script: script,
	})
	err := t.adaptors.Task.Add(task)
	So(err, ShouldBeNil)

	return
}
