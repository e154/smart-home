package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/lib/common"
	"net/url"
	"github.com/e154/smart-home/api/log"
)

// UserController operations for User
type UserController struct {
	CommonController
}

// URLMapping ...
func (c *UserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("UpdateStatus", c.UpdateStatus)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {object} models.User
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Post() {

	user := &models.User{}
	json.Unmarshal(c.Ctx.Input.RequestBody, user)

	incoming := map[string]string{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &incoming)

	// construct user data
	//---------------------------------------------------------
	status := incoming["status"]
	if len(status) > 0 && (status == "active" || status == "blocked") {
		user.Status = status
	} else {
		user.Status = "blocked"
	}

	password := incoming["password"]
	if len(password) >= 6 && password == incoming["password_repeat"] {
		user.EncryptedPassword = common.Pwdhash(password)
	}

	// validation
	//---------------------------------------------------------
	valid := validation.Validation{}
	b, err := valid.Valid(user)
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

	// create user
	//---------------------------------------------------------
	_, err = models.AddUser(user)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	for _, meta := range user.Meta {
		user.SetMeta(meta.Key, meta.Value)
	}

	if err = user.LoadRelated(); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"user": user}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get User by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UserController) GetOne() {
	id, _ := c.GetInt(":id")
	user, err := models.GetUserById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if err = user.LoadRelated(); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"user": user}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 403
// @router / [get]
func (c *UserController) GetAll() {
	users, meta, err := models.GetAllUser(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"users": users, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UserController) Put() {
	id, _ := c.GetInt(":id")
	var user *models.User
	var err error

	if user, err = models.GetUserById(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if err = json.Unmarshal(c.Ctx.Input.RequestBody, user); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	incoming := map[string]interface{}{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &incoming); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	// construct user data
	//---------------------------------------------------------
	user.Id = int64(id)
	password, ok := incoming["password"].(string)
	password_repeat, ok := incoming["password_repeat"].(string)
	if ok &&  password != "" && len(password) < 6  {
		c.ErrHan(403, "The password should be at least six characters.")
		return
	}

	if ok && password == password_repeat {
		log.Warnf("account: updated password for '%s'", user.Email)
		user.EncryptedPassword = common.Pwdhash(password)
	} else {
		c.ErrHan(403, "Passwords do not match.")
		return
	}

	// update user
	//---------------------------------------------------------
	if err := models.UpdateUserById(user); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	for _, meta := range user.Meta {
		user.SetMeta(meta.Key, meta.Value)
	}

	c.ServeJSON()
}

func (c *UserController) UpdateStatus() {
	var user, _user *models.User
	var err error

	user = &models.User{}

	id, _ := c.GetInt(":id")
	json.Unmarshal(c.Ctx.Input.RequestBody, user)

	if _user, err = models.GetUserById(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	_user.Status = user.Status
	if err = models.UpdateUserById(_user); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if err = _user.LoadRelated(); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"user": _user}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UserController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteUser(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	log.Warnf("account: removed user with id %d", id)

	c.ServeJSON()
}

func (c *UserController) Search() {

	query, fields, sortby, order, offset, limit := c.pagination()
	link, _ := url.ParseRequestURI(c.Ctx.Request.URL.String())
	q := link.Query()

	if val, ok := q["query"]; ok {
		for _, v := range val {
			query["nickname__icontains"] = v
		}
	}

	ml, meta, err := models.GetAllUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"users": ml, "meta": meta}
	c.ServeJSON()
}