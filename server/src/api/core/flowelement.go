package core

import (
	"errors"
	"sync"
	"../models"
	"../scripts"
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

	//flowElement.Script("device", device)
	//flowElement.Script("flow", flow.Model)
	//flowElement.Script("node", node)
	//flowElement.Script("message", message)

	//debug(flowElement.Script)

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

	m.mutex.Lock()
	m.status = DONE
	m.mutex.Unlock()
	return m.Prototype.Before(message, m.Flow)
}

// run internal process
func (m *FlowElement) Run(message *Message) (err error) {

	m.mutex.Lock()
	m.status = IN_PROCESS
	m.mutex.Unlock()

	//m.Flow.PushCursor(m)
	err = m.Before(message)
	err = m.Prototype.Run(message, m.Flow)

	//run script if exist
	if m.Script != nil {
		m.Script.Do()
	}

	err = m.After(message)

	// each connections
	var elements []*FlowElement
	for _, conn := range m.Flow.Connections {
		if conn.ElementFrom != m.Model.Uuid || conn.ElementTo == m.Model.Uuid {
			continue
		}

		for _, element := range m.Flow.FlowElements {
			if conn.ElementTo != element.Model.Uuid {
				continue
			}
			elements = append(elements, element)
		}
	}

	for _, element := range elements {
		go element.Run(message)
	}

	//m.Flow.PopCursor(m)

	m.mutex.Lock()
	m.status = ENDED
	m.mutex.Unlock()

	if len(elements) == 0 {
		err = errors.New("Element not found")
	}

	return
}

func (m *FlowElement) After(message *Message) error {

	m.mutex.Lock()
	m.status = STARTED
	m.mutex.Unlock()

	return  m.Prototype.After(message, m.Flow)
}

func (m *FlowElement) GetStatus() (status Status) {

	m.mutex.Lock()
	status = m.status
	m.mutex.Unlock()

	return
}