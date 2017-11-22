package dasboard

import (
	"reflect"
	"encoding/json"
	"github.com/e154/smart-home/api/core"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/stream"
	"time"
)

func NewNode() (node *Nodes) {

	node = &Nodes{
		Status: make(map[int64]string),
	}

	return
}

type Nodes struct {
	Total	int64			`json:"total"`
	Status	map[int64]string	`json:"status"`
}

func (n *Nodes) Update() {
	n.Total, _ = models.GetNodesCount()
	nodes := core.CorePtr().GetNodes()

	n.Status = make(map[int64]string)
	for _, node := range nodes {
		n.Status[node.Id] = node.GetConnectStatus()
	}
}

func (n *Nodes) Broadcast() {

	time.Sleep(time.Second)

	n.Update()

	msg, _ := json.Marshal(map[string]interface{}{"type": "broadcast",
		"value": map[string]interface{}{"type": "telemetry", "body": map[string]interface{}{
			"nodes": n,
		}}},
	)
	Hub.Broadcast(string(msg))
}

// only on request: 'dashboard.get.nodes.status'
//
func (n *Nodes) streamNodesStatus(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	n.Update()

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "nodes": n})
	client.Send(string(msg))
}