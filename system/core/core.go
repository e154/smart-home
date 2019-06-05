package core

import (
	"errors"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	cr "github.com/e154/smart-home/system/cron"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
	"github.com/op/go-logging"
	"sync"
)

var (
	log = logging.MustGetLogger("core")
)

type IGetMap interface {
	GetMap() *Map
}

type ITelemetry interface {
	Broadcast(pack string)
	BroadcastOne(pack string, deviceId int64, elementName string)
	Run(core *Core)
	Stop()
}

type Core struct {
	sync.Mutex
	nodes         map[int64]*Node
	workflows     map[int64]*Workflow
	adaptors      *adaptors.Adaptors
	scripts       *scripts.ScriptService
	cron          *cr.Cron
	mqtt          *mqtt.Mqtt
	telemetry     ITelemetry
	streamService *stream.StreamService
	Map           *Map
}

func NewCore(adaptors *adaptors.Adaptors,
	scripts *scripts.ScriptService,
	graceful *graceful_service.GracefulService,
	cron *cr.Cron,
	mqtt *mqtt.Mqtt,
	telemetry ITelemetry,
	streamService *stream.StreamService) (core *Core, err error) {

	core = &Core{
		nodes:         make(map[int64]*Node),
		workflows:     make(map[int64]*Workflow),
		adaptors:      adaptors,
		scripts:       scripts,
		cron:          cron,
		mqtt:          mqtt,
		telemetry:     telemetry,
		streamService: streamService,
		Map: &Map{
			telemetry: telemetry,
		},
	}

	graceful.Subscribe(core)

	scripts.PushStruct("Map", &MapBind{Map: core.Map})

	return
}

func (c *Core) Run() (err error) {
	if err = c.initNodes(); err != nil {
		return
	}

	if err = c.InitWorkflows(); err != nil {
		return
	}

	c.telemetry.Run(c)

	// register steam actions
	c.streamService.Subscribe("do.worker", streamDoWorker(c))
	c.streamService.Subscribe("do.action", streamDoAction(c))

	return
}

func (b *Core) Stop() (err error) {

	for _, workflow := range b.workflows {
		if err = b.DeleteWorkflow(workflow.model); err != nil {
			return
		}
	}

	for _, node := range b.nodes {
		if err = b.RemoveNode(&m.Node{Id: node.Id, Name: node.Name}); err != nil {
			return
		}
	}

	b.telemetry.Stop()

	// unregister steam actions
	b.streamService.UnSubscribe("do.worker")
	b.streamService.UnSubscribe("do.action")

	return
}

func (b *Core) Restart() (err error) {

	if err = b.Stop(); err != nil {
		log.Error(err.Error())
	}

	if err = b.Run(); err != nil {
		log.Error(err.Error())
	}

	return
}

func (b *Core) Shutdown() {
	b.Stop()
}

// ------------------------------------------------
// Nodes
// ------------------------------------------------

func (c *Core) initNodes() (err error) {

	var nodes []*m.Node
	if nodes, err = c.adaptors.Node.GetAllEnabled(); err != nil {
		return
	}

	for _, modelNode := range nodes {
		c.AddNode(modelNode)
	}

	return
}

func (b *Core) AddNode(node *m.Node) (n *Node, err error) {

	if _, exist := b.nodes[node.Id]; exist {
		err = b.ReloadNode(node)
		return
	}

	log.Infof("Add node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		return
	}

	b.Lock()
	n = NewNode(node, b.mqtt)
	b.nodes[node.Id] = n.Connect()
	b.Unlock()

	b.telemetry.Broadcast("nodes")

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

	b.telemetry.Broadcast("nodes")

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
	//b.nodes[node.Id].SetConnectStatus("wait")
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

	b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) DisconnectNode(node *m.Node) (err error) {

	log.Infof("Disconnect from node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes[node.Id].Disconnect()
	}

	b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) GetNodes() (nodes map[int64]*Node) {

	nodes = make(map[int64]*Node)

	b.Lock()
	for id, node := range b.nodes {
		nodes[id] = node
	}
	b.Unlock()

	return
}

func (b *Core) GetNodeById(nodeId int64) *Node {

	b.Lock()
	for id, node := range b.nodes {
		if id == nodeId {
			b.Unlock()
			return node
		}
	}
	b.Unlock()

	return nil
}

// ------------------------------------------------
// Workflows
// ------------------------------------------------

// инициализация всего рабочего процесса, с запуском
// дочерни подпроцессов
func (b *Core) InitWorkflows() (err error) {

	workflows, err := b.adaptors.Workflow.GetAllEnabled()
	if err != nil {
		return
	}

	for _, workflow := range workflows {
		if err = b.AddWorkflow(workflow); err != nil {
			return
		}
	}

	return
}

// добавление рабочего процесс
func (b *Core) AddWorkflow(workflow *m.Workflow) (err error) {

	log.Infof("Add workflow: %s", workflow.Name)

	if _, ok := b.workflows[workflow.Id]; ok {
		return
	}

	wf := NewWorkflow(workflow, b.adaptors, b.scripts, b.cron, b)

	if err = wf.Run(); err != nil {
		return
	}

	b.workflows[workflow.Id] = wf

	return
}

