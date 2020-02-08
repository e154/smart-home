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

import (
	"context"
	"errors"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
)

type FlowElement struct {
	Model        *m.FlowElement
	Flow         *Flow
	Workflow     *Workflow
	ScriptEngine *scripts.Engine
	Prototype    ActionPrototypes
	status       Status
	Action       *Action
	adaptors     *adaptors.Adaptors
}

func NewFlowElement(model *m.FlowElement,
	flow *Flow,
	workflow *Workflow,
	adaptors *adaptors.Adaptors) (flowElement *FlowElement, err error) {

	flowElement = &FlowElement{
		Model:        model,
		Flow:         flow,
		Workflow:     workflow,
		adaptors:     adaptors,
		ScriptEngine: flow.scriptEngine,
	}

	switch flowElement.Model.PrototypeType {
	case "MessageHandler":
		flowElement.Prototype = &MessageHandler{}
		break
	case "MessageEmitter":
		flowElement.Prototype = &MessageEmitter{}
		break
	case "Task":
		flowElement.Prototype = &Task{}
		break
	case "Gateway":
		flowElement.Prototype = &Gateway{}
		break
	case "Flow":
		flowElement.Prototype = &FlowLink{}
		break
	}

	return
}

func (m *FlowElement) Before() error {

	m.status = DONE
	return m.Prototype.Before(m.Flow)
}

// run internal process
func (m *FlowElement) Run(ctx context.Context) (b bool, err error) {

	if err = ctx.Err(); err != nil {
		return
	}

	m.status = IN_PROCESS

	//???
	m.Flow.cursor = m.Model.Uuid

	if err = m.Before(); err != nil {
		m.status = ERROR
		return
	}

	if err = m.Prototype.Run(m.Flow); err != nil {
		m.status = ERROR
		return
	}

	//run script if exist
	if m.Model.Script != nil {

		if err = m.ScriptEngine.EvalString(m.Model.Script.Compiled); err != nil {
			m.status = ERROR
			return
		}

		if m.Flow.message.Error != "" {
			err = errors.New(m.Flow.message.Error)
			m.status = ERROR
			return
		}

		b = m.Flow.message.Success
	}

	if err = m.After(); err != nil {
		m.status = ERROR
		return
	}

	m.status = ENDED

	return
}

func (m *FlowElement) After() error {
	m.status = DONE
	return m.Prototype.After(m.Flow)
}

func (m *FlowElement) GetStatus() (status Status) {

	status = m.status
	return
}
