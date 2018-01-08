package dasboard

import (
	"reflect"
	"encoding/json"
	"github.com/e154/smart-home/api/core"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/stream"
	"sync"
	"time"
)

func NewNode() (node *Nodes) {

	node = &Nodes{
		Status: make(map[int64]string),
	}

	return
}

type Nodes struct {
	sync.Mutex
	Total      int64            `json:"total"`
	Status     map[int64]string `json:"status"`
	lastUpdate time.Time
}

func (n *Nodes) Update()  {

	n.Lock()
	defer n.Unlock()

	if time.Now().Sub(n.lastUpdate).Seconds() < 15 {
		return
	}

	n.lastUpdate = time.Now()

	n.Total, _ = models.GetNodesCount()
	nodes := core.CorePtr().GetNodes()

	n.Status = make(map[int64]string)

	for _, node := range nodes {
		n.Status[node.Id] = node.GetConnectStatus()
	}
}

func (n *Nodes) Broadcast() (interface{}, bool) {

	n.Update()

	return map[string]interface{}{
		"nodes": n,
	}, true
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
	client.Send <- msg
}