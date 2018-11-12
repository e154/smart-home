package core

// Javascript Binding
//
//IC.Workflow()
//	 .getName()
//	 .getDescription()
//	 .setVar(string, interface)
//	 .getVar(string)
//	 .getScenario() string
//	 .GetScenarioName() string
//	 .setScenario(string)
//
type WorkflowBind struct {
	wf *Workflow
}

func (w *WorkflowBind) GetName() string {
	return w.wf.model.Name
}

func (w *WorkflowBind) GetDescription() string {
	return w.wf.model.Description
}

func (w *WorkflowBind) SetVar(key string, value interface{}) {
	w.wf.SetVar(key, value)
}

func (w *WorkflowBind) GetVar(key string) interface{} {
	return w.wf.GetVar(key)
}

func (w *WorkflowBind) GetScenario() string {
	return w.wf.model.Scenario.SystemName
}

func (w *WorkflowBind) GetScenarioName() string {
	return w.wf.model.Scenario.Name
}

func (w *WorkflowBind) SetScenario(system_name string) {
	w.wf.SetScenario(system_name)
}