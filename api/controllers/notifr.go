package controllers

import (
	"strconv"
	"github.com/e154/smart-home/api/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"github.com/e154/smart-home/api/notifr"
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
	c.Mapping("Repeat", c.Repeat)
}

// Post ...
// @Title Create
// @Description create Notifr
// @Param	body		body 	models.MessageDeliverie	true		"body for Notifr content"
// @Success 201 {object} models.MessageDeliverie
// @Failure 403 body is empty
// @router / [post]
func (c *NotifrController) Post() {

	type NewMessage struct {
		Address		string	`json:"address"`
		Title		string	`json:"title"`
		Body		string	`json:"body"`
		Template	*models.EmailItem	`json:"template"`
		Params		map[string]interface{}	`json:"params"`
		Type		string	`json:"type"`
	}

	message := &NewMessage{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, message); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	switch message.Type {
	case "email":
		email := notifr.NewEmail()
		email.To = message.Address
		email.Body = message.Body
		email.Params = message.Params
		email.Subject = message.Title
		if message.Template != nil {
			email.Template = message.Template.Name
		}
		notifr.Send(email)

	case "sms":
		sms := notifr.NewSms()
		sms.Body = message.Body
		sms.To = message.Address
		notifr.Send(sms)
	case "push":
		push := notifr.NewPush()
		push.Body = message.Body
		push.To = message.Address
		notifr.Send(push)
	}

	c.Ctx.Output.SetStatus(201)
	c.ServeJSON()
}

func (c NotifrController) Repeat() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	notifr.Repeat(id)

	c.Ctx.Output.SetStatus(201)
	c.ServeJSON()
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
