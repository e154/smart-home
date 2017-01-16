package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/astaxie/beego"
	"../models"
)

type DashboardController struct {
	BaseController
}

// URLMapping ...
func (c *DashboardController) URLMapping() {
	c.Mapping("Index", c.Index)
	c.Mapping("Signin", c.Signin)
	c.Mapping("Signout", c.Signout)
	c.Mapping("Recovery", c.Recovery)
	c.Mapping("Reset", c.Reset)
}

func (h *DashboardController) Index() {

	userinfo := h.GetSession("userinfo")
	if userinfo == nil {
		h.Ctx.Redirect(302, "/signin")
		return
	}

	h.Layout = h.GetTemplate() + "/private/base.tpl.html"
	h.TplName = h.GetTemplate() + "/private/index.tpl.html"
	h.UpdateTemplate()
}

func (h *DashboardController) Signin() {

	userinfo := h.GetSession("userinfo")
	if userinfo != nil {
		h.Ctx.Redirect(302, "/")
		return
	}

	if isajax := h.Ctx.Input.IsAjax(); isajax {

		var err error
		input := map[string]string{}
		if err = json.Unmarshal(h.Ctx.Input.RequestBody, &input); err != nil {
			h.ErrHan(403, err.Error())
			return
		}

		server_url := fmt.Sprintf("%s:%s/api/v1/signin", beego.AppConfig.String("serveraddr"), beego.AppConfig.String("serverport"))

		j, _ := json.Marshal(input)
		result, err := h.SendRequest("POST", server_url, j)
		if err != nil {
			h.ErrHan(403, err.Error())
			return
		}

		signin := &models.Signin{}
		if err = json.Unmarshal(result, signin); err != nil {
			h.ErrHan(403, err.Error())
			return
		}

		if signin.User != nil {
			h.SetSession("userinfo", signin.User)
			h.SetSession("token", signin.Token)
			h.ServeJSON()
			return
		}

		return
	}

	h.Layout = h.GetTemplate() + "/public/base.tpl.html"
	h.TplName = h.GetTemplate() + "/public/login.tpl.html"
	h.UpdateTemplate()
}

func (h *DashboardController) Signout() {
	h.DelSession("userinfo")
	h.DelSession("token")
	h.Ctx.Redirect(302, "/signin")
}

func (h *DashboardController) Recovery() {}

func (h *DashboardController) Reset() {}