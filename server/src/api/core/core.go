package core

import (
	"../log"
	"time"
	"../models"
	"../stream"
	cr "github.com/e154/cron"
	"sync"
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
	nodes      	map[int64]*models.Node
	nodes_chan 	map[int64]chan string
	workflows  	map[int64]*Workflow
	mu		sync.Mutex
	deviceStates	map[int64]*models.DeviceState
	telemetry	Telemetry
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
				for _, action := range worker.actions {
					if action.Device.Id == device.Id {
						workflow.UpdateFlow(flow.Model)
						continue
					}

					if action.Device != nil && action.Device.Id == device.Id {
						workflow.UpdateFlow(flow.Model)
						continue
					}
				}

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

func (b *Core) DoWorker(model *models.Worker) (err error) {

	for _, workflow := range b.workflows {
		for _, flow := range workflow.Flows {
			if worker, ok := flow.Workers[model.Id]; ok {
				worker.Do()
				break
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

	b.telemetry.Broadcast("nodes")

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

	delete(b.nodes, node.Id)

	b.telemetry.Broadcast("nodes")

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

	b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) ConnectNode(node *models.Node) (err error) {

	log.Infof("Connect to node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes_chan[node.Id] <- "connect"
	}

	b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) DisconnectNode(node *models.Node) (err error) {

	log.Infof("Disconnect from node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes_chan[node.Id] <- "disconnect"
	}

	b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) GetNodes() (map[int64]*models.Node) {
	return b.nodes
}

// ------------------------------------------------
// Device states
// ------------------------------------------------
func (b *Core) SetDeviceState(id int64, state *models.DeviceState) {
	b.mu.Lock()

	if old_state, ok := b.deviceStates[id]; ok {
		if old_state.Id == state.Id {
			b.mu.Unlock()
			return
		}
	}

	b.deviceStates[id] = state
	b.mu.Unlock()

	b.telemetry.BroadcastOne("devices", id)
}

func (b *Core) GetDevicesStates() map[int64]*models.DeviceState {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.deviceStates
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

func Initialize(telemetry Telemetry) (err error) {
	log.Info("Core initialize...")

	if cron == nil {
		cron = cr.NewCron()
		cron.Run()
	}

	corePtr = &Core{
		nodes: make(map[int64]*models.Node),
		nodes_chan: make(map[int64]chan string),
		workflows: make(map[int64]*Workflow),
		deviceStates: make(map[int64]*models.DeviceState),
		telemetry: telemetry,
	}
	if err = corePtr.Run(); err != nil {
		return
	}

	Hub = stream.GetHub()
	Hub.Subscribe("do.worker", streamDoWorker)
	Hub.Subscribe("do.action", streamDoAction)

	return
}