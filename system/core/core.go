package core

import (
	"sync"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("core")
)

type Core struct {
	nodes map[int64]*m.Node
	//workflows  		map[int64]*Workflow
	//telemetry		Telemetry
	//Map				*Map
	sync.Mutex
	adaptors *adaptors.Adaptors
}

func NewCore(adaptors *adaptors.Adaptors) (core *Core, err error) {
	core = &Core{
		nodes:    make(map[int64]*m.Node),
		adaptors: adaptors,
	}

	return
}

func (c *Core) Start() (err error) {
	err = c.initNodes()
	return
}

// ------------------------------------------------
// Nodes
// ------------------------------------------------

func (c *Core) initNodes() (err error) {

	var nodes []*m.Node
	if nodes, err = c.adaptors.Node.GetAllEnabled(); err != nil {
		return
	}

	for _, model_node := range nodes {
		c.AddNode(model_node)
	}

	return
}

func (b *Core) AddNode(node *m.Node) (err error) {

	if _, exist := b.nodes[node.Id]; exist {
		return b.ReloadNode(node)
	}

	log.Infof("Add node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		return
	}

	b.Lock()
	node.Start()
	b.nodes[node.Id] = node.Connect()
	b.Unlock()

	//TODO add telemetry
	//b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) RemoveNode(node *m.Node) (err error) {

	log.Infof("Remove node: \"%s\"", node.Name)

	if _, exist := b.nodes[node.Id]; !exist {
		return
	}

	b.Lock()
	if _, ok := b.nodes[node.Id]; ok {
		b.nodes[node.Id].Disconnect()
		delete(b.nodes, node.Id)
	}

	delete(b.nodes, node.Id)
	b.Unlock()

	//TODO add telemetry
	//b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) ReloadNode(node *m.Node) (err error) {

	log.Infof("Reload node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; !ok {
		b.AddNode(node)
		return
	}

	b.Lock()
	b.nodes[node.Id].Status = node.Status
	b.nodes[node.Id].Ip = node.Ip
	b.nodes[node.Id].Port = node.Port
	b.nodes[node.Id].SetConnectStatus("wait")
	b.Unlock()

	if b.nodes[node.Id].Status == "disabled" {
		b.nodes[node.Id].Disconnect()
	} else {
		b.nodes[node.Id].Connect()
	}

	return
}

func (b *Core) ConnectNode(node *m.Node) (err error) {

	log.Infof("Connect to node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes[node.Id].Connect()
	}

	//TODO add telemetry
	//b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) DisconnectNode(node *m.Node) (err error) {

	log.Infof("Disconnect from node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes[node.Id].Disconnect()
	}

	//TODO add telemetry
	//b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) GetNodes() (nodes map[int64]*m.Node) {

	nodes = make(map[int64]*m.Node)

	b.Lock()
	for id, node := range b.nodes {
		nodes[id] = node
	}
	b.Unlock()

	return
}