package core

import (
	"../models"
	"errors"
	"sync"
)

func NewFlow(model *models.Flow, workflow *Workflow) (flow *Flow, err error) {

	flow = &Flow{
		Model: model,
		workflow: workflow,
		cursor: []*models.FlowElement{},
		quit: make(chan bool),
		push: make(chan *models.FlowElement),
		pop: make(chan *models.FlowElement),
		stat: make(chan []*models.FlowElement),
	}

	// get flow elements
	var flowelements []*models.FlowElement
	if flowelements, err = models.GetFlowElementsByFlow(model); err != nil {
		return
	}

	for _, element := range flowelements {
		flow.FlowElements = append(flow.FlowElements, &FlowElement{Model:element, Flow: flow, Workflow: workflow})
	}

	// get connections
	if flow.Connections, err = models.GetConnectionsByFlow(model); err != nil {
		return
	}

	//for _, conn := range flow.Connections {
	//	for _, element := range flow.FlowElements {
	//		if element.Model.Uuid == conn.ElementFrom {
	//			conn.FlowElementFrom = element.Model
	//		} else if element.Model.Uuid == conn.ElementTo {
	//			conn.FlowElementTo = element.Model
	//		}
	//	}
	//}


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
	mutex        	sync.RWMutex
	cursor       	[]*models.FlowElement
	quit         	chan bool
	push         	chan *models.FlowElement
	pop          	chan *models.FlowElement
	stat         	chan []*models.FlowElement
}

func (f *Flow) Close() {
	f.quit <- true
}

func NewMessage(flow *Flow, message *Message) error {
	return flow.NewMessage(message)
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

func (f *Flow) PushCursor(element *models.FlowElement) {
	f.push <- element
}

func (f *Flow) PopCursor(element *models.FlowElement) {
	f.pop <- element
}

func (f *Flow) GetCursor() (cursor []*models.FlowElement) {
	f.stat <- []*models.FlowElement{}
	cursor = <- f.stat
	return
}

func (f *Flow) loop() {
	for {
		select {
		case <- f.quit:
			break
		case element := <- f.push:
			f.cursor = append(f.cursor, element)
		case element := <- f.pop:
			for i, cursor := range f.cursor {
				if cursor.Uuid == element.Uuid {
					f.cursor = append(f.cursor[:i], f.cursor[i+1:]...)
				}
			}
		case <- f.stat:
			f.stat <- f.cursor
		}
	}
}