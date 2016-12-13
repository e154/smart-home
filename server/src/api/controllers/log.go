package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
)

// LogController operations for Log
type LogController struct {
	CommonController
}

// URLMapping ...
func (c *LogController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Log
// @Param	body		body 	models.Log	true		"body for Log content"
// @Success 201 {object} models.Log
// @Failure 403 body is empty
// @router / [post]
func (c *LogController) Post() {
	var log models.Log
	json.Unmarshal(c.Ctx.Input.RequestBody, &log)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&log)
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

	_, err = models.AddLog(&log)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"log": log}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Log by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Log
// @Failure 403 :id is empty
// @router /:id [get]
func (c *LogController) GetOne() {
	id, _ := c.GetInt(":id")
	log, err := models.GetLogById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"log": log}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Log
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Log
// @Failure 403
// @router / [get]
func (c *LogController) GetAll() {
	ml, meta, err := models.GetAllLog(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"logs": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Log
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Log	true		"body for Log content"
// @Success 200 {object} models.Log
// @Failure 403 :id is not int
// @router /:id [put]
func (c *LogController) Put() {
	id, _ := c.GetInt(":id")
	var log models.Log
	json.Unmarshal(c.Ctx.Input.RequestBody, &log)
	log.Id = int64(id)
	if err := models.UpdateLogById(&log); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Log
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *LogController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteLog(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
