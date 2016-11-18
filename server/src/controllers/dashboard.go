package controllers

import (
	"github.com/astaxie/beego"
)

type DashboardController struct {
	BaseController
}

func (h *DashboardController) Index() {
	h.Layout = h.GetTemplate() + "/base.tpl.html"
	h.TplName = h.GetTemplate() + "/index.tpl.html"
	beego.BuildTemplate(beego.BConfig.WebConfig.ViewsPath)
	h.Render()
}
