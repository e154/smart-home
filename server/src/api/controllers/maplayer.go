package controllers

import (
	"../models"
)
import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
)

// MapLayerController operations for MapLayer
type MapLayerController struct {
	CommonController
}

// URLMapping ...
func (c *MapLayerController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create MapLayer
// @Param	body		body 	models.MapLayer	true		"body for MapLayer content"
// @Success 201 {object} models.MapLayer
// @Failure 403 body is empty
// @router / [post]
func (c *MapLayerController) Post() {
	var map_layer models.MapLayer
	json.Unmarshal(c.Ctx.Input.RequestBody, &map_layer)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&map_layer)
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

	if _, err = models.AddMapLayer(&map_layer); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map_layer": map_layer}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get MapLayer by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.MapLayer
// @Failure 403 :id is empty
// @router /:id [get]
func (c *MapLayerController) GetOne() {
	id, _ := c.GetInt(":id")
	map_layer, err := models.GetMapLayerById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map": map_layer}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get MapLayer
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.MapLayer
// @Failure 403
// @router / [get]
func (c *MapLayerController) GetAll() {
	ml, meta, err := models.GetAllMapLayer(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"maps": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the MapLayer
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.MapLayer	true		"body for MapLayer content"
// @Success 200 {object} models.MapLayer
// @Failure 403 :id is not int
// @router /:id [put]
func (c *MapLayerController) Put() {
	id, _ := c.GetInt(":id")
	var map_layer models.MapLayer
	json.Unmarshal(c.Ctx.Input.RequestBody, &map_layer)
	map_layer.Id = int64(id)
	if err := models.UpdateMapLayerById(&map_layer); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the MapLayer
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *MapLayerController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteMapLayer(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
