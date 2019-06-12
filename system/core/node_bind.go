package core

// Javascript Binding
//
// node
//	 .name()
//	 .ip()
//	 .status()
//	 .stat()
//	 .description()
//	 .isConnected()
//
type NodeBind struct {
	node *Node
}

//func (n *NodeBind) Send(device *DeviceBind, command []byte) NodeBindResult {
//	return n.node.Send(device.model, command)
//}

func (n *NodeBind) IsConnected() bool {
	return n.node.IsConnected
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

func (n *NodeBind) Stat() *NodeStatModel {
	return n.node.stat
}

func (n *NodeBind) Description() string {
	return n.node.Description
}
