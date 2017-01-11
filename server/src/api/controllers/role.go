package controllers

import (
	"../models"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"net/url"
)

// RoleController operations for Role
type RoleController struct {
	CommonController
}

// URLMapping ...
func (c *RoleController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Role
// @Param	body		body 	models.Role	true		"body for Role content"
// @Success 201 {object} models.Role
// @Failure 403 body is empty
// @router / [post]
func (c *RoleController) Post() {
	var role models.Role
	json.Unmarshal(c.Ctx.Input.RequestBody, &role)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&role)
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

	_, err = models.AddRole(&role)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"role": role}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Role by name
// @Param	name		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Role
// @Failure 403 :name is empty
// @router /:name [get]
func (c *RoleController) GetOne() {
	name := c.GetString(":name")
	role, err := models.GetRoleByName(name)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"role": role}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Role
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Role
// @Failure 403
// @router / [get]
func (c *RoleController) GetAll() {
	roles, meta, err := models.GetAllRole(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"roles": roles, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Role
// @Param	name		path 	string	true		"The name you want to update"
// @Param	body		body 	models.Role	true		"body for Role content"
// @Success 200 {object} models.Role
// @Failure 403 :name is not int
// @router /:name [put]
func (c *RoleController) Put() {
	name := c.GetString(":name")
	var role models.Role
	json.Unmarshal(c.Ctx.Input.RequestBody, &role)
	role.Name = name
	if err := models.UpdateRoleByName(&role); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Role
// @Param	name		path 	string	true		"The name you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 name is empty
// @router /:name [delete]
func (c *RoleController) Delete() {
	name := c.GetString(":name")
	if err := models.DeleteRole(name); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

func (c *RoleController) Search() {

	query, fields, sortby, order, offset, limit := c.pagination()
	link, _ := url.ParseRequestURI(c.Ctx.Request.URL.String())
	q := link.Query()

	if val, ok := q["query"]; ok {
		for _, v := range val {
			query["name__icontains"] = v
		}
	}

	ml, meta, err := models.GetAllRole(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"roles": ml, "meta": meta}
	c.ServeJSON()
}