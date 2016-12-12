package core

import (
	"../models"
	"errors"
	"log"
)

func NewFlow(model *models.Flow, workflow *Workflow) (flow *Flow, err error) {

	flow = &Flow{
		Model: model,
		workflow: workflow,
		cursor: []*FlowElement{},
		quit: make(chan bool),
	}

	// get flow elements
	var flowelements []*models.FlowElement
	if flowelements, err = models.GetFlowElementsByFlow(model); err != nil {
		return
	}

	for _, element := range flowelements {
		flowElement, err := NewFlowElement(element,flow, workflow)
		if err == nil {
			flow.FlowElements = append(flow.FlowElements, flowElement)
		} else {
			log.Println("error", err.Error())
		}
	}

	// get connections
	if flow.Connections, err = models.GetConnectionsByFlow(model); err != nil {
		return
	}

	for _, element := range flow.FlowElements {
		element.Flow = flow
		switch element.Model.PrototypeType  {
		case "MessageHandler":
			element.Prototype = &MessageHandler{}
			break
		case "MessageEmitter":
			element.Prototype = &MessageEmitter{}
			break
		case "Task":
			element.Prototype = &Task{}
			break
		case "Gateway":
			element.Prototype = &Gateway{}
			break
		}
	}

	go flow.loop()

	return
}

type Flow struct {
	Model        	*models.Flow
	workflow     	*Workflow
	Connections  	[]*models.Connection
	FlowElements 	[]*FlowElement
	Workers     	[]*Worker
	cursor       	[]*FlowElement
	quit         	chan bool
}

func (f *Flow) Close() {
	f.quit <- true
}

func (f *Flow) NewMessage(message *Message) (err error) {

	var exist bool
	for _, element := range f.FlowElements {
		if element.Prototype == nil {
			continue
		}

		if element.Model.PrototypeType != "MessageHandler" {
			continue
		}

		element.Run(message)

		exist = true
	}

	if !exist {
		err = errors.New("Message handler not found")
	}

	return
}

func (f *Flow) loop() {
	for {
		select {
		case <- f.quit:
			break
		}
	}
}