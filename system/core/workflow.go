package core

import (
	"sync"
	m "github.com/e154/smart-home/models"
)

type Workflow struct {
	Storage
	model *m.Workflow
	sync.Mutex
	//Flows   map[int64]*Flow
}

func NewWorkflow(model *m.Workflow) (workflow *Workflow) {

	workflow = &Workflow{
		model: model,
		//Flows: make(map[int64]*Flow),
	}

	workflow.pull = make(map[string]interface{})

	return
}

func (wf *Workflow) Run() (err error) {

	//wf.enterScenario()

	//wf.runScripts()

	//err = wf.initFlows()

	if err != nil {
		return
	}

	return
}

func (wf *Workflow) Stop() (err error) {

	//for _, flow := range wf.Flows {
	//	wf.RemoveFlow(flow.Model)
	//}

	return
}

func (wf *Workflow) Restart() (err error) {

	wf.Stop()
	err = wf.Run()

	return
}