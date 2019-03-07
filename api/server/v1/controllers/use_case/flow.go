package use_case

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"encoding/json"
	"errors"
	"github.com/e154/smart-home/common/debug"
	"fmt"
)

func GetFlowById(flowId int64, adaptors *adaptors.Adaptors) (flow *m.Flow, err error) {

	flow, err = adaptors.Flow.GetById(flowId)

	return
}

func GetFlowRedactor(flowId int64, adaptors *adaptors.Adaptors) (redactorFlow *m.RedactorFlow, err error) {

	var flow *m.Flow
	if flow, err = adaptors.Flow.GetById(flowId); err != nil {
		return
	}

	redactorFlow, err = ExportToRedactor(flow, adaptors)

	return
}

func GetFlowList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*m.Flow, total int64, err error) {

	items, total, err = adaptors.Flow.List(limit, offset, order, sortBy)

	return
}

func UpdateFlowRedactor(redactor *m.RedactorFlow, adaptors *adaptors.Adaptors) (result *m.RedactorFlow, err error) {

	debug.Println(redactor)
	fmt.Println("--------")

	newFlow := &m.Flow{
		Id:                 redactor.Id,
		Name:               redactor.Name,
		Description:        redactor.Description,
		Status:             redactor.Status,
		WorkflowId:         redactor.Workflow.Id,
		WorkflowScenarioId: redactor.Scenario.Id,
	}

	if err = adaptors.Flow.Update(newFlow); err != nil {
		return
	}

	return
}

func SearchFlow(query string, limit, offset int, adaptors *adaptors.Adaptors) (flows []*m.Flow, total int64, err error) {

	flows, total, err = adaptors.Flow.Search(query, limit, offset)

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
		Objects:     make([]*m.RedactorObject, 0),
		Connectors:  make([]*m.RedactorConnector, 0),
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}

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

		gst := &m.RedactorGrapSettings{}
		if err = json.Unmarshal([]byte(el.GraphSettings), &gst); err != nil {
			return
		}

		object.Position = gst.Position

		redactorFlow.Objects = append(redactorFlow.Objects, object)
	}

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

	return
}
