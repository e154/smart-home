package dashboard

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
	"reflect"
	"sync"
	"time"
)

type Nodes struct {
	sync.Mutex
	Total      int64              `json:"total"`
	Status     map[int64]string   `json:"status"`
	lastUpdate time.Time          `json:"-"`
	adaptors   *adaptors.Adaptors `json:"-"`
	core       *core.Core         `json:"-"`
}

func NewNode(adaptors *adaptors.Adaptors) (node *Nodes) {

	node = &Nodes{
		Status:   make(map[int64]string),
		adaptors: adaptors,
	}

	return
}

func (n *Nodes) Update() {

	if n.core == nil {
		return
	}

	n.Lock()
	defer n.Unlock()

	if time.Now().Sub(n.lastUpdate).Seconds() < 15 {
		return
	}

	n.lastUpdate = time.Now()

	_, n.Total, _ = n.adaptors.Node.List(999, 0, "", "")
	nodes := n.core.GetNodes()

	n.Status = make(map[int64]string)

	for _, node := range nodes {
		n.Status[node.Id] = node.Status
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
