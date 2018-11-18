package core

import (
	"sync"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/scripts"
	"errors"
)

type Workflow struct {
	Storage
	model *m.Workflow
	sync.Mutex
	adaptors *adaptors.Adaptors
	scripts  *scripts.ScriptService
	Flows    map[int64]*Flow
	engine   *scripts.Engine
}

func NewWorkflow(model *m.Workflow,
	adaptors *adaptors.Adaptors,
	scripts *scripts.ScriptService) (workflow *Workflow) {

	workflow = &Workflow{
		model:    model,
		adaptors: adaptors,
		scripts:  scripts,
		Flows:    make(map[int64]*Flow),
	}

	workflow.pull = make(map[string]interface{})

	return
}

func (wf *Workflow) Run() (err error) {

	wf.runScripts()

	wf.enterScenario()

	err = wf.initFlows()

	return
}

func (wf *Workflow) Stop() (err error) {

	for _, flow := range wf.Flows {
		wf.RemoveFlow(flow.Model)
	}

	err = wf.exitScenario()

	return
}

func (wf *Workflow) Restart() (err error) {

	wf.Stop()
	err = wf.Run()

	return
}

// ------------------------------------------------
// Flows
// ------------------------------------------------

// получаем все связанные процессы
func (wf *Workflow) initFlows() (err error) {

	var flows []*m.Flow
	if flows, err = wf.adaptors.Flow.GetAllEnabledByWorkflow(wf.model.Id); err != nil {
		return
	}

	for _, flow := range flows {
		wf.AddFlow(flow)
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

	log.Infof("Add flow: %s", flow.Name)

	wf.Lock()
	if _, ok := wf.Flows[flow.Id]; ok {
		return
	}
	wf.Unlock()

	var model *Flow
	if model, err = NewFlow(flow, wf, wf.adaptors, wf.scripts); err != nil {
		log.Error(err.Error())
		return
	}

	wf.Lock()
	wf.Flows[flow.Id] = model
	wf.Unlock()

	return
}

func (wf *Workflow) UpdateFlow(flow *m.Flow) (err error) {

	err = wf.RemoveFlow(flow)
	if err != nil {
		return
	}

	err = wf.AddFlow(flow)

	return
}

func (wf *Workflow) RemoveFlow(flow *m.Flow) (err error) {

	log.Infof("Remove flow: %s", flow.Name)

	wf.Lock()
	defer wf.Unlock()

	if _, ok := wf.Flows[flow.Id]; !ok {
		return
	}

	wf.Flows[flow.Id].Remove()
	delete(wf.Flows, flow.Id)

	return
}

func (wf *Workflow) GetFLow(flowId int64) (flow *Flow, err error) {

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

	var scenario *m.WorkflowScenario
	for _, scenario = range wf.model.Scenarios {
		if scenario.SystemName != systemName {
			continue
		}

		workflow := &m.Workflow{}
		*workflow = *wf.model
		workflow.Scenario = scenario

		if err = wf.adaptors.Workflow.Update(workflow); err != nil {
			return
		}

		wf.UpdateScenario()

		break
	}

	return
}

func (wf *Workflow) enterScenario() (err error) {

	if wf.model.Scenario == nil {
		return
	}

	log.Infof("Workflow '%s': enter scenario", wf.model.Name)

	err = wf.runScenarioScripts(wf.model.Scenario, "on_enter")

	return
}

func (wf *Workflow) exitScenario() (err error) {

	if wf.model.Scenario == nil {
		return
	}

	log.Infof("Workflow '%s': exit from scenario", wf.model.Name)

	err = wf.runScenarioScripts(wf.model.Scenario, "on_exit")

	return
}

func (wf *Workflow) UpdateScenario() (err error) {

	// get workflow from base
	var model *m.Workflow
	if model, err = wf.adaptors.Workflow.GetById(wf.model.Id); err != nil {
		return
	}

	// exit if scenario is loaded
	if wf.model.Scenario.SystemName == model.Scenario.SystemName {
		return
	}

	log.Infof("Workflow '%s': update scenario", wf.model.Name)

	if wf.model.Scenario != nil {
		if err = wf.runScenarioScripts(wf.model.Scenario, "on_exit"); err != nil {
			return
		}
	}

	*wf.model = *model

	err = wf.enterScenario()

	return
}

func (wf *Workflow) runScenarioScripts(scenario *m.WorkflowScenario, f string) (err error) {

	var script *scripts.Engine
	for _, scenarioScript := range scenario.Scripts {

		if script, err = wf.NewScript(scenarioScript); err != nil {
			log.Errorf("compile script %d, message: %s", scenarioScript.Id, err.Error())
		}

		if _, err = script.DoCustom(f); err != nil {
			log.Errorf("on run script %s scenario, message: %s", f, err.Error())
		}
	}

	return
}

// ------------------------------------------------
// Scripts
// ------------------------------------------------

func (wf *Workflow) runScripts() (err error) {

	var engine *scripts.Engine
	for _, scenarioScript := range wf.model.Scripts {
		if engine, err = wf.NewScript(scenarioScript); err != nil {
			continue
		}

		if _, err = engine.DoFull(); err != nil {
			log.Errorf("on run script %s", err.Error())
		}
	}

	return
}

func (wf *Workflow) NewScript(model *m.Script) (script *scripts.Engine, err error) {

	if script, err = wf.scripts.NewEngine(model); err != nil {
		return
	}

	javascript := script.Get().(*scripts.Javascript)
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

	return
}
