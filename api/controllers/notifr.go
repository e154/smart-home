package controllers

import (
	"strconv"
	"github.com/e154/smart-home/api/models"
	"github.com/astaxie/beego/orm"
)

// NotifrController operations for Notifr
type NotifrController struct {
	CommonController
}

// URLMapping ...
func (c *NotifrController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Notifr
// @Param	body		body 	models.MessageDeliverie	true		"body for Notifr content"
// @Success 201 {object} models.MessageDeliverie
// @Failure 403 body is empty
// @router / [post]
func (c *NotifrController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Notify by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.MessageDeliverie
// @Failure 403 :id is empty
// @router /:id [get]
func (c *NotifrController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	notify, err := models.GetMessageDeliverieById(id)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	o := orm.NewOrm()
	if _, err = o.LoadRelated(notify, "Message");err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"notify": notify}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Notifr
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.MessageDeliverie
// @Failure 403
// @router / [get]
func (c *NotifrController) GetAll() {
	ml, meta, err := models.GetAllMessageDeliverie(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	var notifications []models.MessageDeliverie
	o := orm.NewOrm()
	for _, message := range ml {
		md := message.(models.MessageDeliverie)
		if _, err = o.LoadRelated(&md, "Message", 2);err != nil {
			c.ErrHan(403, err.Error())
			continue
		}
		notifications = append(notifications, md)
	}

	c.Data["json"] = &map[string]interface{}{"notifications": notifications, "meta": meta}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Notifr
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *NotifrController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteMessageDeliverie(id); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
