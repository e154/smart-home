package core

import (
	"sync"
	"../models"
	"../scripts"
	r "../../lib/rpc"
)

func NewFlowElement(model *models.FlowElement, flow *Flow, workflow *Workflow) (flowElement *FlowElement, err error) {

	flowElement = &FlowElement{
		Model:model,
		Flow: flow,
		Workflow: workflow,
	}

	if model.Script == nil {
		return
	}

	var script *models.Script
	if script, err = models.GetScriptById(model.Script.Id); err != nil {
		return
	}

	if flowElement.Script, err = scripts.New(script); err != nil {
		return
	}

	flowElement.Script.PushStruct("request", &r.Request{})
	flowElement.Script.PushFunction("modbus_send", func(args *r.Request) (result r.Result) {

		if flow.Node == nil {
			result.Error = "Node is nil pointer"
			return
		}

		if err := flow.Node.ModbusSend(args, &result); err != nil {
			result.Error = err.Error()
		}

		return
	})


	return
}

type FlowElement struct {
	Model 		*models.FlowElement
	Flow		*Flow
	Workflow	*Workflow
	Script		*scripts.Engine
	Prototype	ActionPrototypes
	status		Status
	mutex     	sync.Mutex
}

func (m *FlowElement) Before(message *Message) error {

	m.status = DONE
	return m.Prototype.Before(message, m.Flow)
}

// run internal process
func (m *FlowElement) Run(message *Message) (res string, err error) {

	m.status = IN_PROCESS

	//m.Flow.PushCursor(m)
	err = m.Before(message)
	err = m.Prototype.Run(message, m.Flow)

	if m.Script != nil {
		m.Script.PushStruct("message", message)
		res, _ = m.Script.Do()
	} else {
		res = "false"
	}

	err = m.After(message)
	//m.Flow.PopCursor(m)

	m.status = ENDED

	return
}

func (m *FlowElement) After(message *Message) error {
	m.status = STARTED
	return  m.Prototype.After(message, m.Flow)
}

func (m *FlowElement) GetStatus() (status Status) {

	status = m.status
	return
}