package controllers

import (
	"encoding/json"
	"github.com/e154/smart-home/api/models"
	"strconv"
	"net/url"
	"github.com/astaxie/beego/orm"
)

//  ScenarioController oprations for Scenario
type ScenarioController struct {
	CommonController
}

// URLMapping ...
func (c *ScenarioController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Search", c.Search)
}

// Post ...
// @Title Post
// @Description create Scenario
// @Param	body		body 	models.Scenario	true		"body for Scenario content"
// @Success 201 {int} models.Scenario
// @Failure 403 body is empty
// @router / [post]
func (c *ScenarioController) Post() {
	var v models.Scenario
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if _, err := models.AddScenario(&v); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = map[string]interface{}{"scenario": v}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Scenario by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Scenario
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ScenarioController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	scenario, err := models.GetScenarioById(id)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	o := orm.NewOrm()
	if _, err = o.LoadRelated(scenario, "Scripts", 2);err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"scenario": scenario}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Scenario
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Scenario
// @Failure 403
// @router / [get]
func (c *ScenarioController) GetAll() {
	ml, meta, err := models.GetAllScenario(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	var scenarios []models.Scenario
	o := orm.NewOrm()
	for _, bp := range ml {
		scenario := bp.(models.Scenario)
		if _, err = o.LoadRelated(&scenario, "Scripts", 2);err != nil {
			c.ErrHan(403, err.Error())
			continue
		}
		scenarios = append(scenarios, scenario)
	}

	c.Data["json"] = &map[string]interface{}{"scenarios": scenarios, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Scenario
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Scenario	true		"body for Scenario content"
// @Success 200 {object} models.Scenario
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ScenarioController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Scenario{Id: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateScenarioById(&v); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Scenario
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ScenarioController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteScenario(id); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

func (c *ScenarioController) Search() {

	query, fields, sortby, order, offset, limit := c.pagination()
	link, _ := url.ParseRequestURI(c.Ctx.Request.URL.String())
	q := link.Query()

	if val, ok := q["query"]; ok {
		for _, v := range val {
			query["name__icontains"] = v
		}
	}

	ml, meta, err := models.GetAllScenario(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"scenarios": ml, "meta": meta}
	c.ServeJSON()
}
