package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
)

// DeviceController operations for Device
type DeviceController struct {
	CommonController
}

// URLMapping ...
func (c *DeviceController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Get", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Device
// @Param	body		body 	models.Device	true		"body for Device content"
// @Success 201 {object} models.Device
// @Failure 403 body is empty
// @router / [post]
func (c *DeviceController) Post() {
	var device models.Device
	json.Unmarshal(c.Ctx.Input.RequestBody, &device)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&device)
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

	nid, err := models.AddDevice(&device)
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
// @Description get Device by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Device
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DeviceController) GetOne() {
	id, _ := c.GetInt(":id")
	device, err := models.GetDeviceById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"device": device}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Device
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Device
// @Failure 403
// @router / [get]
func (c *DeviceController) GetAll() {

	ml, meta, err := models.GetAllDevice(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"devices": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Device
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Device	true		"body for Device content"
// @Success 200 {object} models.Device
// @Failure 403 :id is not int
// @router /:id [put]
func (c *DeviceController) Put() {
	id, _ := c.GetInt(":id")
	var device models.Device
	json.Unmarshal(c.Ctx.Input.RequestBody, &device)
	device.Id = int64(id)
	if err := models.UpdateDeviceById(&device); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Device
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DeviceController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteDevice(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
