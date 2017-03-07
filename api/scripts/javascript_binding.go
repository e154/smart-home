package scripts

import (
	"github.com/e154/smart-home/api/models"
	r "github.com/e154/smart-home/lib/rpc"
	"github.com/e154/smart-home/api/log"
	"time"
	"sync"
)

type JavascriptBinding struct {
	mu	sync.Mutex
	pool	map[string]interface{}
}

// Logging
// method: log
//
func (m *JavascriptBinding) Log(s string, v ...interface{}) {

	switch s {
	case "info":
		log.Info(v)
	case "warn":
		log.Warn(v)
	case "debug":
		log.Debug(v)
	case "error":
		log.Error(v)
	default:
		log.Debug(v)
	}
}

// Get new device instance
// method: new_device
//
func (m *JavascriptBinding) New_device() *models.Device {
	return  &models.Device{}
}

//TODO check
// Get new flow instance
// method: new_flow
//
func (m *JavascriptBinding) New_flow() *models.Flow {
	return  &models.Flow{}
}

// Get new node instance
// method: new_node
//
func (m *JavascriptBinding) New_node() *models.Node {
	return  &models.Node{}
}

// send date to node
// method: node_send
//
func (m *JavascriptBinding) Node_send(protocol string, node *models.Node, device *models.Device, command []byte,) (result r.Result) {

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

//
func initJsBinds(j *Javascript) (jsBinds *JavascriptBinding) {

	// print
	j.ctx.PushGlobalGoFunction("print", func(a ...interface{}){
		j.engine.Print(a...)
	})

	// global main object
	jsBinds = &JavascriptBinding{
		pool: make(map[string]interface{}),
	}
	j.PushStruct("smart", jsBinds)
	j.PushStruct("device", &models.Device{})
	j.PushFunction("to_time", func(i int64) time.Duration {
		return time.Duration(i)
	})

	// etc
	j.ctx.PevalString(`helper = {}`)
	j.ctx.PevalString(`helper.hex2arr = function (hexString) {
   var result = [];
   while (hexString.length >= 2) {
       result.push(parseInt(hexString.substring(0, 2), 16));
       hexString = hexString.substring(2, hexString.length);
   }
   return result;
}`)

	return
}