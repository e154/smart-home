package core

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/surgemq/message"
)

type Nodes []*Node

type Node struct {
	*m.Node
	errors     int64
	connStatus string
	mqttClient *mqtt.Client
}

func NewNode(model *m.Node,
	mqtt *mqtt.Mqtt) *Node {

	node := &Node{
		Node: model,
	}

	mqttClient, err := mqtt.NewClient(model.Name, node.onComplete, node.onPublish)
	if err != nil {
		log.Error(err.Error())
	}

	node.mqttClient = mqttClient

	return node
}

func (n *Node) Send(device *m.Device, command []byte) (err error) {
	log.Debugf("send device(%v) command(%v)", device, command)
	//n.mqttClient.Publish(command)
	return
}

func (n *Node) onComplete(msg, ack message.Message, err error) error {
	log.Debug("onComplete")
	return nil
}

func (n *Node) onPublish(msg *message.PublishMessage) error {
	log.Debug("onPublish")
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

}
