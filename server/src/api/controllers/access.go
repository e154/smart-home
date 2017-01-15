package controllers

import (
	"../models"
)

// AccessController operations for Access
type AccessController struct {
	CommonController
}

// URLMapping ...
func (c *AccessController) URLMapping() {
	c.Mapping("GetOne", c.Get)
}

func (c *AccessController) Get() {
	c.Data["json"] = &map[string]interface{}{"access_list": models.AccessConfigList}
	c.ServeJSON()
}
