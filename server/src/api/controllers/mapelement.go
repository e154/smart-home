package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
)

// MapElementController operations for MapElement
type MapElementController struct {
	CommonController
}

// URLMapping ...
func (c *MapElementController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create MapElement
// @Param	body		body 	models.MapElement	true		"body for MapElement content"
// @Success 201 {object} models.MapElement
// @Failure 403 body is empty
// @router / [post]
func (c *MapElementController) Post() {
	var entity models.MapElement
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

	if _, err = models.AddMapElement(&entity); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map_element": entity}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get MapElement by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.MapElement
// @Failure 403 :id is empty
// @router /:id [get]
func (c *MapElementController) GetOne() {
	id, _ := c.GetInt(":id")
	entity, err := models.GetMapElementById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map_element": entity}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get MapElement
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.MapElement
// @Failure 403
// @router / [get]
func (c *MapElementController) GetAll() {
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
// @Description update the MapElement
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.MapElement	true		"body for MapElement content"
// @Success 200 {object} models.MapElement
// @Failure 403 :id is not int
// @router /:id [put]
func (c *MapElementController) Put() {
	id, _ := c.GetInt(":id")
	var mapElement models.MapElement
	json.Unmarshal(c.Ctx.Input.RequestBody, &mapElement)
	mapElement.Id = int64(id)

	// load old data
	//
	var oldMapElement *models.MapElement
	var err error
	if oldMapElement, err = models.GetMapElementById(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	//
	switch mapElement.PrototypeType {
	case "image":

		prototype := &models.MapImage{}
		b, _ := json.Marshal(mapElement.Prototype)
		json.Unmarshal(b, &prototype)
		if prototype != nil {
			if prototype.Id == 0 {
				if _, err := models.AddMapImage(prototype); err != nil {
					c.ErrHan(403, err.Error())
					return
				}
			} else {
				if err := models.UpdateMapImageById(prototype); err != nil {
					c.ErrHan(403, err.Error())
					return
				}
			}

			mapElement.PrototypeId = prototype.Id
			mapElement.PrototypeType = "image"
		} else {
			mapElement.PrototypeId = 0
			mapElement.PrototypeType = ""
		}

	}

	// update mam element
	//
	if err := models.UpdateMapElementById(&mapElement); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	// delete old map image
	//
	if oldMapElement.PrototypeId == 0 {
		oldMapElement.PrototypeType = ""
	}

	switch oldMapElement.PrototypeType {
	case "image":
		models.DeleteMapImage(oldMapElement.PrototypeId)
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the MapElement
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *MapElementController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteMapElement(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
