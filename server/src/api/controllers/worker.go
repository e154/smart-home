package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/orm"
	"../models"
	"../core"
)

// WorkerController operations for Worker
type WorkerController struct {
	CommonController
}

// URLMapping ...
func (c *WorkerController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Worker
// @Param	body		body 	models.Worker	true		"body for Worker content"
// @Success 201 {object} models.Worker
// @Failure 403 body is empty
// @router / [post]
func (c *WorkerController) Post() {
	var worker models.Worker
	json.Unmarshal(c.Ctx.Input.RequestBody, &worker)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&worker)
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

	nid, err := models.AddWorker(&worker)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	} else {
		c.Data["json"] = map[string]interface{}{"id": nid}

	}

	core.CorePtr().AddWorker(&worker)

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Worker by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Worker
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WorkerController) GetOne() {
	id, _ := c.GetInt(":id")
	worker, err := models.GetWorkerById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	o := orm.NewOrm()
	if _, err = o.LoadRelated(worker, "Flow"); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if _, err = o.LoadRelated(worker, "WorkFlow"); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"worker": worker}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Worker
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Worker
// @Failure 403
// @router / [get]
func (c *WorkerController) GetAll() {
	ml, meta, err := models.GetAllWorker(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"workers": ml, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Worker
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Worker	true		"body for Worker content"
// @Success 200 {object} models.Worker
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WorkerController) Put() {
	id, _ := c.GetInt(":id")
	var worker models.Worker
	json.Unmarshal(c.Ctx.Input.RequestBody, &worker)
	worker.Id = int64(id)
	if err := models.UpdateWorkerById(&worker); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	core.CorePtr().UpdateWorker(&worker)

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Worker
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WorkerController) Delete() {
	id, _ := c.GetInt(":id")

	worker, err  := models.GetWorkerById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if err := models.DeleteWorker(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	core.CorePtr().RemoveWorker(worker)

	c.ServeJSON()
}
