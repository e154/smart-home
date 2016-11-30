package controllers

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/orm"
	"../models"
	"../core"
)

// DeviceController operations for Device
type DeviceController struct {
	CommonController
}

// URLMapping ...
func (c *DeviceController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetGroup", c.GetGroup)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Get", c.GetAll)
	c.Mapping("Get", c.GetActions)
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

	core.CorePtr().UpdateWorkerFromDevice(&device)

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
	}

	o := orm.NewOrm()

	var count int64
	if count, err = o.QueryTable(device).Filter("device_id", id).Count(); err != nil {
		return
	}

	if count > 0 {
		device.IsGroup = true
	}

	c.Data["json"] = map[string]interface{}{"device": device}

	c.ServeJSON()
}

// GetGroup ...
// @Title GetGroup
// @Description get Device by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Device
// @Failure 403 :id is empty
// @router / [get]
func (c *DeviceController) GetGroup() {

	o := orm.NewOrm()

	devices := []*models.Device{}
	if _, err := o.QueryTable(&models.Device{}).Filter("address__isnull", true).Filter("device_id__isnull", true).RelatedSel("Node").All(&devices); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	for _, device := range devices {
		device.IsGroup = true
	}

	c.Data["json"] = map[string]interface{}{"devices": devices}

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

	o := orm.NewOrm()
	devices := []models.Device{}
	for _, m := range ml {
		device := m.(models.Device)
		devices = append(devices, device)
	}


	var count int64
	for key, device := range devices {

		if count, err = o.QueryTable(device).Filter("device_id", device.Id).Count(); err != nil {
			return
		}

		if count > 0 {
			devices[key].IsGroup = true
		}
		count = 0
	}

	c.Data["json"] = &map[string]interface{}{"devices": devices, "meta": meta}
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

	core.CorePtr().UpdateWorkerFromDevice(&device)

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
	device, err := models.GetDeviceById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	core.CorePtr().UpdateWorkerFromDevice(device)

	if err := models.DeleteDevice(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// GetActions ...
// @Title GetActions
// @Description get device Actions by id
// @Param	id		path 	string	true
// @Success 200 {object} models.Device
// @Failure 403 :id is empty
// @router /:id/actions [get]
func (c *DeviceController) GetActions() {
	id, _ := c.GetInt(":id")
	device, err := models.GetDeviceById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	ids := []int64{int64(id)}
	log.Println("device",device)
	if device.Device != nil {
		ids = append(ids, device.Device.Id)
	}

	actions, err := models.GetDeviceActionsByDeviceId(ids)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"actions": actions}
	c.ServeJSON()
}