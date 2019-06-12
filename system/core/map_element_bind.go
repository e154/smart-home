package core

// Javascript Binding
//
// mapElement
//	.setState()
//	.getState()
//	.setOptions()
//	.getOptions()
//
type MapElementBind struct {
	element *MapElement
}

func (e *MapElementBind) SetState(name string) {
	e.element.SetState(name)
}

func (e *MapElementBind) GetState() interface{} {
	return e.element.State
}

func (e *MapElementBind) SetOptions(options interface{}) {
	e.element.SetOptions(options)
}

func (e *MapElementBind) GetOptions() interface{} {
	return e.element.GetOptions()
}
