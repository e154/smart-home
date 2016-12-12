package core

import (
	"../models"
	"errors"
	"sync"
)

type FlowElement struct {
	Model 		*models.FlowElement
	Flow		*Flow
	Workflow	*Workflow
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

	m.Flow.PushCursor(m)
	err = m.Before(message)
	err = m.Prototype.Run(message, m.Flow)
	err = m.After(message)
	m.Flow.PopCursor(m)

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