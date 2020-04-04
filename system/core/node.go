// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/mqtt"
	"sync"
	"time"
)

type Nodes []*Node

type Node struct {
	modelLock  sync.Mutex
	model      *m.Node
	mqttClient *mqtt.Client
	quit       chan struct{}
	metric     *metrics.MetricManager
	chLock     sync.Mutex
	ch         map[int64]chan *NodeResponse
	statLock   sync.Mutex
	stat       NodeStat
}

func NewNode(model *m.Node,
	mqtt *mqtt.Mqtt,
	metric *metrics.MetricManager) *Node {

	node := &Node{
		model: model,
		stat: NodeStat{
			ConnStatus: "disabled",
			LastPing:   time.Now(),
		},
		ch:         make(map[int64]chan *NodeResponse, 0),
		quit:       make(chan struct{}),
		metric:     metric,
		mqttClient: mqtt.NewClient(fmt.Sprintf("node_%v", model.Name)),
	}

	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				node.updateStatus()

			case _, ok := <-node.quit:
				if !ok {
					return
				}
				close(node.quit)
				return
			}
		}
	}()

	return node
}

func (n *Node) Remove() {

	log.Infof("Remove node %v", n.model.Id)

	n.quit <- struct{}{}
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

	n.MqttPublish(n.topic(fmt.Sprintf("req/device%d", device.Id)), msg)

	// wait response
	ticker := time.NewTimer(time.Second * 5)
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
	n.chLock.Lock()
	defer n.chLock.Unlock()

	if _, ok := n.ch[deviceId]; ok {
		return
	}

	n.ch[deviceId] = ch
}

func (n *Node) delCh(deviceId int64) {
	n.chLock.Lock()
	defer n.chLock.Unlock()

	if _, ok := n.ch[deviceId]; !ok {
		return
	}

	close(n.ch[deviceId])
	delete(n.ch, deviceId)
}

func (n *Node) Connect() *Node {

	n.mqttClient.Subscribe(n.topic("resp/#"), n.onPublish)
	n.mqttClient.Subscribe(n.topic("ping"), n.ping)

	return n
}

func (n *Node) IsConnected() bool {
	n.statLock.Lock()
	defer n.statLock.Unlock()
	return n.stat.IsConnected
}

func (n *Node) onPublish(client *mqtt.Client, msg mqtt.Message) {

	n.chLock.Lock()
	defer n.chLock.Unlock()

	resp := &NodeResponse{}
	if err := json.Unmarshal(msg.Payload, resp); err != nil {
		log.Error(err.Error())
		return
	}

	if _, ok := n.ch[resp.DeviceId]; !ok {
		return
	}

	n.ch[resp.DeviceId] <- resp
}

func (n *Node) ping(client *mqtt.Client, msg mqtt.Message) {

	var stat NodeStatModel
	_ = json.Unmarshal(msg.Payload, &stat)

	n.statLock.Lock()
	defer n.statLock.Unlock()

	//n.stat.Status = stat.Status //????
	n.stat.Thread = stat.Thread
	n.stat.Rps = stat.Rps
	n.stat.Min = stat.Min
	n.stat.Max = stat.Max
	n.stat.StartedAt = stat.StartedAt
	n.stat.LastPing = time.Now()

	return
}

func (n *Node) MqttPublish(topic string, msg interface{}) {

	data, _ := json.Marshal(msg)

	if err := n.mqttClient.Publish(topic, data); err != nil {
		log.Error(err.Error())
		return
	}
}

func (n *Node) topic(r string) string {
	return fmt.Sprintf("home/node/%s/%s", n.model.Name, r)
}

func (n *Node) GetConnStatus() string {
	n.statLock.Lock()
	defer n.statLock.Unlock()
	return n.stat.ConnStatus
}

func (n *Node) GetStat() NodeStat {
	n.statLock.Lock()
	defer n.statLock.Unlock()
	return n.stat
}

func (n *Node) UpdateClientParams(params *m.Node) {
	n.modelLock.Lock()
	defer n.modelLock.Unlock()

	n.model = params

	// unsubscribe all mqtt client
	if n.mqttClient != nil {
		n.mqttClient.UnsubscribeAll()
	}

	if n.model.Status != "disabled" {
		n.Connect()
	}

}

func (n *Node) updateStatus() {
	n.statLock.Lock()
	defer n.statLock.Unlock()

	n.stat.IsConnected = time.Now().Sub(n.stat.LastPing).Seconds() < 2

	n.modelLock.Lock()
	modelStatus := n.model.Status
	n.modelLock.Unlock()

	if modelStatus == "enabled" {
		if n.stat.IsConnected {
			n.stat.ConnStatus = "connected"
		} else {
			n.stat.ConnStatus = "wait"
		}
	} else {
		n.stat.ConnStatus = "disabled"
	}
	go n.metric.Update(metrics.NodeUpdateStatus{Id: n.model.Id, Status: n.stat.ConnStatus})
}

func (n *Node) Model() *m.Node {
	n.modelLock.Lock()
	defer n.modelLock.Unlock()
	return n.model
}
