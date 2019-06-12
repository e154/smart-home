package core

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/adaptors"
	"errors"
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

func (m *FlowElement) Before(message *Message) error {

	m.status = DONE
	return m.Prototype.Before(message, m.Flow)
}

// run internal process
func (m *FlowElement) Run(msg *Message) (b bool, returnMessage *Message, err error) {

	message := msg.Copy()

	m.status = IN_PROCESS

	//???
	m.Flow.cursor = m.Model.Uuid

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
	if m.Model.Script != nil {

		if _, err = m.ScriptEngine.PushStruct("message", message); err != nil {
			m.status = ERROR
			return
		}

		if err = m.ScriptEngine.EvalString(m.Model.Script.Compiled); err != nil {
			m.status = ERROR
			return
		}

		if message.Error != "" {
			err = errors.New(message.Error)
			m.status = ERROR
			return
		}

		b = message.Success
	}

	if err = m.After(message); err != nil {
		m.status = ERROR
		return
	}

	m.status = ENDED

	return
}

func (m *FlowElement) After(message *Message) error {
	m.status = DONE
	return m.Prototype.After(message, m.Flow)
}

func (m *FlowElement) GetStatus() (status Status) {

	status = m.status
	return
}