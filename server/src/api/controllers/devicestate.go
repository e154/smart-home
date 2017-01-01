package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
)

// DeviceStateController operations for DeviceState
type DeviceStateController struct {
	CommonController
}

// URLMapping ...
func (c *DeviceStateController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create DeviceState
// @Param	body		body 	models.DeviceState	true		"body for DeviceState content"
// @Success 201 {object} models.DeviceState
// @Failure 403 body is empty
// @router / [post]
func (c *DeviceStateController) Post() {
	var device_state models.DeviceState
	json.Unmarshal(c.Ctx.Input.RequestBody, &device_state)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&device_state)
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

	if _, err = models.AddDeviceState(&device_state); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"device_state": device_state}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get DeviceState by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DeviceState
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DeviceStateController) GetOne() {
	id, _ := c.GetInt(":id")
	device_state, err := models.GetDeviceStateById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"device_state": device_state}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get DeviceState
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.DeviceState
// @Failure 403
// @router / [get]
func (c *DeviceStateController) GetAll() {
	ml, meta, err := models.GetAllDeviceState(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"device_states": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the DeviceState
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.DeviceState	true		"body for DeviceState content"
// @Success 200 {object} models.DeviceState
// @Failure 403 :id is not int
// @router /:id [put]
func (c *DeviceStateController) Put() {
	id, _ := c.GetInt(":id")
	var device_state models.DeviceState
	json.Unmarshal(c.Ctx.Input.RequestBody, &device_state)
	device_state.Id = int64(id)
	if err := models.UpdateDeviceStateById(&device_state); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the DeviceState
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DeviceStateController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteDeviceState(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// GetByDevice ...
// @Title GetByDevice
// @Description get DeviceState by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DeviceState
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DeviceStateController) GetByDevice() {
	id, _ := c.GetInt(":id")
	device_states, err := models.GetAllDeviceStateByDevice(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"device_states": device_states}
	c.ServeJSON()
}