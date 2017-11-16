package core

import (
	"sync"
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/scripts"
)

func NewWorkflow(model *models.Workflow) (workflow *Workflow) {

	workflow = &Workflow{
		model: model,
		Flows: make(map[int64]*Flow),
	}

	workflow.pull = make(map[string]interface{})

	return
}

type Workflow struct {
	Storage
	model   *models.Workflow
	sync.Mutex
	Flows   map[int64]*Flow
}

func (wf *Workflow) Run() (err error) {

	wf.enterScenario()

	wf.runScripts()

	err = wf.initFlows()

	if err != nil {
		return
	}

	return
}

func (wf *Workflow) Stop() (err error) {

	for _, flow := range wf.Flows {
		wf.RemoveFlow(flow.Model)
	}

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

	var flows []*models.Flow
	if flows, err = wf.model.GetAllEnabledFlows(); err != nil {
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
func (wf *Workflow) AddFlow(flow *models.Flow) (err error) {

	if flow.Status != "enabled" {
		return
	}

	log.Info("Add flow:", flow.Name)

	wf.Lock()
	if _, ok := wf.Flows[flow.Id]; ok {
		return
	}
	wf.Unlock()

	var model *Flow
	if model, err = NewFlow(flow, wf); err != nil {
		return
	}

	wf.Lock()
	wf.Flows[flow.Id] = model
	wf.Unlock()


	return
}

func (wf *Workflow) UpdateFlow(flow *models.Flow) (err error) {

	err = wf.RemoveFlow(flow)
	if err != nil {
		return
	}

	err = wf.AddFlow(flow)

	return
}

func (wf *Workflow) RemoveFlow(flow *models.Flow) (err error) {

	log.Info("Remove flow:", flow.Name)

	wf.Lock()
	defer wf.Unlock()

	if _, ok := wf.Flows[flow.Id]; !ok {
		return
	}

	wf.Flows[flow.Id].Remove()
	delete(wf.Flows, flow.Id)

	return
}

// ------------------------------------------------
// Scenarios
// ------------------------------------------------

func (wf *Workflow) SetScenario(system_name string) (err error) {

	if _, err = wf.model.GetScenarios(); err != nil {
		return
	}

	var scenario *models.WorkflowScenario
	for _, scenario = range wf.model.Scenarios {
		if scenario.SystemName != system_name {
			continue
		}

		workflow := &models.Workflow{}
		*workflow = *wf.model
		workflow.Scenario = scenario

		if err = models.UpdateWorkflowById(workflow); err != nil {
			return
		}

		wf.UpdateScenario()
	}

	return
}

func (wf *Workflow) enterScenario() (err error) {

	if wf.model.Scenario == nil {
		return
	}

	log.Infof("Workflow '%s': enter scenario", wf.model.Name)

	if _, err = wf.model.GetScenario(); err != nil {
		log.Errorf("on update scenario, message: %s", err.Error())
		return
	}

	if _, err = wf.model.Scenario.GetScripts(); err != nil {
		log.Errorf("on update scenario, message: %s", err.Error())
		return
	}

	wf.runScenarioScripts(wf.model.Scenario, "on_enter")

	return
}

func (wf *Workflow) UpdateScenario() (err error) {

	// get workflow from base
	var model *models.Workflow
	if model, err = models.GetWorkflowById(wf.model.Id); err != nil {
		return
	}

	if _, err = model.GetScenario(); err != nil {
		log.Errorf("on update scenario, message: %s", err.Error())
		return
	}

	// exit if scenario is loaded
	if wf.model.Scenario.SystemName == model.Scenario.SystemName {
		return
	}

	log.Infof("Workflow '%s': update scenario", wf.model.Name)

	if wf.model.Scenario != nil {
		wf.runScenarioScripts(wf.model.Scenario, "on_exit")
	}

	*wf.model = *model

	if _, err = wf.model.GetScenario(); err != nil {
		log.Errorf("on update scenario, message: %s", err.Error())
		return
	}

	if _, err = wf.model.Scenario.GetScripts(); err != nil {
		log.Errorf("on update scenario, message: %s", err.Error())
		return
	}

	wf.Restart()

	return
}

func (wf *Workflow) runScenarioScripts(scenario *models.WorkflowScenario, f string) (err error) {

	var script *scripts.Engine
	for _, scenario_script := range scenario.Scripts {

		if script, err = wf.NewScript(scenario_script); err != nil {
			log.Errorf("compile script %d, message: %s", scenario_script.Id, err.Error())
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
	if _, err = wf.model.GetScripts(); err != nil {
		return
	}

	var script *scripts.Engine
	for _, scenario_script := range wf.model.Scripts {

		if script, err = wf.NewScript(scenario_script); err != nil {
			log.Errorf("compile script %d, message: %s", scenario_script.Id, err.Error())
		}

		if _, err = script.Do(); err != nil {
			log.Errorf("on run script %s", err.Error())
		}
	}

	return
}

func (wf *Workflow) NewScript(model *models.Script) (script *scripts.Engine, err error) {

	if script, err = scripts.New(model); err != nil {
		return
	}

	javascript := script.Get().(*scripts.Javascript)
	ctx := javascript.Ctx()
	if b := ctx.GetGlobalString("IC"); !b {
		return
	}
	ctx.PushObject()
	ctx.PushGoFunction(func() *WorkflowBind {
		return &WorkflowBind{wf:wf}
	})
	ctx.PutPropString(-3, "Workflow")
	ctx.Pop()

	return
}