package use_case

import (
	"errors"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/uuid"
	"github.com/e154/smart-home/system/core"
)

func AddFlow(newFlow *models.NewFlowModel, adaptors *adaptors.Adaptors, core *core.Core) (id int64, errs []*validation.Error, err error) {

	flow := &m.Flow{}
	if err = common.Copy(&flow, &newFlow); err != nil {
		return
	}

	if newFlow.Workflow.Id != 0 {
		flow.WorkflowId = newFlow.Workflow.Id
	}

	if newFlow.Scenario.Id != 0 {
		flow.WorkflowScenarioId = newFlow.Scenario.Id
	}

	_, errs = flow.Valid()
	if len(errs) > 0 {
		return
	}

	if id, err = adaptors.Flow.Add(flow); err != nil {
		return
	}

	err = core.AddFlow(flow)

	return
}
func GetFlowById(flowId int64, adaptors *adaptors.Adaptors) (flowDto *models.FlowModel, err error) {

	var flow *m.Flow
	if flow, err = adaptors.Flow.GetById(flowId); err != nil {
		return
	}

	flowDto = &models.FlowModel{}
	err = common.Copy(&flowDto, &flow, common.JsonEngine)

	return
}

func GetFlowRedactor(flowId int64, adaptors *adaptors.Adaptors) (redactorFlowDto *models.RedactorFlowModel, err error) {

	var flow *m.Flow
	if flow, err = adaptors.Flow.GetById(flowId); err != nil {
		return
	}

	var redactorFlow *m.RedactorFlow
	if redactorFlow, err = ExportToRedactor(flow, adaptors); err != nil {
		return
	}

	redactorFlowDto = &models.RedactorFlowModel{}
	err = common.Copy(&redactorFlowDto, &redactorFlow, common.JsonEngine)

	return
}

func GetFlowList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (listDto []*models.FlowListModel, total int64, err error) {

	var list []*m.Flow
	if list, total, err = adaptors.Flow.List(limit, offset, order, sortBy); err != nil {
		return
	}

	listDto = make([]*models.FlowListModel, 0)
	err = common.Copy(&listDto, &list)

	return
}

