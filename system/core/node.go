package core

import (
	"encoding/json"
	"errors"
	"fmt"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/surgemq/message"
	"sync"
	"time"
)

type Nodes []*Node

type Node struct {
	*m.Node
	errors      int64
	ConnStatus  string
	qClient     *mqtt.Client
	IsConnected bool
	lastPing    time.Time
	quit        bool
	stat        *NodeStatModel
	sync.Mutex
	ch map[int64]chan *NodeResponse
}

func NewNode(model *m.Node,
	mqtt *mqtt.Mqtt) *Node {

	node := &Node{
		Node:       model,
		ConnStatus: "disabled",
		ch:         make(map[int64]chan *NodeResponse, 0),
		stat:       &NodeStatModel{},
	}

	topic := fmt.Sprintf("/home/%s", model.Name)
	mqttClient, err := mqtt.NewClient(topic)
	if err != nil {
		log.Error(err.Error())
	}

	node.qClient = mqttClient

	go func() {
		for ; ; {
			if node.quit {
				break
			}
			time.Sleep(time.Second)
			node.IsConnected = time.Now().Sub(node.lastPing).Seconds() < 2 && node.stat.Thread > 0
		}
	}()

	return node
}

func (n *Node) Send(device *m.Device, command []byte) (result NodeResponse, err error) {

	//log.Debugf("send device(%v) command(%v)", device.Id, command)

	// time metric
	startTime := time.Now()

	ch := make(chan *NodeResponse)
	n.addCh(device.Id, ch)
	defer n.delCh(device.Id)

	// send message to node
	msg := &NodeMessage{
		DeviceId:   device.Id,
		DeviceType: device.Type,
		Properties: device.Properties,
		Command:    command,
	}

	data, _ := json.Marshal(msg)
	if err = n.qClient.Publish(data); err != nil {
		log.Error(err.Error())
		return
	}

	// wait response
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	var done bool
	for ; ; {
		if done {
			break
		}
		select {
		case <-ticker.C:
			//log.Debugf("request timeout device(%d)", device.Id)
			err = errors.New("timeout")
			done = true
		case resp := <-ch:
			if resp == nil {
				return
			}
			if resp.DeviceId != device.Id {
				continue
			}

			// response from node
			result = NodeResponse{
				DeviceId:   resp.DeviceId,
				Status:     resp.Status,
				DeviceType: resp.DeviceType,
				Response:   resp.Response,
				Properties: resp.Properties,
			}
			done = true
		}
	}

	result.Time = time.Since(startTime).Seconds()

	return
}

func (n *Node) addCh(deviceId int64, ch chan *NodeResponse) {
	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[deviceId]; ok {
		return
	}

	n.ch[deviceId] = ch
}

func (n *Node) delCh(deviceId int64) {
	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[deviceId]; !ok {
		return
	}

	close(n.ch[deviceId])
	delete(n.ch, deviceId)
}

//func (n *NodeModel) onComplete(msg, ack message.Message, err error) error {
//	log.Debug("onComplete")
//	return nil
//}

func (n *Node) onPublish(msg *message.PublishMessage) (err error) {

	resp := &NodeResponse{}
	if err = json.Unmarshal(msg.Payload(), resp); err != nil {
		return
	}

	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[resp.DeviceId]; !ok {
		return
	}

	n.ch[resp.DeviceId] <- resp

	return
}

func (n *Node) Connect() *Node {

	if err := n.qClient.Connect(); err != nil {
		log.Error(err.Error())
	}

	time.Sleep(time.Second)

	topic := fmt.Sprintf("/home/%s", n.Name)
	if err := n.qClient.Subscribe(topic+"/resp", nil, n.onPublish); err != nil {
		log.Warning(err.Error())
	}

	if err := n.qClient.Subscribe(topic+"/ping", nil, n.ping); err != nil {
		log.Warning(err.Error())
	}

	return n
}

func (n *Node) Disconnect() {
	n.quit = true
	if n.qClient != nil {
		n.qClient.Disconnect()
	}
}

func (n *Node) ping(msg *message.PublishMessage) (err error) {

	_ = json.Unmarshal(msg.Payload(), n.stat)
	n.lastPing = time.Now()

	switch n.stat.Status {
	case "busy":
		n.ConnStatus = "wait"
	case "enabled":
		n.ConnStatus = "connected"
	default:
		n.ConnStatus = "enabled"
	}
	return
}
