package core

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/surgemq/message"
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common"
	"time"
)

type Nodes []*Node

type NodeMessage struct {
	DeviceId   int64             `json:"device_id"`
	DeviceType common.DeviceType `json:"device_type"`
	Properties json.RawMessage   `json:"properties"`
	Command    []byte            `json:"command"`
}

type Node struct {
	*m.Node
	errors      int64
	connStatus  string
	mqttClient  *mqtt.Client
	IsConnected bool
	lastPing time.Time
}

func NewNode(model *m.Node,
	mqtt *mqtt.Mqtt) *Node {

	node := &Node{
		Node: model,
	}

	topic := fmt.Sprintf("/home/%s", model.Name)
	mqttClient, err := mqtt.NewClient(topic, nil, node.onPublish)
	if err != nil {
		log.Error(err.Error())
	}

	node.mqttClient = mqttClient

	go node.pong()

	return node
}

func (n *Node) Send(device *m.Device, command []byte) (err error) {

	log.Debugf("send device(%v) command(%v)", device.Id, command)

	msg := &NodeMessage{
		DeviceId:   device.Id,
		DeviceType: device.Type,
		Properties: device.Properties,
		Command:    command,
	}

	data, _ := json.Marshal(msg)
	err = n.mqttClient.Publish(data)
	return
}

//func (n *Node) onComplete(msg, ack message.Message, err error) error {
//	log.Debug("onComplete")
//	return nil
//}

func (n *Node) onPublish(msg *message.PublishMessage) error {
	if string(msg.Payload()) == "live" {
		n.lastPing = time.Now()
		return nil
	}
	return nil
}

func (n *Node) Connect() *Node {
	log.Debug("Connect")
	if err := n.mqttClient.Connect(); err != nil {
		log.Error(err.Error())
	}
	return n
}

func (n *Node) Disconnect() {
	if n.mqttClient != nil {
		n.mqttClient.Disconnect()
	}
}

func (n *Node) pong() {
	go func() {
		for ;; {
			time.Sleep(time.Second)
			n.IsConnected = time.Now().Sub(n.lastPing).Seconds() < 10
		}
	}()
}
