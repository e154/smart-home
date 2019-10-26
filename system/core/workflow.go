package core

import (
	"errors"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	cr "github.com/e154/smart-home/system/cron"
	"github.com/e154/smart-home/system/scripts"
	"sync"
	"time"
)

type Workflow struct {
	Storage
	model *m.Workflow
	sync.Mutex
	adaptors        *adaptors.Adaptors
	scripts         *scripts.ScriptService
	Flows           map[int64]*Flow
	engine          *scripts.Engine
	cron            *cr.Cron
	core            *Core
	nextScenario    *m.WorkflowScenario
	isRuning        bool
	scenarioEntered bool
}

func NewWorkflow(model *m.Workflow,
	adaptors *adaptors.Adaptors,
	scripts *scripts.ScriptService,
	cron *cr.Cron,
	core *Core) (workflow *Workflow) {

	workflow = &Workflow{
		model:    model,
		adaptors: adaptors,
		scripts:  scripts,
		Flows:    make(map[int64]*Flow),
		cron:     cron,
		core:     core,
	}

	workflow.pull = make(map[string]interface{})

	return
}

func (wf *Workflow) Run() (err error) {

	if wf.isRuning {
		return
	}

	wf.isRuning = true

	defer func() {
		if err != nil {
			wf.isRuning = false
		}
	}()

	log.Infof("Run workflow '%v'", wf.model.Name)

	if err = wf.runScripts(); err != nil {
		return
	}

	if err = wf.enterScenario(); err != nil {
		return
	}

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

	wf.isRuning = false

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

	var flows []*m.Flow
	if flows, err = wf.adaptors.Flow.GetAllEnabledByWorkflow(wf.model.Id); err != nil {
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

	wf.Lock()
	if _, ok := wf.Flows[flow.Id]; ok {
		return
	}
	wf.Unlock()

	var model *Flow
	if model, err = NewFlow(flow, wf, wf.adaptors, wf.scripts, wf.cron, wf.core); err != nil {
		log.Error(err.Error())
		return
	}

	wf.Lock()
	wf.Flows[flow.Id] = model
	wf.Unlock()

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

	//wf.Lock()
	//defer wf.Unlock()

	if _, ok := wf.Flows[flow.Id]; !ok {
		return
	}

	wf.Flows[flow.Id].Remove()
	delete(wf.Flows, flow.Id)

	return
}

func (wf *Workflow) GetFLow(flowId int64) (flow *Flow, err error) {

	log.Infof("GetFLow: id(%v)", flowId)

	if _, ok := wf.Flows[flowId]; !ok {
		err = errors.New("not found")
		return
	}

	flow = wf.Flows[flowId]

	return
}

// ------------------------------------------------
// Scenarios
// ------------------------------------------------

func (wf *Workflow) SetScenario(systemName string) (err error) {

	log.Infof("workflow(%s) set scenario '%s'", wf.model.Name, systemName)

	var scenario *m.WorkflowScenario
	for _, scenario = range wf.model.Scenarios {
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

func (wf *Workflow) enterScenario() (err error) {

	if wf.model.Scenario == nil {
		return
	}

	log.Infof("Workflow '%s', scenario '%s'", wf.model.Name, wf.model.Scenario.Name)

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
	log.Infof("run scenario: %s", f)

	for _, scenarioScript := range wf.model.Scenario.Scripts {
		if err = wf.engine.EvalString(scenarioScript.Compiled); err != nil {
			log.Errorf("eval script(%d), message: %s", scenarioScript.Id, err.Error())
			continue
		}

		if _, err = wf.engine.DoCustom(f); err != nil {
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

	javascript := wf.engine.Get().(*scripts.Javascript)
	ctx := javascript.Ctx()
	if b := ctx.GetGlobalString("IC"); !b {
		return
	}
	ctx.PushObject()
	ctx.PushGoFunction(func() *WorkflowBind {
		return &WorkflowBind{wf: wf}
	})
	ctx.PutPropString(-3, "Workflow")
	ctx.Pop()

	for _, wfScript := range wf.model.Scripts {
		if err = wf.engine.EvalString(wfScript.Compiled); err != nil {
			log.Errorf(err.Error())
		}
	}

	return
}

func (wf *Workflow) NewScript(model *m.Script) (engine *scripts.Engine, err error) {
	engine, err = wf.scripts.NewEngine(model)
	return
}
