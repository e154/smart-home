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

func (n *NodeBind) Smartbus(device *DeviceBind, return_result bool, command []byte) SmartbusResult {
	return n.node.Smartbus(device.model, return_result, command)
}

func (n *NodeBind) Modbus(device *DeviceBind, return_result bool, command []byte) SmartbusResult {
	return n.node.Modbus(device.model, return_result, command)
}

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