package controllers

import (
	"encoding/json"
	"github.com/e154/smart-home/api/models"
	"strconv"
	"net/url"
)

//  WorkflowScenarioController oprations for WorkflowScenario
type WorkflowScenarioController struct {
	CommonController
}

// URLMapping ...
func (c *WorkflowScenarioController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create WorkflowScenario
// @Param	body		body 	models.WorkflowScenario	true		"body for WorkflowScenario content"
// @Success 201 {int} models.WorkflowScenario
// @Failure 403 body is empty
// @router / [post]
func (c *WorkflowScenarioController) Post() {
	workflow_id, _ := c.GetInt(":workflow_id")
	var v models.WorkflowScenario
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	v.Workflow = &models.Workflow{Id: int64(workflow_id)}
	if _, err := models.AddWorkflowScenario(&v); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Ctx.Output.SetStatus(201)

	c.Data["json"] = map[string]interface{}{"scenario": v}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get WorkflowScenario by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.WorkflowScenario
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WorkflowScenarioController) GetOne() {
	workflow_id, _ := c.GetInt(":workflow_id")
	workflow, err := models.GetWorkflowById(int64(workflow_id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	id, _ := c.GetInt(":id")
	workflow_scenario, err := workflow.GetScenarioById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if _, err = workflow_scenario.GetScripts(); err != nil {
		return
	}

	c.Data["json"] = map[string]interface{}{"scenario": workflow_scenario}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get WorkflowScenario
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.WorkflowScenario
// @Failure 403
// @router / [get]
func (c *WorkflowScenarioController) GetAll() {
	workflow_id, _ := c.GetInt(":workflow_id")
	workflow, err := models.GetWorkflowById(int64(workflow_id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if _, err = workflow.GetScenarios(); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"scenarios": workflow.Scenarios}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the WorkflowScenario
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.WorkflowScenario	true		"body for WorkflowScenario content"
// @Success 200 {object} models.WorkflowScenario
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WorkflowScenarioController) Put() {
	workflow_id, _ := c.GetInt(":workflow_id")
	workflow, err := models.GetWorkflowById(int64(workflow_id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	id, _ := c.GetInt(":id")
	scenario, err := workflow.GetScenarioById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, scenario)
	if err := models.UpdateWorkflowScenarioById(scenario); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	_scripts := scenario.Scripts
	scenario, err = workflow.GetScenarioById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if _, err := scenario.UpdateScripts(_scripts); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the WorkflowScenario
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WorkflowScenarioController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteWorkflowScenario(id); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

func (c *WorkflowScenarioController) Search() {

	workflow_id, _ := c.GetInt(":workflow_id")
	query, fields, sortby, order, offset, limit := c.pagination()
	link, _ := url.ParseRequestURI(c.Ctx.Request.URL.String())
	q := link.Query()

	if val, ok := q["query"]; ok {
		for _, v := range val {
			query["name__icontains"] = v
			query["workflow_id__exact"] = strconv.Itoa(workflow_id)
		}
	}

	ml, meta, err := models.GetAllWorkflowScenario(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"scenarios": ml, "meta": meta}
	c.ServeJSON()
}