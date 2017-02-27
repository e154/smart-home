package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/e154/smart-home/api/models"
)

//  ScenarioScriptController oprations for ScenarioScript
type ScenarioScriptController struct {
	CommonController
}

// URLMapping ...
func (c *ScenarioScriptController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create ScenarioScript
// @Param	body		body 	models.ScenarioScript	true		"body for ScenarioScript content"
// @Success 201 {int} models.ScenarioScript
// @Failure 403 body is empty
// @router / [post]
func (c *ScenarioScriptController) Post() {
	var v models.ScenarioScript
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if _, err := models.AddScenarioScript(&v); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = v
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get ScenarioScript by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ScenarioScript
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ScenarioScriptController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetScenarioScriptById(id)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = v
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get ScenarioScript
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ScenarioScript
// @Failure 403
// @router / [get]
func (c *ScenarioScriptController) GetAll() {
	ml, meta, err := models.GetAllScenario(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"scenario_scripts": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the ScenarioScript
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ScenarioScript	true		"body for ScenarioScript content"
// @Success 200 {object} models.ScenarioScript
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ScenarioScriptController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.ScenarioScript{Id: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateScenarioScriptById(&v); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = "OK"
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the ScenarioScript
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ScenarioScriptController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteScenarioScript(id); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = "OK"
	c.ServeJSON()
}
