package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
)

// MapEntityController operations for MapEntity
type MapEntityController struct {
	CommonController
}

// URLMapping ...
func (c *MapEntityController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create MapEntity
// @Param	body		body 	models.MapEntity	true		"body for MapEntity content"
// @Success 201 {object} models.MapEntity
// @Failure 403 body is empty
// @router / [post]
func (c *MapEntityController) Post() {
	var entity models.MapEntity
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

	if _, err = models.AddMapEntity(&entity); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map_entity": entity}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get MapEntity by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.MapEntity
// @Failure 403 :id is empty
// @router /:id [get]
func (c *MapEntityController) GetOne() {
	id, _ := c.GetInt(":id")
	entity, err := models.GetMapEntityById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map_entity": entity}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get MapEntity
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.MapEntity
// @Failure 403
// @router / [get]
func (c *MapEntityController) GetAll() {
	ml, meta, err := models.GetAllMap(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"map_entities": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the MapEntity
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.MapEntity	true		"body for MapEntity content"
// @Success 200 {object} models.MapEntity
// @Failure 403 :id is not int
// @router /:id [put]
func (c *MapEntityController) Put() {
	id, _ := c.GetInt(":id")
	var entity models.MapEntity
	json.Unmarshal(c.Ctx.Input.RequestBody, &entity)
	entity.Id = int64(id)
	if err := models.UpdateMapEntityById(&entity); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the MapEntity
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *MapEntityController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteMapEntity(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
