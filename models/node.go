package models

import (
	"net"
	"sync"
	"time"
	"net/rpc"
	"fmt"
	"errors"
	"github.com/e154/smart-home/system/validation"
)

type SmartbusRequest struct {
	Line     string        `json:"line"`
	Device   string        `json:"device"`
	Baud     int           `json:"baud"`
	StopBits int64         `json:"stop_bits"`
	Sleep    int64         `json:"sleep"`
	Timeout  time.Duration `json:"timeout"`
	Command  []byte        `json:"command"`
	Result   bool          `json:"result"`
}

type SmartbusResult struct {
	Command   []byte `json:"command"`
	Device    string `json:"device"`
	Result    string `json:"result"`
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type Nodes []*Node

type Node struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name" valid:"MaxSize(254);Required"`
	Ip          string     `json:"ip" valid:"IP;Required"` // Must be a valid IPv4 address
	Port        int        `json:"port" valid:"Range(1, 65535);Required"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdateAt    *time.Time `json:"update_at"`
	sync.Mutex              //
	netConn     net.Conn    //
	errors      int64       //
	connStatus  string      //
	ch          chan string //
	rpcClient   *rpc.Client //
}

func NewNode() (node *Node) {
	node = &Node{
		ch: make(chan string),
	}
	return
}

func (d *Node) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}

func (n *Node) Start() {
	go n.run()
}

func (n *Node) run() {
	var quit, disconnect, connect bool

	for ; ; {

		select {
		case c := <-n.ch:
			switch c {
			case "quit":
				quit = true
			case "disconnect":
				disconnect = true
			case "connect":
				connect = true
			default:

			}

		default:

		}

		if quit {
			n.TcpClose()
			break
		}

		if n.errors > 5 {
			connect = true
		}

		if disconnect {
			n.TcpClose()
			connect = false
			disconnect = false
			n.SetConnectStatus("disconnected")
		}

		if connect {
			disconnect = false
			n.TcpClose()

			if _, err := n.RpcDial(); err == nil {
				n.errors = 0
				log.Warningf("Node dial tcp %s:%d ... ok", n.Ip, n.Port)
				connect = false
				n.SetConnectStatus("connected")
			} else {
				n.errors++
				if n.errors == 7 {
					log.Errorf("Node error %s", err.Error())
				}
				n.SetConnectStatus("error")
			}
		}

		time.Sleep(time.Second)
	}
}

func (n *Node) Connect() *Node {
	n.ch <- "connect"
	return n
}

func (n *Node) Disconnect() {
	n.ch <- "disconnect"
}

func (n *Node) Quit() {
	n.ch <- "quit"
}

func (n *Node) RpcDial() (*rpc.Client, error) {
	var err error
	if _, err = n.TcpDial(); err != nil {
		return nil, err
	}
	if n.rpcClient == nil {
		n.rpcClient = rpc.NewClient(n.netConn)
	}
	return n.rpcClient, err
}

func (n *Node) TcpDial() (net.Conn, error) {
	var err error
	if n.netConn == nil {
		n.netConn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", n.Ip, n.Port), time.Second*2)
		if err != nil {
			return nil, err
		}
	}
	//defer n.netConn.Close()
	return n.netConn, err
}

func (n *Node) TcpClose() {
	if n.netConn == nil {
		return
	}
	n.netConn.Close()
	n.netConn = nil
	n.rpcClient = nil
}

func (n *Node) GetConnectStatus() string {
	n.Lock()
	defer n.Unlock()

	if n.Status == "disabled" {
		n.connStatus = "disabled"
	}

	return n.connStatus
}

func (n *Node) SetConnectStatus(st string) {
	n.Lock()
	n.connStatus = st
	n.Unlock()

	//TODO fix call core telemetry
	//CorePtr().telemetry.Broadcast("nodes")
}

func (n *Node) Smartbus(device *Device, return_result bool, command []byte) (result SmartbusResult) {

	request := &SmartbusRequest{
		Baud:     device.Baud,
		Device:   device.Tty,
		Timeout:  device.Timeout,
		StopBits: device.StopBite,
		Sleep:    device.Sleep,
		Result:   return_result,
		Command:  command,
	}

	if err := n.RpcCall("Smartbus.Send", request, &result); err != nil {
		result.Error = err.Error()
	}

	return
}

//TODO update modbus method
func (n *Node) Modbus(device *Device, return_result bool, command []byte) (result SmartbusResult) {

	request := &SmartbusRequest{}

	if err := n.RpcCall("Modbus.Send", request, &result); err != nil {
		result.Error = err.Error()
	}

	return
}

func (n *Node) RpcCall(method string, request interface{}, reply interface{}) error {

	if n.rpcClient == nil {
		return errors.New("rpc.client is nil")
	}

	if n.netConn == nil {
		n.errors++
		return errors.New("node not connected")
	}

	if err := n.rpcClient.Call(method, request, reply); err != nil {
		n.errors++
		return err
	}

	return nil
}
