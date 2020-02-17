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

package core

// Javascript Binding
//
//Workflow
//	 .GetName()
//	 .GetDescription()
//	 .SetVar(string, interface)
//	 .GetVar(string)
//	 .GetScenario() string
//	 .GetScenarioName() string
//	 .SetScenario(string)
//
type WorkflowBind struct {
	wf *Workflow
}

func (w *WorkflowBind) GetName() string {
	return w.wf.model.Name
}

func (w *WorkflowBind) GetDescription() string {
	return w.wf.model.Description
}

func (w *WorkflowBind) SetVar(key string, value interface{}) {
	w.wf.SetVar(key, value)
}

func (w *WorkflowBind) GetVar(key string) interface{} {
	return w.wf.GetVar(key)
}

func (w *WorkflowBind) GetScenario() string {
	return w.wf.model.Scenario.SystemName
}

func (w *WorkflowBind) GetScenarioName() string {
	return w.wf.model.Scenario.Name
}

func (w *WorkflowBind) SetScenario(system_name string) {
	//bug if call this method from scenario
	go w.wf.SetScenario(system_name)
}
