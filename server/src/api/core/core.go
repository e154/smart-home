package core

import (
	"../log"
	"time"
	"encoding/json"
	"reflect"
	"../models"
	"../stream"
	cr "github.com/e154/cron"
)

var (
	Hub		stream.Hub
	corePtr         *Core = nil
	cron		*cr.Cron = nil
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

	log.Info("Add workflow:", workflow.Name)

	if _, ok := b.workflows[workflow.Id]; ok {
		return
	}

	wf := NewWorkflow(workflow, b.nodes)

	if err = wf.Run(); err != nil {
		return
	}

	b.workflows[workflow.Id] = wf

	return
}

// нельзя удалить workflow, если присутствуют связанные сущности
func (b *Core) RemoveWorkflow(workflow *models.Workflow) (err error) {

	log.Info("Remove workflow:", workflow.Name)

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

func (b *Core) UpdateFlowFromDevice(device *models.Device) (err error) {

	for _, workflow := range b.workflows {
		for _, flow := range workflow.Flows {
			for _, worker := range flow.Workers {
				//TODO fix
				//for _, dev := range worker.Devices {
				//	if dev.Device.Id == device.Id {
				//		workflow.UpdateFlow(flow.Model)
				//		continue
				//	}
				//
				//	if dev.Device != nil && dev.Device.Id == device.Id {
				//		workflow.UpdateFlow(flow.Model)
				//		continue
				//	}
				//}

				if device.Device != nil && worker.Model.DeviceAction.Device.Id == device.Device.Id {
					workflow.UpdateFlow(flow.Model)
					continue
				}
			}
		}
	}

	return
}

func (b *Core) UpdateWorker(_worker *models.Worker) (err error) {

	for _, workflow := range b.workflows {
		for _, flow := range workflow.Flows {
			for _, worker := range flow.Workers {
				if worker.Model.Id == _worker.Id {
					flow.UpdateWorker(_worker)
					break
				}
			}
		}
	}

	return
}

func (b *Core) RemoveWorker(_worker *models.Worker) (err error) {

	for _, workflow := range b.workflows {
		for _, flow := range workflow.Flows {
			for _, worker := range flow.Workers {
				if worker.Model.Id == _worker.Id {
					flow.RemoveWorker(_worker)
					break
				}
			}
		}
	}

	return
}

// ------------------------------------------------
// Flows
// ------------------------------------------------

func (b *Core) AddFlow(flow *models.Flow) (err error) {

	if _, ok := b.workflows[flow.Workflow.Id]; !ok {
		return
	}

	if err = b.workflows[flow.Workflow.Id].AddFlow(flow); err != nil {
		return
	}

	return
}

func (b *Core) UpdateFlow(flow *models.Flow) (err error) {

	if _, ok := b.workflows[flow.Workflow.Id]; !ok {
		return
	}

	if err = b.workflows[flow.Workflow.Id].UpdateFlow(flow); err != nil {
		return
	}

	return
}

func (b *Core) RemoveFlow(flow *models.Flow) (err error) {

	if _, ok := b.workflows[flow.Workflow.Id]; !ok {
		return
	}

	if err = b.workflows[flow.Workflow.Id].RemoveFlow(flow); err != nil {
		return
	}

	return
}

// ------------------------------------------------
// Nodes
// ------------------------------------------------

func (b *Core) InitNodes() (err error) {

	var nodes []*models.Node
	b.nodes = make(map[int64]*models.Node)
	b.nodes_chan = make(map[int64]chan string)

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

	log.Infof("Add node: \"%s\"", node.Name)

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
					log.Infof("Node dial tcp %s:%d ... ok",node.Ip, node.Port)
					connect = false
					node.SetConnectStatus("connected")
				} else {
					node.Errors++
					if node.Errors == 7 {
						log.Errorf("Node error %s", err.Error())
					}
					node.SetConnectStatus("error")
				}
			}

			time.Sleep(time.Second)
		}
	}(b.nodes_chan[node.Id])

	return
}

func (b *Core) RemoveNode(node *models.Node) (err error) {

	log.Infof("Remove node: \"%s\"", node.Name)

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

	log.Infof("Reload node: \"%s\"", node.Name)

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

	log.Infof("Connect to node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes_chan[node.Id] <- "connect"
	}

	BroadcastNodesStatus()

	return
}

func (b *Core) DisconnectNode(node *models.Node) (err error) {

	log.Infof("Disconnect from node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes_chan[node.Id] <- "disconnect"
	}

	BroadcastNodesStatus()

	return
}

// ------------------------------------------------
// Script
// ------------------------------------------------
func (b *Core) UpdateScript(script *models.Script) (err error) {

	for _, workflow := range b.workflows {
		for _, flow := range workflow.Flows {

			for _, worker := range flow.Workers {
				if worker.Model.DeviceAction.Script.Id == script.Id {
					workflow.UpdateFlow(flow.Model)
				}
			}

			for _, flowElement := range flow.FlowElements {
				if flowElement.Model.Script == nil {
					continue
				}

				if flowElement.Model.Script.Id == script.Id {
					workflow.UpdateFlow(flow.Model)
				}
			}
		}
	}

	return
}

// ------------------------------------------------
// etc
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
	log.Info("Core initialize...")

	if cron == nil {
		cron = cr.NewCron()
		cron.Run()
	}

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