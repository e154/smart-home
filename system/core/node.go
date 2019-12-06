package core

import (
	"encoding/json"
	"errors"
	"fmt"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/mqtt_client"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"sync"
	"time"
)

type Nodes []*Node

type Node struct {
	*m.Node
	errors     int64
	ConnStatus string
	mqttClient *mqtt_client.Client2
	lastPing   time.Time
	stat       *NodeStatModel
	sync.Mutex
	ch map[int64]chan *NodeResponse
}

func NewNode(model *m.Node, mqtt *mqtt.Mqtt) *Node {

	mqttClient, err := mqtt.NewClient(nil)
	if err != nil {
		log.Error(err.Error())
	}

	node := &Node{
		Node:       model,
		ConnStatus: "disabled",
		ch:         make(map[int64]chan *NodeResponse, 0),
		stat:       &NodeStatModel{},
		mqttClient: mqttClient,
	}

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

	n.MqttPublish(msg)

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

func (n *Node) Connect() *Node {

	if err := n.mqttClient.Connect(); err != nil {
		log.Error(err.Error())
	}

	time.Sleep(time.Second)

	// /home/node/resp
	sub1 := mqtt_client.Subscribe{
		Qos:      0,
		Callback: n.onPublish,
	}
	topic := fmt.Sprintf("/home/%s/resp", n.Name)
	if err := n.mqttClient.Subscribe(topic, sub1); err != nil {
		log.Warning(err.Error())
	}

	// /home/node/ping
	sub2 := mqtt_client.Subscribe{
		Qos:      0,
		Callback: n.ping,
	}
	topic = fmt.Sprintf("/home/%s/ping", n.Name)
	if err := n.mqttClient.Subscribe(topic, sub2); err != nil {
		log.Warning(err.Error())
	}

	return n
}

func (n *Node) Disconnect() {
	if n.mqttClient != nil {
		n.mqttClient.Disconnect()
	}
}

func (n *Node) IsConnected() bool {
	return time.Now().Sub(n.lastPing).Seconds() < 2
}

func (n *Node) onPublish(client MQTT.Client, msg MQTT.Message) {

	resp := &NodeResponse{}
	if err := json.Unmarshal(msg.Payload(), resp); err != nil {
		log.Error(err.Error())
		return
	}

	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[resp.DeviceId]; !ok {
		return
	}

	n.ch[resp.DeviceId] <- resp
}

func (n *Node) ping(client MQTT.Client, msg MQTT.Message) {

	_ = json.Unmarshal(msg.Payload(), n.stat)
	n.lastPing = time.Now()

	switch n.stat.Status {
	case "enabled":
		n.ConnStatus = "connected"
	default:
		n.ConnStatus = "enabled"
	}
	return
}

func (n *Node) MqttPublish(msg interface{}) {

	data, _ := json.Marshal(msg)
	topic := fmt.Sprintf("/home/%s/req", n.Node.Name)
	if err := n.mqttClient.Publish(topic, data); err != nil {
		log.Error(err.Error())
		return
	}
}