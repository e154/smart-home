package controllers

import (
	"github.com/astaxie/beego"
)

type DashboardController struct {
	BaseController
}

func (h *DashboardController) Index() {

	h.Layout = h.GetTemplate() + "/frontend/base.tpl.html"
	h.TplName = h.GetTemplate() + "/frontend/index.tpl.html"
	beego.BuildTemplate(beego.BConfig.WebConfig.ViewsPath)
	h.Render()
}
