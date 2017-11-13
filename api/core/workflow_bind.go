package core

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