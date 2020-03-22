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
	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	cr "github.com/e154/smart-home/system/cron"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/uuid"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"sync"
	"time"
)

type Flow struct {
	Storage
	Model            *m.Flow
	workflow         *Workflow
	Connections      []*m.Connection
	FlowElements     []*FlowElement
	cursor           uuid.UUID
	Node             *Node
	adaptors         *adaptors.Adaptors
	scriptService    *scripts.ScriptService
	scriptEngine     *scripts.Engine
	cron             *cr.Cron
	core             *Core
	nextScenario     bool
	mqttClient       *mqtt.Client
	mqttMessageQueue chan *Message
	mqttWorkerQuit   chan struct{}
	message          *Message
	zigbee2mqtt      *zigbee2mqtt.Zigbee2mqtt
	sync.Mutex
	isRunning bool
	Workers   map[int64]*Worker
}

func NewFlow(model *m.Flow,
	workflow *Workflow,
	adaptors *adaptors.Adaptors,
	scripts *scripts.ScriptService,
	cron *cr.Cron,
	core *Core,
	mqtt *mqtt.Mqtt,
	zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt) (flow *Flow, err error) {

	flow = &Flow{
		Storage:          NewStorage(),
		Model:            model,
		workflow:         workflow,
		adaptors:         adaptors,
		scriptService:    scripts,
		Workers:          make(map[int64]*Worker),
		cron:             cron,
		core:             core,
		mqttMessageQueue: make(chan *Message),
		mqttWorkerQuit:   make(chan struct{}),
		message:          NewMessage(),
		zigbee2mqtt:      zigbee2mqtt,
	}

	if flow.scriptEngine, err = flow.NewScript(); err != nil {
		return
	}

	for _, element := range flow.Model.FlowElements {
		var flowElement *FlowElement
		if flowElement, err = NewFlowElement(element, flow, workflow, adaptors); err == nil {
			flow.FlowElements = append(flow.FlowElements, flowElement)
		} else {
			log.Warn(err.Error())
		}
	}

	for _, conn := range flow.Model.Connections {
		flow.Connections = append(flow.Connections, conn)
	}

	// add worker
	err = flow.InitWorkers()

	go flow.mqttMessageWorker()

	// mqtt client
	flow.mqttClient = mqtt.NewClient(fmt.Sprintf("flow_%v", flow.Model.Name))

	// raw topic subscriptions
	for _, subParams := range flow.Model.Subscriptions {

		topic := fmt.Sprintf("%s", subParams.Topic)
		flow.mqttClient.Subscribe(topic, flow.mqttOnPublish)
	}

	// zigbee2mqtt devices
	var topic string
	for _, device := range flow.Model.Zigbee2mqttDevices {
		if topic, err = zigbee2mqtt.GetTopicByDevice(device); err != nil {
			log.Error(err.Error())
			continue
		}
		flow.mqttClient.Subscribe(topic, flow.mqttOnPublish)
	}

	return
}

func (f *Flow) Remove() {

	log.Infof("Remove flow '%v'", f.Model.Name)

	for _, worker := range f.Workers {
		f.RemoveWorker(worker.Model)
	}

	f.mqttWorkerQuit <- struct{}{}

	if f.mqttClient != nil {
		f.mqttClient.UnsubscribeAll()
	}

	timeout := time.After(3 * time.Second)
	for {
		time.Sleep(time.Second * 1)
		if !f.isRunning {
			log.Infof("flow %v ... ok", f.Model.Id)
			break
		}

		select {
		case <-timeout:
			return
		default:

		}
	}

	//close(f.mqttMessageQueue)
	//close(f.mqttWorkerQuit)
}

