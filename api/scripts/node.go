package scripts

import (
	"github.com/e154/smart-home/api/models"
	r "github.com/e154/smart-home/lib/rpc"
)

type Node struct {

}

// send date to node
// method: node.send
//
func (m *Node) Send(protocol string, node *models.Node, device *models.Device, command []byte,) (result r.Result) {

	var request *r.Request
	request = &r.Request{}
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