func UpdateFlowRedactor(params *models.RedactorFlowModel,
	adaptors *adaptors.Adaptors,
	core *core.Core) (result *models.RedactorFlowModel, errs []*validation.Error, err error) {

	//debug.Println(params)
	//fmt.Println("--------")

	var oldFlow *m.Flow
	if oldFlow, err = adaptors.Flow.GetById(params.Id); err != nil {
		return
	}

	newFlow := &m.Flow{}
	if err = common.Copy(&newFlow, &params, common.JsonEngine); err != nil {
		return
	}
	if params.Scenario != nil {
		newFlow.WorkflowScenarioId = params.Scenario.Id
	}
	if params.Workflow != nil {
		newFlow.WorkflowId = params.Workflow.Id
	}

	_, errs = newFlow.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.Flow.Update(newFlow); err != nil {
		return
	}

	//update flow elements
	flowTodoRemove := make([]uuid.UUID, 0)
	for _, element := range oldFlow.FlowElements {
		exist := false
		for _, object := range params.Objects {
			if object.Id == element.Uuid.String() {
				exist = true
				break
			}
		}

		if !exist {
			flowTodoRemove = append(flowTodoRemove, element.Uuid)
		}
	}

	if len(flowTodoRemove) > 0 {
		if err = adaptors.FlowElement.Delete(flowTodoRemove); err != nil {
			return
		}
	}

	// insert or update flow elements
	for _, element := range params.Objects {

		fl := &m.FlowElement{}
		common.Copy(&fl, &element)
		common.Copy(&fl.GraphSettings.Position, &element.Position)
		fl.Uuid.Scan(element.Id)
		fl.FlowId = newFlow.Id
		fl.Name = element.Title

		if element.FlowLink != nil && element.FlowLink.Id != 0 {
			fl.FlowLink = &element.FlowLink.Id
		}

		if element.Script != nil {
			fl.ScriptId = &element.Script.Id
		}

		switch element.Type.Name {
		case "event":
			if element.Type.Start != nil {
				fl.PrototypeType = common.FlowElementsPrototypeMessageHandler
			} else if element.Type.End != nil {
				fl.PrototypeType = common.FlowElementsPrototypeMessageEmitter
			}
		case "task":
			fl.PrototypeType = common.FlowElementsPrototypeTask
		case "gateway":
			fl.PrototypeType = common.FlowElementsPrototypeGateway
		case "flow":
			fl.PrototypeType = common.FlowElementsPrototypeFlow
		default:
			fl.PrototypeType = common.FlowElementsPrototypeDefault
		}

		_, errs = fl.Valid()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%s %s", err.Key, err.Message)
			}
			return
		}

		if err = adaptors.FlowElement.AddOrUpdateFlowElement(fl); err != nil {
			return
		}
	}

	// connectors
	connTodoRemove := make([]uuid.UUID, 0)
	for _, oldConn := range oldFlow.Connections {
		exist := false
		for _, newConn := range params.Connectors {
			if oldConn.Uuid.String() == newConn.Id {
				exist = true
				break
			}
		}

		if !exist {
			connTodoRemove = append(connTodoRemove, oldConn.Uuid)
		}
	}

	if len(connTodoRemove) > 0 {
		if err = adaptors.Connection.Delete(connTodoRemove); err != nil {
			return
		}
	}

	for _, c := range params.Connectors {

		conn := &m.Connection{
			Name:      c.Title,
			PointFrom: c.Start.Point,
			PointTo:   c.End.Point,
			FlowId:    newFlow.Id,
			Direction: c.Direction,
		}
		conn.Uuid.Scan(c.Id)
		conn.ElementFrom.Scan(c.Start.Object)
		conn.ElementTo.Scan(c.End.Object)

		_, errs = conn.Valid()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%s %s", err.Key, err.Message)
			}
			return
		}

		if err = adaptors.Connection.AddOrUpdateConnection(conn); err != nil {
			return
		}
	}

	// workers
	workersTodoRemove := make([]*m.Worker, 0)
	for _, oldWorker := range oldFlow.Workers {
		exist := false
		for _, newWorker := range params.Workers {
			if newWorker.Id == oldWorker.Id {
				exist = true
				break
			}
		}

		if !exist {
			workersTodoRemove = append(workersTodoRemove, oldWorker)
		}
	}

	for _, worker := range workersTodoRemove {
		if err = core.RemoveWorker(worker); err == nil {
			if err = adaptors.Worker.Delete([]int64{worker.Id}); err != nil {
				return
			}
		}
	}

	for _, w := range params.Workers {
		worker := &m.Worker{}
		common.Copy(&worker, &w)
		worker.WorkflowId = newFlow.Workflow.Id
		worker.FlowId = newFlow.Id
		worker.DeviceActionId = w.DeviceAction.Id

		_, errs = worker.Valid()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%s %s", err.Key, err.Message)
			}
			return
		}

		if worker.Id == 0 {
			if _, err = adaptors.Worker.Add(worker); err != nil {
				return
			}
		} else {
			if err = adaptors.Worker.Update(worker); err != nil {
				return
			}
		}
	}

	// exit
	if newFlow, err = adaptors.Flow.GetById(params.Id); err != nil {
		return
	}

	if err = core.UpdateFlow(newFlow); err != nil {
		return
	}

	var redactorFlow *m.RedactorFlow
	if redactorFlow, err = ExportToRedactor(newFlow, adaptors); err != nil {
		return
	}

	result = &models.RedactorFlowModel{}
	err = common.Copy(&result, &redactorFlow)

	return
}

