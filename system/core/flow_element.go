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
	"fmt"
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

func (f *FlowElement) Before(ctx context.Context) (newCtx context.Context, err error) {

	// circular dependency search
	if newCtx, err = f.defineCircularConnection(ctx); err != nil {
		return
	}

	if err = f.Prototype.Before(f.Flow); err != nil {
		return
	}

	f.status = DONE

	return
}

// run internal process
func (f *FlowElement) Run(ctx context.Context) (newCtx context.Context, b bool, err error) {

	if err = ctx.Err(); err != nil {
		return
	}

	f.status = IN_PROCESS

	//???
	f.Flow.cursor = f.Model.Uuid

	if newCtx, err = f.Before(ctx); err != nil {
		f.status = ERROR
		return
	}

	if err = f.Prototype.Run(f.Flow); err != nil {
		f.status = ERROR
		return
	}

	//run script if exist
	if f.Model.Script != nil {

		if err = f.ScriptEngine.EvalString(f.Model.Script.Compiled); err != nil {
			f.status = ERROR
			return
		}

		if f.Flow.message.Error != "" {
			err = errors.New(f.Flow.message.Error)
			f.status = ERROR
			return
		}

		b = f.Flow.message.Success
	}

	if err = f.After(); err != nil {
		f.status = ERROR
		return
	}

	f.status = ENDED

	return
}

func (f *FlowElement) After() error {
	f.status = DONE
	return f.Prototype.After(f.Flow)
}

func (f *FlowElement) GetStatus() (status Status) {

	status = f.status
	return
}

func (f *FlowElement) defineCircularConnection(ctx context.Context) (newCtx context.Context, err error) {

	const max = 0
	var counter = 0

	if v := ctx.Value("flow_elements"); v != nil {

		if flowElements, ok := v.([]string); ok {

			for _, parentId := range flowElements {
				if parentId == f.Model.Uuid.String() {
					counter++
				}
			}

			if counter > max {
				depends := fmt.Sprintf("%s", flowElements[0])
				for _, parentId := range flowElements[1:] {
					depends = fmt.Sprintf("%s -> %s", depends, parentId)
				}
				err = fmt.Errorf("circular relationship detected flow(%v): tree(%v -> %s)", f.Flow.Model.Name, depends, f.Model.Uuid.String())
				return
			}

			if flowElements, ok := v.([]string); ok {
				flowElements = append(flowElements, f.Model.Uuid.String())
				newCtx = context.WithValue(ctx, "flow_elements", flowElements)
			} else {
				err = fmt.Errorf("bad flow_elements context value: flow_elements(%v)", flowElements)
				return
			}

		} else {
			err = fmt.Errorf("bad parent context value: flowElements(%v), flow(%v)", flowElements, f.Flow.Model.Name)
		}

		return
	}

	newCtx = context.WithValue(ctx, "flow_elements", []string{f.Model.Uuid.String()})

	return
}
