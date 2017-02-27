package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"github.com/e154/smart-home/api/models"
)

// DashboardController operations for Dashboard
type DashboardController struct {
	CommonController
}

// URLMapping ...
func (c *DashboardController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Dashboard
// @Param	body		body 	models.Dashboard	true		"body for Dashboard content"
// @Success 201 {object} models.Dashboard
// @Failure 403 body is empty
// @router / [post]
func (c *DashboardController) Post() {
	var board models.Dashboard
	json.Unmarshal(c.Ctx.Input.RequestBody, &board)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&board)
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

	_, err = models.AddDashboard(&board)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"dashboard": board}
	c.Ctx.Output.SetStatus(201)
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Dashboard by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Dashboard
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DashboardController) GetOne() {
	id, _ := c.GetInt(":id")
	board, err := models.GetDashboardById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"dashboard": board}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Dashboard
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Dashboard
// @Failure 403
// @router / [get]
func (c *DashboardController) GetAll() {
	boards, meta, err := models.GetAllDashboard(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"dashboards": boards, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Dashboard
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Dashboard	true		"body for Dashboard content"
// @Success 200 {object} models.Dashboard
// @Failure 403 :id is not int
// @router /:id [put]
func (c *DashboardController) Put() {
	id, _ := c.GetInt(":id")
	var board models.Dashboard
	json.Unmarshal(c.Ctx.Input.RequestBody, &board)
	board.Id = int64(id)
	if err := models.UpdateDashboardById(&board); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Dashboard
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DashboardController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteDashboard(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
