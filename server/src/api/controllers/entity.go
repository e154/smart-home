package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
)

// EntityController operations for Entity
type EntityController struct {
	CommonController
}

// URLMapping ...
func (c *EntityController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Entity
// @Param	body		body 	models.Entity	true		"body for Entity content"
// @Success 201 {object} models.Entity
// @Failure 403 body is empty
// @router / [post]
func (c *EntityController) Post() {
	var entity models.Entity
	json.Unmarshal(c.Ctx.Input.RequestBody, &entity)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&entity)
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

	nid, err := models.AddEntity(&entity)
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
// @Description get Entity by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Entity
// @Failure 403 :id is empty
// @router /:id [get]
func (c *EntityController) GetOne() {
	id, _ := c.GetInt(":id")
	entity, err := models.GetEntityById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"entity": entity}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Entity
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Entity
// @Failure 403
// @router / [get]
func (c *EntityController) GetAll() {
	ml, meta, err := models.GetAllEntity(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"entities": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Entity
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Entity	true		"body for Entity content"
// @Success 200 {object} models.Entity
// @Failure 403 :id is not int
// @router /:id [put]
func (c *EntityController) Put() {
	id, _ := c.GetInt(":id")
	var entity models.Entity
	json.Unmarshal(c.Ctx.Input.RequestBody, &entity)
	entity.Id = int64(id)
	if err := models.UpdateEntityById(&entity); err != nil {
		c.ErrHan(403, err.Error())
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Entity
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *EntityController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteEntity(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
	}

	c.ServeJSON()
}
