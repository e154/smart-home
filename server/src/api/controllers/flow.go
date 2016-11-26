package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"../models"
	"../bpms"
	"github.com/astaxie/beego/orm"
)

// FlowController operations for Flow
type FlowController struct {
	CommonController
}

// URLMapping ...
func (c *FlowController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetOneFull", c.GetOneFull)
	c.Mapping("GetWorkers", c.GetWorkers)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Flow
// @Param	body		body 	models.Flow	true		"body for Flow content"
// @Success 201 {object} models.Flow
// @Failure 403 body is empty
// @router / [post]
func (c *FlowController) Post() {
	var flow models.Flow
	json.Unmarshal(c.Ctx.Input.RequestBody, &flow)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&flow)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if !b {
		var msg string
		for _, err := range valid.Errors {
			msg += fmt.Sprintf("%s: %s\r", err.Key, err.Message)
		}
		c.ErrHan(403, msg)
		return
	}
	//....

	nid, err := models.AddFlow(&flow)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"id": nid}

	}

	bpms.BpmsPtr().AddFlow(&flow)

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Flow by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Flow
// @Failure 403 :id is empty
// @router /:id [get]
func (c *FlowController) GetOne() {
	id, _ := c.GetInt(":id")
	flow, err := models.GetFlowById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"flow": flow}
	}

	workers, err := models.GetWorkersByFlowId(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}
	flow.Workers = workers

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Flow by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Flow
// @Failure 403 :id is empty
// @router /:id [get]
func (c *FlowController) GetOneFull() {
	id, _ := c.GetInt(":id")
	flow, err := models.GetFullFlowById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"flow": flow}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Flow
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Flow
// @Failure 403
// @router / [get]
func (c *FlowController) GetAll() {
	ml, flows, err := models.GetAllFlow(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"flows": ml, "meta": flows}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Flow
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Flow	true		"body for Flow content"
// @Success 200 {object} models.Flow
// @Failure 403 :id is not int
// @router /:id [put]
func (c *FlowController) Put() {
	id, _ := c.GetInt(":id")
	var flow models.Flow
	json.Unmarshal(c.Ctx.Input.RequestBody, &flow)
	flow.Id = int64(id)
	if err := models.UpdateFlowById(&flow); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	bpms.BpmsPtr().UpdateFlow(&flow)

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Flow
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *FlowController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteFlow(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	bpms.BpmsPtr().RemoveFlow(&models.Flow{Id: int64(id)})

	c.ServeJSON()
}

func (c *FlowController) GetOneRedactor() {
	id, _ := c.GetInt(":id")
	flow, err := models.GetFlowById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	o := orm.NewOrm()
	_, err = o.LoadRelated(flow, "Workflow")
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	var r *models.RedactorFlow
	r, err = ExportToRedactor(flow)

	workers, err := models.GetWorkersByFlowId(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	r.Workers = workers

	c.Data["json"] = map[string]interface{}{"flow": r}
	c.ServeJSON()
}

func (c *FlowController) UpdateRedactor() {
	var flow models.RedactorFlow
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &flow); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	newFlow := &models.Flow{
		Id: flow.Id,
		Name: flow.Name,
		Description: flow.Description,
		Status: flow.Status,
		Workflow: &models.Workflow{Id:flow.Workflow.Id},
	}
	if err := models.UpdateFlowById(newFlow); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	var err error
	var flowElements []*models.FlowElement
	// update flow elements
	//---------------------------------------------------
	if flowElements, err = models.GetFlowElementsByFlow(&models.Flow{Id:flow.Id}); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	flow_todo_remove := []*models.FlowElement{}
	for _, element := range flowElements {
		exist := false
		for _, object := range flow.Objects {
			if object.Id == element.Uuid {
				exist = true
				break
			}
		}

		if !exist {
			flow_todo_remove = append(flow_todo_remove, element)
		}
	}

	// remove flow elements
	for _, element := range flow_todo_remove {
		if err = models.DeleteFlowElement(element.Uuid); err != nil {
			c.ErrHan(403, err.Error())
			return
		}
	}

	var j []byte
	// insert or update flow elements
	for _, element := range flow.Objects {

		if j, err = json.Marshal(element.Position); err != nil {
			c.ErrHan(403, err.Error())
			return
		}

		fl := &models.FlowElement{
			Uuid: element.Id,
			Name: element.Title,
			GraphSettings: fmt.Sprintf("{\"position\":%s}", j),
			Status: element.Status,
			FlowId: flow.Id,
		}

		switch element.Type.Name {
		case "event":
			if element.Type.Start != nil {
				fl.PrototypeType = "MessageHandler"
			} else if element.Type.End != nil {
				fl.PrototypeType = "MessageEmitter"
			}
		case "task":
			fl.PrototypeType = "Task"
		default:
			fl.PrototypeType = "default"
		}

		if _, err = models.AddOrUpdateFlowElement(fl); err != nil {
			c.ErrHan(403, err.Error())
			return
		}
	}

	// connectors
	//---------------------------------------------------
	var connections []*models.Connection
	if connections, err = models.GetConnectionsByFlow(&models.Flow{Id:flow.Id}); err != nil {
		return
	}

	conn_todo_remove := []*models.Connection{}
	for _, old_conn := range connections {
		exist := false
		for _, new_conn := range flow.Connectors {
			if old_conn.Uuid == new_conn.Id {
				exist = true
				break
			}
		}

		if !exist {
			conn_todo_remove = append(conn_todo_remove, old_conn)
		}
	}

	for _, conn := range conn_todo_remove {
		if err = models.DeleteConnection(conn.Uuid); err != nil {
			c.ErrHan(403, err.Error())
			return
		}
	}

	for _, connector := range flow.Connectors {

		conn := &models.Connection{
			Uuid: connector.Id,
			Name: connector.Title,
			ElementFrom: connector.Start.Object,
			ElementTo: connector.End.Object,
			PointFrom: connector.Start.Point,
			PointTo: connector.End.Point,
			FlowId: flow.Id,
			GraphSettings: "",
		}

		if _, err = models.AddOrUpdateConnection(conn); err != nil {
			c.ErrHan(403, err.Error())
			return
		}
	}

	// workers
	//---------------------------------------------------
	//var workers []*models.Worker


	//---------------------------------------------------
	newflow, err := models.GetFlowById(flow.Id)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	var r *models.RedactorFlow
	r, err = ExportToRedactor(newflow)

	bpms.BpmsPtr().UpdateFlow(&models.Flow{Id:flow.Id})

	c.Data["json"] = map[string]interface{}{"flow": r}
	c.ServeJSON()

}

func ExportToRedactor(f *models.Flow) (flow *models.RedactorFlow, err error) {

	flow = &models.RedactorFlow{
		Id: f.Id,
		Name: f.Name,
		Status: f.Status,
		Description: f.Description,
		Workflow: f.Workflow,
		Objects: make([]*models.RedactorObject, 0),
		Connectors: make([]*models.RedactorConnector, 0),
		Created_at: f.Created_at,
		Update_at: f.Update_at,
	}

	var flowElements []*models.FlowElement
	if flowElements, err = models.GetFlowElementsByFlow(f); err != nil {
		return
	}

	for _, el := range flowElements {
		object := &models.RedactorObject{
			Id: el.Uuid,
			Title: el.Name,
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
		default:

		}

		gst := new(models.RedactorGrapSettings)
		if err = json.Unmarshal([]byte(el.GraphSettings), &gst); err != nil {
			return
		}

		object.Position = gst.Position

		flow.Objects = append(flow.Objects, object)
	}

	var connections []*models.Connection
	if connections, err = models.GetConnectionsByFlow(f); err != nil {
		return
	}

	for _, con := range connections {
		connector := &models.RedactorConnector{
			Id: con.Uuid,
			Flow_type: "default",
			Title: con.Name,
		}
		connector.Start.Object = con.ElementFrom
		connector.Start.Point = con.PointFrom

		connector.End.Object = con.ElementTo
		connector.End.Point = con.PointTo

		flow.Connectors = append(flow.Connectors, connector)
	}

	return
}

func (c *FlowController) GetWorkers() {
	id, _ := c.GetInt(":id")
	workers, err := models.GetWorkersByFlowId(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"workers": workers}
	c.ServeJSON()
}