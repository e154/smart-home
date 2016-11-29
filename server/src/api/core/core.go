package core

import (
	"log"
	"../models"
	"../stream"
	"time"
	"encoding/json"
	"reflect"
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

	err = b.InitNodes()

	if err != nil {
		return
	}

	err = b.InitWorkflows()
	return
}

// ------------------------------------------------
// Workflows
// ------------------------------------------------

// инициализация всего рабочего процесса, с запуском
// дочерни подпроцессов
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

// добавление рабочего процесс, без автоматического поиска
// и запуска подпроцессов
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

// нельзя удалить workflow если к нему привязаны связанные сущности
func (b *Core) RemoveWorkflow(workflow *models.Workflow) (err error) {

	log.Println("Remove workflow:", workflow.Name)

	if _, ok := b.workflows[workflow.Id]; !ok {
		return
	}

	b.workflows[workflow.Id].Stop()
	delete(b.workflows, workflow.Id)


	return
}

// ------------------------------------------------
// Workers
// ------------------------------------------------

func (b *Core) AddWorker(flow *models.Worker) (err error) {

	return
}

func (b *Core) UpdateWorker(flow *models.Worker) (err error) {

	return
}

func (b *Core) RemoveWorker(flow *models.Worker) (err error) {

	return
}

// ------------------------------------------------
// Flows
// ------------------------------------------------

func (b *Core) AddFlow(flow *models.Flow) (err error) {

	return
}

func (b *Core) UpdateFlow(flow *models.Flow) (err error) {

	return
}

func (b *Core) RemoveFlow(flow *models.Flow) (err error) {

	return
}

// ------------------------------------------------
// Nodes
// ------------------------------------------------

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

	return
}

func (b *Core) AddNode(node *models.Node) (err error) {

	if _, exist := b.nodes[node.Id]; exist {
		return b.ReloadNode(node)
	}

	log.Printf("Add node: \"%s\"", node.Name)

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

	log.Printf("Remove node: \"%s\"", node.Name)

	if _, exist := b.nodes[node.Id]; !exist {
		return
	}

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

	log.Printf("Reload node: \"%s\"", node.Name)

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

	log.Printf("Connect to node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes_chan[node.Id] <- "connect"
	}

	BroadcastNodesStatus()

	return
}

func (b *Core) DisconnectNode(node *models.Node) (err error) {

	log.Printf("Disconnect from node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes_chan[node.Id] <- "disconnect"
	}

	BroadcastNodesStatus()

	return
}

// ------------------------------------------------
//
// ------------------------------------------------

func streamWorkflowsStatus(client *stream.Client, value interface{}) {

	return
}

func GetNodesStatus() (result map[int64]string) {
	result = make(map[int64]string)
	for _, node := range corePtr.nodes {
		result[node.Id] = node.GetConnectStatus()
	}

	return
}

func streamNodesStatus(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	result := GetNodesStatus()
	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "nodes": result})
	client.Send(string(msg))
}

func BroadcastNodesStatus() {
	result := GetNodesStatus()
	msg, _ := json.Marshal(map[string]interface{}{"type": "broadcast", "value": &map[string]interface{}{"type": "nodes", "body": result}})
	Hub.Broadcast(string(msg))
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