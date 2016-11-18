package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"../models"
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
	c.Mapping("Get", c.GetAll)
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

	//bpms.BpmsPtr().AddFlow(&flow)

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
	ml, meta, err := models.GetAllFlow(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"flows": ml, "meta": meta}
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

	//bpms.BpmsPtr().RemoveWorkflow(&models.Workflow{Id:int64(id)})

	c.ServeJSON()
}
