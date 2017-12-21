package core

import (
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/common"
)

// Javascript Binding
//
// node
//	 .Send()
//
type NodeBind struct {
	node *models.Node
}

func (n *NodeBind) Send(protocol string, device *models.Device, return_result bool, command []byte) common.Result {
	return n.node.Send(protocol, device, return_result, command)
}