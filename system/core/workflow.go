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
	"errors"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	cr "github.com/e154/smart-home/system/cron"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/telemetry"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"sync"
	"time"
)

type Workflow struct {
	Storage
	sync.Mutex
	model           *m.Workflow
	adaptors        *adaptors.Adaptors
	scripts         *scripts.ScriptService
	Flows           map[int64]*Flow
	engine          *scripts.Engine
	cron            *cr.Cron
	core            *Core
	mqtt            *mqtt.Mqtt
	nextScenario    *m.WorkflowScenario
	isRunning       bool
	scenarioEntered bool
	telemetry       telemetry.ITelemetry
	zigbee2mqtt     *zigbee2mqtt.Zigbee2mqtt
	metric          *metrics.MetricManager
}

func NewWorkflow(model *m.Workflow,
	adaptors *adaptors.Adaptors,
	scripts *scripts.ScriptService,
	cron *cr.Cron,
	core *Core,
	mqtt *mqtt.Mqtt,
	telemetry telemetry.ITelemetry,
	zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt,
	metric *metrics.MetricManager) (workflow *Workflow) {

	workflow = &Workflow{
		Storage:     NewStorage(),
		model:       model,
		adaptors:    adaptors,
		scripts:     scripts,
		Flows:       make(map[int64]*Flow),
		cron:        cron,
		core:        core,
		mqtt:        mqtt,
		telemetry:   telemetry,
		zigbee2mqtt: zigbee2mqtt,
		metric:      metric,
	}

	return
}

func (wf *Workflow) Run() (err error) {

	if wf.model == nil {
		err = errors.New("workflow model is nil")
		return
	}

	if wf.model.Scenario == nil {
		err = errors.New("workflow scenario is nil")
		return
	}

	if wf.isRunning {
		return
	}

	wf.isRunning = true

	defer func() {
		if err != nil {
			wf.isRunning = false
		}
	}()

	log.Infof("Run workflow '%v'", wf.model.Name)

	if err = wf.runScripts(); err != nil {
		//return
	}

	if err = wf.enterScenario(); err != nil {
		return
	}

	wf.telemetry.BroadcastOne(telemetry.WorkflowScenario{WorkflowId: wf.model.Id, ScenarioId: wf.model.Scenario.Id})

	if err = wf.initFlows(); err != nil {
		return
	}

	return
}

func (wf *Workflow) Stop() (err error) {

	for _, flow := range wf.Flows {
		wf.RemoveFlow(flow.Model)
	}

	err = wf.exitScenario()

	wf.isRunning = false

	return
}

func (wf *Workflow) Restart() (err error) {

	if err = wf.Stop(); err != nil {
		return
	}

	err = wf.Run()

	return
}

// ------------------------------------------------
// Flows
// ------------------------------------------------

// получаем все связанные процессы
func (wf *Workflow) initFlows() (err error) {

	log.Infof("Get flows")

	wf.Lock()
	workflowId := wf.model.Id
	wf.Unlock()

	var flows []*m.Flow
	if flows, err = wf.adaptors.Flow.GetAllEnabledByWorkflow(workflowId); err != nil {
		return
	}

	for _, flow := range flows {
		if err = wf.AddFlow(flow); err != nil {
			log.Error(err.Error())
		}
	}

	return
}

// Flow должен быть полный:
// с Connections
// с FlowElements
// с Cursor
// с Workers
func (wf *Workflow) AddFlow(flow *m.Flow) (err error) {

	if flow.Status != "enabled" {
		return
	}

	log.Infof("Add flow: '%s'", flow.Name)

	if _, ok := wf.safeGetFlow(flow.Id); ok {
		return
	}

	var model *Flow
	if model, err = NewFlow(flow, wf, wf.adaptors, wf.scripts, wf.cron, wf.core, wf.mqtt, wf.zigbee2mqtt); err != nil {
		log.Error(err.Error())
		return
	}

	wf.safeUpdateFlowMap(flow.Id, model)

	return
}

func (wf *Workflow) UpdateFlow(flow *m.Flow) (err error) {

	if err = wf.RemoveFlow(flow); err != nil {
		return
	}

	err = wf.AddFlow(flow)

	return
}

func (wf *Workflow) RemoveFlow(flow *m.Flow) (err error) {

	log.Infof("RemoveFlow: Name(%v)", flow.Name)

	f, ok := wf.safeGetFlow(flow.Id)
	if !ok {
		return
	}

	f.Remove()

	delete(wf.Flows, flow.Id)

	return
}

func (wf *Workflow) GetFLow(flowId int64) (flow *Flow, err error) {

	log.Infof("GetFLow: id(%v)", flowId)

	var ok bool
	if flow, ok = wf.safeGetFlow(flowId); !ok {
		err = errors.New("not found")
		return
	}

	return
}

// ------------------------------------------------
// Scenarios
// ------------------------------------------------