func (f *Flow) NewMessage(ctx context.Context) (err error) {

	f.Lock()
	if f.isRunning {
		err = errors.New("flow is running")
		f.Unlock()
		return
	}

	defer func() {
		f.Lock()
		f.isRunning = false
		f.Unlock()
	}()

	// circular dependency search
	if ctx, err = f.defineCircularConnection(ctx); err != nil {
		return
	}

	f.isRunning = true
	f.Unlock()

	var _element *FlowElement

	// find message handler
	// ------------------------------------------------
	for _, element := range f.FlowElements {
		if element.Prototype == nil {
			continue
		}

		if element.Model.PrototypeType != "MessageHandler" {
			continue
		}

		_element = element
		break
	}

	if _element == nil {
		err = errors.New("message handler not found")
		return
	}

	// ------------------------------------------------
	getNextElements := func(element *FlowElement, isScripted, isTrue bool) (elements []*FlowElement) {
		// each connections
		for _, conn := range f.Connections {
			if conn.ElementFrom != element.Model.Uuid || conn.ElementTo == element.Model.Uuid {
				continue
			}

			for _, element := range f.FlowElements {
				if conn.ElementTo != element.Model.Uuid {
					continue
				}

				if isScripted {
					if conn.Direction == "true" {
						if !isTrue {
							continue
						}
					} else if conn.Direction == "false" {
						if isTrue {
							continue
						}
					}
				}

				elements = append(elements, element)
			}
		}

		return
	}

	if msg, ok := ctx.Value("msg").(*Message); ok {
		f.SetMessage(msg)
	}

	var runElement func(*FlowElement)
	runElement = func(element *FlowElement) {

		if err = ctx.Err(); err != nil {
			return
		}

		var ok, isScripted bool
		isScripted = element.ScriptEngine != nil

		childCtx, _ := context.WithCancel(ctx)

		if ctx, ok, err = element.Run(childCtx); err != nil {
			//log.Error(err.Error())
			return
		}

		// send message to linked flow
		if element.Model.PrototypeType == "Flow" && element.Model.FlowLink != nil {
			if flow, ok := f.workflow.Flows[*element.Model.FlowLink]; ok {
				childCtx, _ := context.WithCancel(ctx)
				childCtx = context.WithValue(childCtx, "msg", f.message)
				if err = flow.NewMessage(childCtx); err != nil {
					return
				}
				f.SetMessage(flow.message)
			}
		}

		elements := getNextElements(element, isScripted, ok)
		for _, e := range elements {
			runElement(e)
		}
	}

	runElement(_element)

	return
}

// ------------------------------------------------
// Workers
// ------------------------------------------------

func (f *Flow) InitWorkers() (err error) {

	for _, worker := range f.Model.Workers {
		if err = f.AddWorker(worker); err != nil {
			return
		}
	}

	return
}

func (f *Flow) AddWorker(model *m.Worker) (err error) {

	log.Infof("Add worker: \"%s\"", model.Name)

	f.Lock()
	if _, ok := f.Workers[model.Id]; ok {
		f.Unlock()
		return
	}
	f.Unlock()

	if len(f.FlowElements) == 0 {
		err = errors.New("no flow elements")
		return
	}

	// get device
	// ------------------------------------------------
	var devices []*m.Device
	if !model.DeviceAction.Device.IsGroup {
		devices = append(devices, model.DeviceAction.Device)
	} else {
		// значит тут группа устройств
		for _, child := range model.DeviceAction.Device.Devices {
			if child.Status != "enabled" {
				continue
			}

			//if child.Address == nil {
			//	continue
			//}

			device := &m.Device{
				Id:         child.Id,
				Name:       child.Name,
				Properties: child.Properties,
				Type:       model.DeviceAction.Device.Type,
				Device:     &m.Device{Id: model.DeviceAction.Device.Id},
			}

			//*device = *model.DeviceAction.Device
			//device.Id = child.Id
			//device.Name = child.Name
			//device.Address = new(int)
			//*device.Address = *child.Address
			//device.Device = &m.Device{Id: model.DeviceAction.Device.Id}
			//device.Tty = child.Tty
			//device.Sleep = model.DeviceAction.Device.Sleep
			devices = append(devices, device)
		}
	}

	// get node
	// ------------------------------------------------
	nodeId := model.DeviceAction.Device.Node.Id
	var ok bool
	if f.Node, ok = f.core.safeGetOrAddNode(nodeId); !ok {
		err = fmt.Errorf("node %d not found", nodeId)
		log.Error(err.Error())
		return
	}

	// generate new worker
	worker := NewWorker(model, f, f.cron)

	// add devices to worker
	// ------------------------------------------------
	for _, device := range devices {

		var action *Action
		if action, err = NewAction(device, model.DeviceAction, f.Node, f, f.scriptService); err != nil {
			log.Error(err.Error())
			continue
		}

		worker.AddAction(action)
	}

	f.Workers[model.Id] = worker
	f.Workers[model.Id].Start()

	return
}

