package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
	"../bpms"
)

// WorkflowController operations for Workflow
type WorkflowController struct {
	CommonController
}

// URLMapping ...
func (c *WorkflowController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Get", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Workflow
// @Param	body		body 	models.Workflow	true		"body for Workflow content"
// @Success 201 {object} models.Workflow
// @Failure 403 body is empty
// @router / [post]
func (c *WorkflowController) Post() {
	var workflow models.Workflow
	json.Unmarshal(c.Ctx.Input.RequestBody, &workflow)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&workflow)
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

	nid, err := models.AddWorkflow(&workflow)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"id": nid}

	}

	bpms.BpmsPtr().AddWorkflow(&workflow)

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Workflow by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Workflow
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WorkflowController) GetOne() {
	id, _ := c.GetInt(":id")
	node, err := models.GetWorkerById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"workflow": node}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Workflow
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Workflow
// @Failure 403
// @router / [get]
func (c *WorkflowController) GetAll() {
	ml, meta, err := models.GetAllWorkflow(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"workflows": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Workflow
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Workflow	true		"body for Workflow content"
// @Success 200 {object} models.Workflow
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WorkflowController) Put() {
	id, _ := c.GetInt(":id")
	var workflow models.Workflow
	json.Unmarshal(c.Ctx.Input.RequestBody, &workflow)
	workflow.Id = int64(id)
	if err := models.UpdateWorkflowById(&workflow); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Workflow
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WorkflowController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteWorkflow(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
