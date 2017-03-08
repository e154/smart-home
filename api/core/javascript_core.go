package core

type JavascriptCore struct {
	workflow *Workflow
}

func (j *JavascriptCore) SetVariable(key string, value interface{}) {
	j.workflow.SetVariable(key, value)
}

func (j *JavascriptCore) GetVariable(key string) interface{} {
	return j.workflow.GetVariable(key)
}