func SearchFlow(query string, limit, offset int, adaptors *adaptors.Adaptors) (listDto []*models.FlowModel, total int64, err error) {
	var list []*m.Flow
	if list, total, err = adaptors.Flow.Search(query, limit, offset); err != nil {
		return
	}

	listDto = make([]*models.FlowModel, 0)
	err = common.Copy(&listDto, &list)
	return
}

func DeleteFlowById(flowId int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	var flow *m.Flow
	if flow, err = adaptors.Flow.GetById(flowId); err != nil {
		return
	}

	if err = core.RemoveFlow(flow); err != nil {
		return
	}

	err = adaptors.Flow.Delete(flowId)
	return
}

func UpdateFlow(params *models.UpdateFlowModel,
	adaptors *adaptors.Adaptors,
	core *core.Core) (result *models.FlowModel, errs []*validation.Error, err error) {

	flow := &m.Flow{}
	if err = common.Copy(&flow, &params); err != nil {
		return
	}

	_, errs = flow.Valid()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Errorf("%s %s", err.Key, err.Message)
		}
		return
	}

	if _, err = adaptors.Flow.GetById(flow.Id); err != nil {
		return
	}

	if err = adaptors.Flow.Update(flow); err != nil {
		return
	}

	err = core.UpdateFlow(flow)

	result = &models.FlowModel{}
	common.Copy(&result, &flow)

	return
}

func ExportToRedactor(f *m.Flow, adaptors *adaptors.Adaptors) (redactorFlow *m.RedactorFlow, err error) {

	if f == nil {
		err = errors.New("ExportToRedactor: Nil point")
		return
	}

	var scenario *m.WorkflowScenario
	if scenario, err = adaptors.WorkflowScenario.GetById(f.WorkflowScenarioId); err != nil {
		return
	}

	redactorFlow = &m.RedactorFlow{
		Id:          f.Id,
		Name:        f.Name,
		Status:      f.Status,
		Description: f.Description,
		Workflow:    f.Workflow,
		Scenario:    scenario,
		Workers:     make([]*m.Worker, 0),
		Objects:     make([]*m.RedactorObject, 0),
		Connectors:  make([]*m.RedactorConnector, 0),
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}

	// elements
	for _, el := range f.FlowElements {
		object := &m.RedactorObject{
			Id:            el.Uuid,
			Title:         el.Name,
			Description:   el.Description,
			PrototypeType: el.PrototypeType,
			Script:        el.Script,
		}

		if el.FlowLink != nil {
			var flow *m.Flow
			if flow, err = adaptors.Flow.GetById(*el.FlowLink); err != nil {
				return
			}
			object.FlowLink = flow
		}

		switch el.PrototypeType {
		case "MessageHandler":
			object.Type.Name = "event"
			object.Type.Start = map[int64]interface{}{0: &map[int64]interface{}{0: true}}
		case "MessageEmitter":
			object.Type.Name = "event"
			object.Type.End = map[string]interface{}{"simply": &map[string]interface{}{"top_level": true}}
		case "Task":
			object.Type.Name = "task"
		case "Flow":
			object.Type.Name = "flow"
		case "Gateway":
			object.Type.Name = "gateway"
			object.Type.Start = map[int64]interface{}{0: &map[int64]interface{}{0: true}}
		default:

		}

		common.Copy(&object.Position, &el.GraphSettings.Position)
		redactorFlow.Objects = append(redactorFlow.Objects, object)
	}

	// connections
	for _, con := range f.Connections {
		connector := &m.RedactorConnector{
			Id:        con.Uuid,
			FlowType:  "default",
			Title:     con.Name,
			Direction: con.Direction,
		}
		connector.Start.Object = con.ElementFrom
		connector.Start.Point = con.PointFrom

		connector.End.Object = con.ElementTo
		connector.End.Point = con.PointTo

		redactorFlow.Connectors = append(redactorFlow.Connectors, connector)
	}

	// workers
	common.Copy(&redactorFlow.Workers, &f.Workers)

	return
}
