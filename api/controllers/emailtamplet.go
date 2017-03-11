package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/e154/smart-home/api/models"
	"net/url"
)

type EmailTemplateController struct {
	CommonController
}

func (e *EmailTemplateController) Post() {

	 _, b, valid, err := models.EmailTemplateAddNew(e.Ctx.Input.RequestBody)
	if err != nil {
		e.ErrHan(403, err.Error())
		return
	}

	if !b {
		var msg string
		for _, err := range valid.Errors {
			msg += fmt.Sprintf( "%s: %s\r", err.Key, err.Message)
		}
		e.ErrHan(403, err.Error())
		return
	}

	e.Ctx.Output.SetStatus(201)
	e.ServeJSON()
}

func (e *EmailTemplateController) GetOne() {

	name := e.Ctx.Input.Param(":name")
	template, err := models.EmailTemplateGetByName(name)
	if err != nil {
		e.ErrHan(403, err.Error())
		return
	}

	template.GetMarkers()

	e.Data["json"] = &map[string]interface {}{"template": template}
	e.ServeJSON()
}

// @Title Get All
// @Description get Email Template
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.EmailTemplate
// @Failure 403
// @router / [get]
func (e *EmailTemplateController) GetAll() {
	ml, meta, err := models.GetAllEmailTemplate(e.pagination())
	if err != nil {
		e.ErrHan(403, err.Error())
		return
	}

	templates := []models.EmailItem{}
	for _, m := range ml {
		item := m.(models.EmailItem)
		item.GetMarkers()
		templates = append(templates, item)
	}

	e.Data["json"] = &map[string]interface{}{"templates": templates, "meta":meta}
	e.ServeJSON()
}

func (e *EmailTemplateController) Put() {

	name := e.Ctx.Input.Param(":name")
	b, valid, err := models.EmailTemplateUpdate(e.Ctx.Input.RequestBody, name)
	if err != nil {
		e.ErrHan(403, err.Error())
		return
	}

	if !b {
		var msg string
		for _, err := range valid.Errors {
			msg += fmt.Sprintf( "%s: %s\r", err.Key, err.Message)
		}
		e.ErrHan(403, err.Error())
		return
	}
}

func (e *EmailTemplateController) Delete() {

	name := e.Ctx.Input.Param(":name")
	if err := models.EmailTemplateDelete(name); err != nil {
		e.ErrHan(403, err.Error())
		return
	}
}

func (e *EmailTemplateController) Preview() {

	tpl := new(models.EmailTemplate)
	if err := json.Unmarshal(e.Ctx.Input.RequestBody, tpl); err != nil {
		e.ErrHan(403, err.Error())
		return
	}

	if len(tpl.Items) == 0 {
		return
	}

	buf, err := models.EmailTemplatePreview(tpl)
	if err != nil {
		e.ErrHan(403, err.Error())
		return
	}

	e.Ctx.Output.Body([]byte(buf))
}

func (c *EmailTemplateController) Search() {

	query, fields, sortby, order, offset, limit := c.pagination()
	link, _ := url.ParseRequestURI(c.Ctx.Request.URL.String())
	q := link.Query()

	if val, ok := q["query"]; ok {
		for _, v := range val {
			query["name__icontains"] = v
		}
	}

	ml, meta, err := models.GetAllEmailTemplate(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	templates := []models.EmailItem{}
	for _, m := range ml {
		item := m.(models.EmailItem)
		item.GetMarkers()
		templates = append(templates, item)
	}

	c.Data["json"] = &map[string]interface{}{"templates": templates, "meta":meta}
	c.ServeJSON()
}
