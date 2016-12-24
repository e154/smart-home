package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
	"github.com/astaxie/beego/orm"
)

// MapController operations for Map
type MapController struct {
	CommonController
}

// URLMapping ...
func (c *MapController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetFull", c.GetFull)
	c.Mapping("Put", c.Put)
	c.Mapping("PutFull", c.PutFull)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Map
// @Param	body		body 	models.Map	true		"body for Map content"
// @Success 201 {object} models.Map
// @Failure 403 body is empty
// @router / [post]
func (c *MapController) Post() {
	var _map models.Map
	json.Unmarshal(c.Ctx.Input.RequestBody, &_map)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&_map)
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

	if _, err = models.AddMap(&_map); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map": _map}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Map by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Map
// @Failure 403 :id is empty
// @router /:id [get]
func (c *MapController) GetOne() {
	id, _ := c.GetInt(":id")
	_map, err := models.GetMapById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map": _map}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Map
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Map
// @Failure 403
// @router / [get]
func (c *MapController) GetAll() {
	ml, meta, err := models.GetAllMap(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"maps": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Map
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Map	true		"body for Map content"
// @Success 200 {object} models.Map
// @Failure 403 :id is not int
// @router /:id [put]
func (c *MapController) Put() {
	id, _ := c.GetInt(":id")
	var _map models.Map
	json.Unmarshal(c.Ctx.Input.RequestBody, &_map)
	_map.Id = int64(id)
	if err := models.UpdateMapById(&_map); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Map
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *MapController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteMap(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

func (c *MapController) GetFull() {
	id, _ := c.GetInt(":id")
	_map, err := models.GetMapById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	o := orm.NewOrm()
	if _, err = o.LoadRelated(_map, "Layers", false); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	for _, layer := range _map.Layers {
		if _, err = o.LoadRelated(layer, "Elements", false); err != nil {
			c.ErrHan(403, err.Error())
			return
		}
	}

	c.Data["json"] = map[string]interface{}{"map": _map}
	c.ServeJSON()
}

func (c *MapController) PutFull() {
	id, _ := c.GetInt(":id")
	var newMap models.Map
	json.Unmarshal(c.Ctx.Input.RequestBody, &newMap)

	// get old map model
	// -----------------------------------
	var err error
	oldMap := &models.Map{Id:int64(id)}
	o := orm.NewOrm()
	if _, err = o.LoadRelated(oldMap, "Layers", false); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	for _, layer := range oldMap.Layers {
		if _, err = o.LoadRelated(layer, "Elements", false); err != nil {
			c.ErrHan(403, err.Error())
			return
		}
	}

	// and compare with new data
	// -----------------------------------
	var exist bool
	for _, oldLayer := range oldMap.Layers {
		exist = false
		for _, newLayer := range newMap.Layers {
			if oldLayer.Id == newLayer.Id {
				exist = true
				break
			}
		}
		if !exist {
			// remove old layers
			if err = models.DeleteMapLayer(oldLayer.Id); err != nil {
				c.ErrHan(403, err.Error())
				return
			}
		}
	}

	// crud operations
	// -----------------------------------
	for _, newLayer := range newMap.Layers {
		if newLayer.Id == 0 {
			// add new
			if _, err = models.AddMapLayer(newLayer); err != nil {
				c.ErrHan(403, err.Error())
				return
			}

		} else {
			// update old
			if err = models.UpdateMapLayerById(newLayer); err != nil {
				c.ErrHan(403, err.Error())
				return
			}
		}
	}

	// update map
	// -----------------------------------
	newMap.Id = int64(id)
	if err := models.UpdateMapById(&newMap); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}