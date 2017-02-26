package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/core"
	"net/url"
)

// DeviceActionController operations for DeviceAction
type DeviceActionController struct {
	CommonController
}

// URLMapping ...
func (c *DeviceActionController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Search", c.Search)
}

// Post ...
// @Title Create
// @Description create DeviceAction
// @Param	body		body 	models.DeviceAction	true		"body for DeviceAction content"
// @Success 201 {object} models.DeviceAction
// @Failure 403 body is empty
// @router / [post]
func (c *DeviceActionController) Post() {
	var action models.DeviceAction
	json.Unmarshal(c.Ctx.Input.RequestBody, &action)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&action)
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

	nid, err := models.AddDeviceAction(&action)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"id": nid}

	}
	c.Ctx.Output.SetStatus(201)
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get DeviceAction by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DeviceAction
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DeviceActionController) GetOne() {
	id, _ := c.GetInt(":id")
	action, err := models.GetDeviceActionById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"action": action}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get DeviceAction
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.DeviceAction
// @Failure 403
// @router / [get]
func (c *DeviceActionController) GetAll() {

	ml, meta, err := models.GetAllDeviceAction(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"actions": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the DeviceAction
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.DeviceAction	true		"body for DeviceAction content"
// @Success 200 {object} models.DeviceAction
// @Failure 403 :id is not int
// @router /:id [put]
func (c *DeviceActionController) Put() {

	id, _ := c.GetInt(":id")
	var action models.DeviceAction
	json.Unmarshal(c.Ctx.Input.RequestBody, &action)
	action.Id = int64(id)
	if err := models.UpdateDeviceActionById(&action); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	// restart worker
	workers, err := models.GetWorkersByDeviceAction(&action)
	for _, worker := range workers {
		if err = core.CorePtr().UpdateWorker(worker); err != nil {
			c.ErrHan(403, err.Error())
			return
		}
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the DeviceAction
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DeviceActionController) Delete() {
	id, _ := c.GetInt(":id")

	// restart worker
	workers, err := models.GetWorkersByDeviceAction(&models.DeviceAction{Id: int64(id)})
	for _, worker := range workers {
		if err = core.CorePtr().RemoveWorker(worker); err != nil {
			c.ErrHan(403, err.Error())
			return
		}

		if err = models.DeleteWorker(worker.Id); err != nil {
			c.ErrHan(403, err.Error())
			return
		}
	}

	if err := models.DeleteDeviceAction(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

func (c *DeviceActionController) Search() {

	query, fields, sortby, order, offset, limit := c.pagination()
	link, _ := url.ParseRequestURI(c.Ctx.Request.URL.String())
	q := link.Query()

	if val, ok := q["query"]; ok {
		for _, v := range val {
			query["name__icontains"] = v
			query["description__icontains"] = v
		}
	}

	ml, meta, err := models.GetAllDeviceAction(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"actions": ml, "meta": meta}
	c.ServeJSON()
}

// GetByDevice ...
// @Title GetByDevice
// @Description get DeviceState by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DeviceState
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DeviceActionController) GetByDevice() {
	id, _ := c.GetInt(":id")
	device_states, err := models.GetAllDeviceActionByDevice(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"device_actions": device_states}
	c.ServeJSON()
}