package core

import (
	"sync"
	"../models"
	"errors"
)

func NewFlow(model *models.Flow, workflow *Workflow) (flow *Flow) {

	flow = &Flow{
		model: model,
		workflow: workflow,
		cursor: []*models.FlowElement{},
		quit: make(chan bool),
		push: make(chan *models.FlowElement),
		pop: make(chan *models.FlowElement),
		stat: make(chan []*models.FlowElement),
	}

	go flow.loop()

	return
}

type Flow struct {
	model 		*models.Flow
	workflow	*Workflow
	//Connections 	[]*models.Connection
	//FlowElements	[]*models.FlowElement
	//Workers     	[]*Worker
	mutex       	sync.RWMutex
	cursor      	[]*models.FlowElement
	quit		chan bool
	push 		chan *models.FlowElement
	pop 		chan *models.FlowElement
	stat		chan []*models.FlowElement
}

func (f *Flow) Close() {
	f.quit <- true
}

func NewMessage(flow *Flow, message *models.Message) error {
	return flow.NewMessage(message)
}

func (f *Flow) NewMessage(message *models.Message) (err error) {

	var exist bool
	for _, element := range f.model.FlowElements {
		if element.Prototype == nil {
			continue
		}

		if element.PrototypeType != "MessageHandler" {
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