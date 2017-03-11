package scripts

import (
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/common"
)

type Node struct {

}

// send date to node
// method: node.send
//
func (m *Node) Send(protocol string, node *models.Node, device *models.Device, command []byte,) (result common.Result) {

	var request *common.Request
	request = &common.Request{}
	request.Baud = device.Baud
	request.Result = true
	request.Device = device.Tty
	request.Timeout = device.Timeout
	request.StopBits = int(device.StopBite)
	request.Sleep = device.Sleep

	// set command
	request.Command = command

	// send command to node
	result = node.Send(protocol, request)

	return
}