func (wf *Workflow) SetScenario(systemName string) (err error) {

	wf.Lock()
	name := wf.model.Name
	scenarios := wf.model.Scenarios
	wf.Unlock()

	log.Infof("workflow(%s) set scenario '%s'", name, systemName)

	var scenario *m.WorkflowScenario
	for _, scenario = range scenarios {
		if scenario.SystemName != systemName {
			continue
		}

		workflow := *wf.model
		workflow.Scenario = scenario

		if err = wf.adaptors.Workflow.Update(&workflow); err != nil {
			return
		}

		wf.nextScenario = scenario

		break
	}

	if !wf.scenarioEntered {
		return
	}

	wg := sync.WaitGroup{}

	wf.Lock()
	wg.Add(len(wf.Flows))
	wf.Unlock()

	go func() {
		for _, flow := range wf.Flows {
			wf.RemoveFlow(flow.Model)
			wg.Done()
		}
	}()

	wg.Wait()

	wf.UpdateScenario()

	return
}

func (wf *Workflow) enterScenario() (err error) {

	wf.Lock()
	scenario := wf.model.Scenario
	wf.Unlock()

	if scenario == nil {
		return
	}

	log.Infof("Workflow '%s', scenario '%s'", wf.model.Name, scenario.Name)

	go wf.metric.Update(metrics.WorkflowUpdateScenario{
		Id:         wf.model.Id,
		ScenarioId: wf.model.Scenario.Id,
	})

	if err = wf.runScenarioScripts("on_enter"); err != nil {
		return
	}

	wf.scenarioEntered = true

	time.Sleep(time.Second)

	if wf.nextScenario == nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(wf.Flows))

	go func() {
		for _, flow := range wf.Flows {
			wf.RemoveFlow(flow.Model)
			wg.Done()
		}
	}()

	wg.Wait()

	wf.UpdateScenario()

	return
}

func (wf *Workflow) exitScenario() (err error) {

	if wf.model.Scenario == nil {
		wf.scenarioEntered = false
		return
	}

	if wf.model.Scenario != nil {
		err = wf.runScenarioScripts("on_exit")
	}

	wf.scenarioEntered = false

	return
}

func (wf *Workflow) UpdateScenario() (err error) {

	// get workflow from base
	var model *m.Workflow
	if model, err = wf.adaptors.Workflow.GetById(wf.model.Id); err != nil {
		return
	}

	// exit if scenario is loaded
	if wf.model.Scenario == nil || wf.model.Scenario.SystemName == model.Scenario.SystemName {
		return
	}

	log.Infof("Workflow '%s' change scenario to '%s'", wf.model.Name, model.Scenario.Name)

	if err = wf.Stop(); err != nil {
		return
	}

	*wf.model = *model

	wf.nextScenario = nil
	err = wf.Run()

	return
}

func (wf *Workflow) runScenarioScripts(f string) (err error) {

	for _, scenarioScript := range wf.model.Scenario.Scripts {
		if _, err = wf.engine.EvalScript(scenarioScript); err != nil {
			log.Errorf("eval script(%d), message: %s", scenarioScript.Id, err.Error())
			continue
		}

		if _, err = wf.engine.AssertFunction(f); err != nil {
			log.Errorf("on run script %s scenario, message: %s", f, err.Error())
		}
	}

	return
}

// ------------------------------------------------
// Scripts
// ------------------------------------------------

func (wf *Workflow) runScripts() (err error) {

	dummy := &m.Script{
		Lang: common.ScriptLangJavascript,
	}
	if wf.engine, err = wf.scripts.NewEngine(dummy); err != nil {
		return
	}

	wf.engine.PushStruct("Workflow", &WorkflowBind{wf: wf})

	for _, wfScript := range wf.model.Scripts {
		if _, err = wf.engine.EvalScript(wfScript); err != nil {
			log.Errorf(err.Error())
		}
	}

	return
}

func (wf *Workflow) NewScript(model *m.Script) (engine *scripts.Engine, err error) {
	engine, err = wf.scripts.NewEngine(model)
	return
}

// ------------------------------------------------
// Workers
// ------------------------------------------------

func (wf *Workflow) UpdateWorker(_worker *m.Worker) (err error) {

	for _, flow := range wf.Flows {
		for _, worker := range flow.Workers {
			if worker.Model.Id == _worker.Id {
				if err = flow.UpdateWorker(_worker); err != nil {
					log.Error(err.Error())
				}
				break
			}
		}
	}

	return
}

func (wf *Workflow) RemoveWorker(_worker *m.Worker) (err error) {

	wf.Lock()
	defer wf.Unlock()

	for _, flow := range wf.Flows {
		for _, worker := range flow.Workers {
			if worker.Model.Id == _worker.Id {
				if err = flow.RemoveWorker(_worker); err != nil {
					log.Error(err.Error())
				}
				break
			}
		}
	}
	return
}

func (wf *Workflow) DoWorker(model *m.Worker) (err error) {

	for _, flow := range wf.Flows {
		if worker, ok := flow.Workers[model.Id]; ok {
			worker.Do()
			break
		}
	}

	return
}

// ------------------------------------------------
// safe methods
// ------------------------------------------------

func (b *Workflow) safeIsRunning() bool {
	b.Lock()
	defer b.Unlock()
	return b.isRunning
}

func (b *Workflow) safeSetIsRunning(v bool) {
	b.Lock()
	b.isRunning = v
	b.Unlock()
}

func (c *Workflow) safeGetFlow(k int64) (w *Flow, ok bool) {
	c.Lock()
	defer c.Unlock()
	w, ok = c.Flows[k]
	return
}

func (c *Workflow) safeUpdateFlowMap(k int64, w *Flow) {
	c.Lock()
	c.Flows[k] = w
	c.Unlock()
}
