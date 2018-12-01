package core

// Javascript Binding
//
// node
//	 .send()
//	 .name()
//	 .description()
//	 .ip()
//	 .status()
//
type NodeBind struct {
	node *Node
}

func (n *NodeBind) Send(device *DeviceBind, command []byte) (err error) {
	return n.node.Send(device.model, command)
}

func (n *NodeBind) IsConnected() bool {
	return n.node.IsConnected
}

//func (n *NodeBind) Smartbus(device *DeviceBind, returnResult bool, command []byte) SmartbusResult {
//	return n.node.Smartbus(device.model, returnResult, command)
//}
//
//func (n *NodeBind) Modbus(device *DeviceBind, returnResult bool, command []byte) SmartbusResult {
//	return n.node.Modbus(device.model, returnResult, command)
//}
//
func (n *NodeBind) Name() string {
	return n.node.Name
}

func (n *NodeBind) Ip() string {
	return n.node.Ip
}

func (n *NodeBind) Status() string {
	return n.node.Status
}

func (n *NodeBind) Description() string {
	return n.node.Description
}
