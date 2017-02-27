package core

import (
	"sync"
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/api/models"
	"github.com/astaxie/beego/orm"
	"github.com/e154/smart-home/api/scripts"
)

func NewWorkflow(model *models.Workflow, nodes map[int64]*models.Node) (workflow *Workflow) {

	workflow = &Workflow{
		model: model,
		Nodes: nodes,
		Flows: make(map[int64]*Flow),
		mutex: &sync.Mutex{},
	}

	workflow.UpdateScenario()

	return
}

type Workflow struct {
	model   	*models.Workflow
	Nodes   	map[int64]*models.Node
	mutex   	*sync.Mutex
	Flows   	map[int64]*Flow
}

func (wf *Workflow) Run() (err error) {

	err = wf.InitFlows()

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
func (wf *Workflow) InitFlows() (err error) {

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

	wf.mutex.Lock()
	if _, ok := wf.Flows[flow.Id]; ok {
		return
	}
	wf.mutex.Unlock()

	var model *Flow
	if model, err = NewFlow(flow, wf); err != nil {
		return
	}

	wf.mutex.Lock()
	wf.Flows[flow.Id] = model
	wf.mutex.Unlock()


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

	wf.mutex.Lock()
	defer wf.mutex.Unlock()

	if _, ok := wf.Flows[flow.Id]; !ok {
		return
	}

	wf.Flows[flow.Id].Remove()
	delete(wf.Flows, flow.Id)

	return
}

func (wf *Workflow) UpdateScenario() (err error) {

	// load related scenario and his scripts
	o := orm.NewOrm()
	if wf.model.Scenario != nil {

		log.Infof("Workflow '%s': update scenario", wf.model.Name)

		var old_scenario *models.Scenario
		old_scenario = wf.model.Scenario

		if _, err = o.LoadRelated(wf.model, "Scenario"); err != nil {
			log.Errorf("on update scenario, message: %s", err.Error())
			return
		}

		if _, err = o.LoadRelated(wf.model.Scenario, "Scripts"); err != nil {
			log.Errorf("on update scenario, message: %s", err.Error())
			return
		}

		wf.runScenarioScript(old_scenario, "on_exit")
		wf.runScenarioScript(wf.model.Scenario, "on_enter")
	}

	return
}

func (wf *Workflow) runScenarioScript(scenario *models.Scenario, state string) (err error) {

	var _script *scripts.Engine
	for _, scenario_script := range scenario.Scripts {
		if scenario_script.State != state {
			continue
		}

		// load script
		o := orm.NewOrm()
		if _, err = o.LoadRelated(scenario_script, "Script"); err != nil {
			log.Errorf("compile script %d, message: %s", scenario_script.Script.Id, err.Error())
			return
		}

		// compile script
		if _script, err = scripts.New(scenario_script.Script); err != nil {
			log.Errorf("compile script %d, message: %s", scenario_script.Script.Id, err.Error())
			continue
		}

		// do script
		if _, err = _script.Do(); err != nil {
			log.Errorf("on run script %s scenario, message: %s", state, err.Error())
		}
	}

	return
}