func (wf *Core) GetWorkflow(workflowId int64) (workflow *Workflow, err error) {

	if _, ok := wf.workflows[workflowId]; !ok {
		err = errors.New("not found")
		return
	}

	workflow = wf.workflows[workflowId]

	return
}

// нельзя удалить workflow, если присутствуют связанные сущности
func (b *Core) DeleteWorkflow(workflow *m.Workflow) (err error) {

	log.Infof("Remove workflow: %s", workflow.Name)

	if _, ok := b.workflows[workflow.Id]; !ok {
		return
	}

	b.workflows[workflow.Id].Stop()
	delete(b.workflows, workflow.Id)

	return
}

func (b *Core) UpdateWorkflowScenario(workflow *m.Workflow) (err error) {

	if _, ok := b.workflows[workflow.Id]; !ok {
		return
	}

	err = b.workflows[workflow.Id].UpdateScenario()

	return
}

func (b *Core) UpdateWorkflow(workflow *m.Workflow) (err error) {

	if workflow.Status == "enabled" {
		if _, ok := b.workflows[workflow.Id]; !ok {
			err = b.AddWorkflow(workflow)
		}
	} else {
		if _, ok := b.workflows[workflow.Id]; ok {
			err = b.DeleteWorkflow(workflow)
		}
	}

	return
}

// ------------------------------------------------
// Flows
// ------------------------------------------------

func (b *Core) AddFlow(flow *m.Flow) (err error) {

	if _, ok := b.workflows[flow.WorkflowId]; !ok {
		return
	}

	if err = b.workflows[flow.WorkflowId].AddFlow(flow); err != nil {
		return
	}

	return
}

func (b *Core) GetFlow(id int64) (*Flow, error) {

	var flow *m.Flow
	var err error
	if flow, err = b.adaptors.Flow.GetById(id); err != nil {
		return nil, err
	}

	if _, ok := b.workflows[flow.WorkflowId]; !ok {
		return nil, nil
	}

	if flow, ok := b.workflows[flow.WorkflowId].Flows[id]; ok {
		return flow, nil
	}

	return nil, nil
}

func (b *Core) UpdateFlow(flow *m.Flow) (err error) {

	if _, ok := b.workflows[flow.WorkflowId]; !ok {
		return
	}

	if err = b.workflows[flow.WorkflowId].UpdateFlow(flow); err != nil {
		return
	}

	return
}

func (b *Core) RemoveFlow(flow *m.Flow) (err error) {

	if _, ok := b.workflows[flow.WorkflowId]; !ok {
		return
	}

	if err = b.workflows[flow.WorkflowId].RemoveFlow(flow); err != nil {
		return
	}

	return
}

// ------------------------------------------------
// Workers
// ------------------------------------------------

func (b *Core) UpdateFlowFromDevice(device *m.Device) (err error) {

	//	var flows map[int64]*m.Flow
	//	flows = make(map[int64]*m.Flow)
	//	childs, _, _ := device.GetChilds()
	//
	//	for _, workflow := range b.workflows {
	//		for _, flow := range workflow.Flows {
	//			for _, worker := range flow.Workers {
	//				for _, action := range worker.actions {
	//					//if action.Device.Id == device.Id {
	//					//	workflow.UpdateFlow(flow.Model)
	//					//	continue
	//					//}
	//
	//					if action.Device != nil && action.Device.Id == device.Id {
	//						//workflow.UpdateFlow(flow.Model)
	//						flows[flow.Model.Id] = flow.Model
	//						continue
	//					}
	//
	//					for _, child := range childs {
	//						if action.Device != nil && action.Device.Id == child.Id {
	//							flows[flow.Model.Id] = flow.Model
	//						}
	//					}
	//				}
	//
	//				if device.Device != nil && worker.Model.DeviceAction.Device.Id == device.Device.Id {
	//					//workflow.UpdateFlow(flow.Model)
	//					flows[flow.Model.Id] = flow.Model
	//					continue
	//				}
	//			}
	//		}
	//
	//		for _, flow := range flows {
	//			workflow.UpdateFlow(flow)
	//		}
	//
	//		flows = make(map[int64]*m.Flow)
	//	}

	return
}

func (b *Core) UpdateWorker(_worker *m.Worker) (err error) {

	//	for _, workflow := range b.workflows {
	//		for _, flow := range workflow.Flows {
	//			for _, worker := range flow.Workers {
	//				if worker.Model.Id == _worker.Id {
	//					flow.UpdateWorker(_worker)
	//					break
	//				}
	//			}
	//		}
	//	}
	//
	return
}

func (b *Core) RemoveWorker(_worker *m.Worker) (err error) {

	//	for _, workflow := range b.workflows {
	//		for _, flow := range workflow.Flows {
	//			for _, worker := range flow.Workers {
	//				if worker.Model.Id == _worker.Id {
	//					flow.RemoveWorker(_worker)
	//					break
	//				}
	//			}
	//		}
	//	}

	return
}

func (b *Core) DoWorker(model *m.Worker) (err error) {

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
