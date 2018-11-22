package core

// Javascript Binding
//
// IC.Action()
//	 .device()
//	 .node()
//
type ActionBind struct {
	action *Action
}

func (a *ActionBind) Device() *DeviceBind {
	return &DeviceBind{model: a.action.GetDevice()}
}

func (a *ActionBind) Node() *NodeBind {
	return &NodeBind{node: a.action.GetNode()}
}
