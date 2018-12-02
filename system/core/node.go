package core

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/surgemq/message"
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common"
	"time"
	"sync"
)

type Nodes []*Node

type NodeMessage struct {
	DeviceId   int64             `json:"device_id"`
	DeviceType common.DeviceType `json:"device_type"`
	Properties json.RawMessage   `json:"properties"`
	Command    []byte            `json:"command"`
}

type NodeResponse struct {
	DeviceId   int64             `json:"device_id"`
	DeviceType common.DeviceType `json:"device_type"`
	Properties json.RawMessage   `json:"properties"`
	Response   []byte            `json:"response"`
	Status     string            `json:"status"`
}

type Node struct {
	*m.Node
	errors      int64
	connStatus  string
	qClient     *mqtt.Client
	IsConnected bool
	lastPing    time.Time
	sync.Mutex
	ch map[int64]chan *NodeResponse
}

func NewNode(model *m.Node,
	mqtt *mqtt.Mqtt) *Node {

	node := &Node{
		Node: model,
		ch:   make(map[int64]chan *NodeResponse, 0),
	}

	topic := fmt.Sprintf("/home/%s", model.Name)
	mqttClient, err := mqtt.NewClient(topic, nil, node.onPublish)
	if err != nil {
		log.Error(err.Error())
	}

	node.qClient = mqttClient

	go node.pong()

	return node
}

func (n *Node) Send(device *m.Device, command []byte) (err error) {

	log.Debugf("send device(%v) command(%v)", device.Id, command)

	ch := make(chan *NodeResponse)
	n.addCh(device.Id, ch)
	defer n.delCh(device.Id)

	msg := &NodeMessage{
		DeviceId:   device.Id,
		DeviceType: device.Type,
		Properties: device.Properties,
		Command:    command,
	}

	data, _ := json.Marshal(msg)
	if err = n.qClient.Publish(data); err != nil {
		log.Error(err.Error())
	}

	timeout := time.Second * 2
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	var done bool
	for ; ; {
		if done {
			break
		}
		select {
		case <-ticker.C:
			log.Debugf("request timeout device(%d)", device.Id)
			done = true
		case resp := <-ch:
			if resp.DeviceId != device.Id {
				continue
			}
			fmt.Println("//////////////")
			fmt.Println(resp)
		}
	}

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

//func (n *Node) onComplete(msg, ack message.Message, err error) error {
//	log.Debug("onComplete")
//	return nil
//}

func (n *Node) onPublish(msg *message.PublishMessage) (err error) {
	if string(msg.Payload()) == "live" {
		n.lastPing = time.Now()
		return
	}

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
	log.Debug("Connect")
	if err := n.qClient.Connect(); err != nil {
		log.Error(err.Error())
	}
	return n
}

func (n *Node) Disconnect() {
	if n.qClient != nil {
		n.qClient.Disconnect()
	}
}

func (n *Node) pong() {
	go func() {
		for ; ; {
			time.Sleep(time.Second)
			n.IsConnected = time.Now().Sub(n.lastPing).Seconds() < 10
		}
	}()
}
