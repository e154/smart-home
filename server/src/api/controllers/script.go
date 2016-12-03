package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
)

// ScriptController operations for Script
type ScriptController struct {
	CommonController
}

// URLMapping ...
func (c *ScriptController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Script
// @Param	body		body 	models.Script	true		"body for Script content"
// @Success 201 {object} models.Script
// @Failure 403 body is empty
// @router / [post]
func (c *ScriptController) Post() {
	var script models.Script
	json.Unmarshal(c.Ctx.Input.RequestBody, &script)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&script)
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

	nid, err := models.AddScript(&script)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"id": nid}

	}

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Script by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Script
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ScriptController) GetOne() {
	id, _ := c.GetInt(":id")
	script, err := models.GetScriptById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"script": script}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Script
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Script
// @Failure 403
// @router / [get]
func (c *ScriptController) GetAll() {
	ml, meta, err := models.GetAllScript(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"scripts": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Script
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Script	true		"body for Script content"
// @Success 200 {object} models.Script
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ScriptController) Put() {
	id, _ := c.GetInt(":id")
	var script models.Script
	json.Unmarshal(c.Ctx.Input.RequestBody, &script)
	script.Id = int64(id)
	if err := models.UpdateScriptById(&script); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Script
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ScriptController) Delete() {
	id, _ := c.GetInt(":id")

	if err := models.DeleteScript(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
