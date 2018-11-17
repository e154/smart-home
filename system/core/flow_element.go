package core

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/adaptors"
	"errors"
)

type FlowElement struct {
	Model     *m.FlowElement
	Flow      *Flow
	Workflow  *Workflow
	Script    *scripts.Engine
	Prototype ActionPrototypes
	status    Status
	Action    *Action
	adaptors  *adaptors.Adaptors
}

func NewFlowElement(model *m.FlowElement,
	flow *Flow,
	workflow *Workflow,
	adaptors *adaptors.Adaptors) (flowElement *FlowElement, err error) {

	flowElement = &FlowElement{
		Model:    model,
		Flow:     flow,
		Workflow: workflow,
		adaptors: adaptors,
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

	if model.Script == nil {
		return
	}

	if flowElement.Script, err = flow.NewScript(model.Script); err != nil {
		return
	}

	return
}

func (m *FlowElement) Before(message *Message) error {

	m.status = DONE
	return m.Prototype.Before(message, m.Flow)
}

// run internal process
func (m *FlowElement) Run(message *Message) (b bool, returnMessage *Message, err error) {

	m.status = IN_PROCESS

	//m.Flow.PushCursor(m)
	if err = m.Before(message); err != nil {
		m.status = ERROR
		return
	}

	if err = m.Prototype.Run(message, m.Flow); err != nil {
		m.status = ERROR
		return
	}

	returnMessage = &Message{}
	*returnMessage = *message

	//run script if exist
	if m.Script != nil {

		m.Script.PushStruct("message", message)

		var ok string
		if ok, err = m.Script.Do(); err != nil {
			m.status = ERROR
			return
		}
		//TODO refactor message system
		if message.Error != "" {
			err = errors.New(message.Error)
			m.status = ERROR
			return
		}

		b = ok == "true"
	}

	if err = m.After(message); err != nil {
		m.status = ERROR
		return
	}

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