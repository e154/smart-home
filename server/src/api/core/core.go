package core

import (
	"log"
	"fmt"
	"time"
	"../models"
	"../stream"
	"github.com/e154/cron"
)

var (
	Hub		stream.Hub
	corePtr         *Core = nil
)

func CorePtr() *Core {
	return corePtr
}

type Core struct {
	nodes      map[int64]*models.Node
	nodes_chan map[int64]chan string
	workflows  map[int64]*Workflow
}

func (b *Core) Run() (err error) {

	if err = b.InitNodes(); err != nil {
		return
	}

	err = b.InitWorkflows()

	return
}

func (b *Core) InitNodes() (err error) {
	var nodes []*models.Node
	b.nodes = make(map[int64]*models.Node)
	b.nodes_chan = make(map[int64]chan string)

	//log.Println("--------------------- NODES ---------------------")
	if nodes, err = models.GetAllEnabledNodes(); err != nil {
		return
	}

	for _, node := range nodes {
		b.AddNode(node)
	}

	//TODO remove
	if len(b.nodes) == 0 {
		return
	}

	return
}

func (b *Core) InitWorkflows() (err error) {

	b.workflows = make(map[int64]*Workflow)
	//log.Println("------------------- WORKFLOW --------------------")
	workflows, err := models.GetAllEnabledWorkflow()
	if err != nil {
		return
	}

	for _, workflow := range workflows {
		b.AddWorkflow(workflow)
	}

	return
}

func (b *Core) Stop() (err error) {

	for _, wf := range b.workflows {
		if err = wf.Stop(); err != nil {
			return
		}
	}
	return
}

func (b *Core) Restart() (err error) {

	for _, wf := range b.workflows {
		if err = wf.Restart(); err != nil {
			return
		}
	}
	return
}

func (b *Core) AddFlow(f *models.Flow) (err error) {

	var flow *models.Flow
	if flow, err = models.GetFlowById(f.Id); err != nil {
		return
	}

	if _, ok := b.workflows[flow.Workflow.Id]; ok {
		if err = b.workflows[flow.Workflow.Id].AddFlow(flow); err != nil {
			return
		}
	}

	return
}

func (b *Core) UpdateFlow(f *models.Flow) (err error) {

	var flow *models.Flow
	if flow, err = models.GetFlowById(f.Id); err != nil {
		return
	}

	if _, ok := b.workflows[flow.Workflow.Id]; ok {
		if err = b.workflows[flow.Workflow.Id].UpdateFlow(flow); err != nil {
			return
		}
	}

	return
}

func (b *Core) RemoveFlow(f *models.Flow) (err error) {

	var flow *models.Flow
	if flow, err = models.GetFlowById(f.Id); err != nil {
		return
	}

	if _, ok := b.workflows[flow.Workflow.Id]; ok {
		if err = b.workflows[flow.Workflow.Id].RemoveFlow(flow); err != nil {
			return
		}
	}

	return
}

func (b *Core) AddWorkflow(workflow *models.Workflow) (err error) {

	log.Println("Add workflow:", workflow.Name)

	if _, ok := b.workflows[workflow.Id]; ok {
		return
	}

	wf := &Workflow{
		model: workflow,
		Nodes: b.nodes,
		CronTasks: make(map[int64]*cron.Task),
	}
	if err = wf.Run(); err != nil {
		return
	}

	b.workflows[workflow.Id] = wf

	return
}

func (b *Core) RemoveWorkflow(workflow *models.Workflow) (err error) {

	log.Println("Remove workflow:", workflow.Name)

	if _, ok := b.workflows[workflow.Id]; !ok {
		return
	}

	b.workflows[workflow.Id].Stop()
	delete(b.workflows, workflow.Id)

	return
}

func (b *Core) AddNode(node *models.Node) (err error) {
	if _, ok := b.nodes[node.Id]; ok {
		return
	}

	b.nodes[node.Id] = node
	b.nodes_chan[node.Id] = make(chan string)

	go func(ch <- chan string) {
		var quit, disconnect bool
		connect := true
		for ;; {

			select {
			case c := <- ch:
				switch c {
				case "quit":
					quit = true
				case "disconnect":
					disconnect = true
				case "connect":
					connect = true
				default:

				}

			default:

			}

			if quit {
				node.TcpClose()
				break
			}

			if node.Errors > 5 {
				connect = true
			}

			if disconnect {
				node.TcpClose()
				connect = false
				disconnect = false
			}

			if connect {
				disconnect = false
				node.TcpClose()

				if _, err := node.RpcDial(); err == nil {
					node.Errors = 0
					//log.Printf("Node dial tcp %s:%d ... ok",node.Ip, node.Port)
					connect = false
					node.SetConnectStatus("connected")
				} else {
					//log.Printf("Node error %s", err.Error())
					node.SetConnectStatus("error")
				}
			}

			time.Sleep(time.Second)
		}
	}(b.nodes_chan[node.Id])

	return
}

func (b *Core) RemoveNode(node *models.Node) (err error) {

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes_chan[node.Id] <- "quit"
		close(b.nodes_chan[node.Id])
		delete(b.nodes_chan, node.Id)
		delete(b.nodes, node.Id)
	}

	BroadcastNodesStatus()

	return
}

func (b *Core) ReloadNode(node *models.Node) (err error) {

	if _, ok := b.nodes[node.Id]; !ok {
		b.AddNode(node)
		return
	}

	b.nodes[node.Id].Status = node.Status
	b.nodes[node.Id].Ip = node.Ip
	b.nodes[node.Id].Port = node.Port
	b.nodes[node.Id].SetConnectStatus("wait")
	if node.Status == "disabled" {
		b.nodes_chan[node.Id] <- "disconnect"
	} else {
		b.nodes_chan[node.Id] <- "connect"
	}

	BroadcastNodesStatus()

	return
}

func (b *Core) ConnectNode(node *models.Node) (err error) {

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes_chan[node.Id] <- "connect"
	}

	BroadcastNodesStatus()

	return
}

func (b *Core) DisconnectNode(node *models.Node) (err error) {

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes_chan[node.Id] <- "disconnect"
	}

	BroadcastNodesStatus()

	return
}

func (b *Core) AddWorker(worker *models.Worker) (err error) {

	if _, ok := b.workflows[worker.Workflow.Id]; ok {
		err = b.workflows[worker.Workflow.Id].AddWorker(worker)
	} else {
		err = fmt.Errorf("workflow id:%d not found", worker.Workflow.Id)
	}

	return
}

func (b *Core) UpdateWorker(worker *models.Worker) (err error) {

	if _, ok := b.workflows[worker.Workflow.Id]; ok {
		err = b.workflows[worker.Workflow.Id].UpdateWorker(worker)
	} else {
		err = fmt.Errorf("workflow id:%d not found", worker.Workflow.Id)
	}

	return
}

func (b *Core) RemoveWorker(worker *models.Worker) (err error) {

	if _, ok := b.workflows[worker.Workflow.Id]; ok {
		err = b.workflows[worker.Workflow.Id].RemoveWorker(worker)
	} else {
		err = fmt.Errorf("workflow id:%d not found", worker.Workflow.Id)
	}

	return
}

func Initialize() (err error) {
	log.Println("Core initialize...")

	corePtr = &Core{}
	if err = corePtr.Run(); err != nil {
		return
	}

	Hub = stream.GetHub()
	Hub.Subscribe("get.nodes.status", streamNodesStatus)
	Hub.Subscribe("get.workflow.status", streamWorkflowsStatus)
	Hub.Subscribe("get.flows.status", streamWorkflowsStatus)

	return
}