package controllers

import (
	"github.com/astaxie/beego/validation"
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/api/core"
	"github.com/e154/smart-home/api/models"
)

// NodeController operations for Node
type NodeController struct {
	CommonController
}

// URLMapping ...
func (c *NodeController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Get", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Node
// @Param	body		body 	models.Node	true		"body for Node content"
// @Success 201 {object} models.Node
// @Failure 403 body is empty
// @router / [post]
func (c *NodeController) Post() {
	var node models.Node
	json.Unmarshal(c.Ctx.Input.RequestBody, &node)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&node)
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

	nid, err := models.AddNode(&node)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"id": nid}

	}

	pm := core.CorePtr()
	node.Id = nid
	pm.AddNode(&node)

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Node by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Node
// @Failure 403 :id is empty
// @router /:id [get]
func (c *NodeController) GetOne() {
	id, _ := c.GetInt(":id")
	node, err := models.GetNodeById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"node": node}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Node
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Node
// @Failure 403
// @router / [get]
func (c *NodeController) GetAll() {

	ml, meta, err := models.GetAllNode(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"nodes": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Node
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Node	true		"body for Node content"
// @Success 200 {object} models.Node
// @Failure 403 :id is not int
// @router /:id [put]
func (c *NodeController) Put() {
	id, _ := c.GetInt(":id")
	var node models.Node
	json.Unmarshal(c.Ctx.Input.RequestBody, &node)
	node.Id = int64(id)
	if err := models.UpdateNodeById(&node); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if err := core.CorePtr().ReloadNode(&node); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Node
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *NodeController) Delete() {
	id, _ := c.GetInt(":id")

	// get node
	node, err := models.GetNodeById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	// remove from process
	if err := core.CorePtr().RemoveNode(node); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if err := models.DeleteNode(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