func (f *Flow) UpdateWorker(worker *m.Worker) (err error) {

	f.Lock()
	if _, ok := f.Workers[worker.Id]; !ok {
		err = fmt.Errorf("worker id:%d not found", worker.Id)
	}
	f.Unlock()

	if err = f.RemoveWorker(worker); err != nil {
		log.Warnf("error: %s", err.Error())
	}

	if err = f.AddWorker(worker); err != nil {
		log.Warnf("error: %s", err.Error())
	}

	return
}

func (f *Flow) RemoveWorker(worker *m.Worker) (err error) {

	log.Infof("Remove worker: \"%s\"", worker.Name)

	if _, ok := f.Workers[worker.Id]; !ok {
		err = fmt.Errorf("worker id:%d not found", worker.Id)
		return
	}

	// stop cron task
	f.Workers[worker.Id].Stop()

	// delete worker
	delete(f.Workers, worker.Id)

	return
}

func (f *Flow) NewScript(s ...*m.Script) (engine *scripts.Engine, err error) {

	var model *m.Script
	if len(s) == 0 {
		model = &m.Script{
			Lang: ScriptLangJavascript,
		}
	} else {
		model = s[0]
	}

	if engine, err = f.workflow.NewScript(model); err != nil {
		return
	}

	engine.PushStruct("Flow", &FlowBind{flow: f})
	engine.PushStruct("Workflow", &WorkflowBind{wf: f.workflow})

	// message
	engine.PushStruct("message", f.message)

	return
}

func (f *Flow) defineCircularConnection(ctx context.Context) (newCtx context.Context, err error) {

	if v := ctx.Value("parents"); v != nil {
		if parents, ok := v.([]int64); ok {
			var exist bool
			for _, parentId := range parents {
				if parentId == f.Model.Id {
					exist = true
				}
			}

			if exist {
				depends := fmt.Sprintf("%d", parents[0])
				for _, parentId := range parents[1:] {
					depends = fmt.Sprintf("%s -> %d", depends, parentId)
				}
				err = fmt.Errorf("circular relationship detected: %s -> flow(%d)", depends, f.Model.Id)
				return
			}

			parents = append(parents, f.Model.Id)
			newCtx = context.WithValue(ctx, "parents", parents)

		} else {
			err = fmt.Errorf("bad parent context value: parents(%v)", parents)
		}

		return
	}

	newCtx = context.WithValue(ctx, "parents", []int64{f.Model.Id})

	return
}

func (f *Flow) mqttOnPublish(client *mqtt.Client, msg mqtt.Message) {

	message := NewMessage()
	message.SetVar("mqtt_payload", string(msg.Payload))
	message.SetVar("mqtt_topic", msg.Topic)
	message.SetVar("mqtt_qos", msg.Qos)
	message.SetVar("mqtt_duplicate", msg.Dup)
	message.Mqtt = true

	f.mqttMessageQueue <- message
}

func (f *Flow) mqttNewMessage(message *Message) {

	// create context
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	ctx = context.WithValue(ctx, "msg", message)

	done := make(chan struct{})
	go func() {
		if err := f.NewMessage(ctx); err != nil {
			log.Errorf("flow '%v' end with error: '%+v'", f.Model.Name, err.Error())
		}

		if ctx.Err() != nil {
			log.Errorf("flow '%v' end with error: '%+v'", f.Model.Name, ctx.Err())
		}

		done <- struct{}{}
	}()

	select {
	case <-done:
		close(done)
	case <-ctx.Done():

	}
}

func (f *Flow) mqttMessageWorker() {

	for {

		f.Lock()
		if f.isRunning {
			time.Sleep(time.Millisecond * 500)
			f.Unlock()
			continue
		}
		f.Unlock()

		select {
		case <-f.mqttWorkerQuit:
			return

		case message := <-f.mqttMessageQueue:
			f.mqttNewMessage(message)
		}
	}
}

func (f *Flow) SetMessage(msg *Message) {
	f.Lock()
	f.message.Update(msg)
	f.Unlock()
}

func (f *Flow) GetMessage() *Message {
	f.Lock()
	defer f.Unlock()
	return f.message.Copy()
}
