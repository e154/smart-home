package controllers

import (
	"../../api"
	"github.com/astaxie/beego"
	"encoding/json"
	"../../polygraphy"
	"../../common"
)

type HomeController struct {
	BaseController
}

func (h *HomeController) Index() {

	h.Data["Title"] = st.SiteName
	h.Data["Menu"] = string(json_menu)
	h.Data["Settings"] = ""
	h.Layout = h.GetTemplate() + "/frontend/base.tpl.html"
	h.TplName = h.GetTemplate() + "/frontend/index.tpl.html"
	beego.BuildTemplate(beego.BConfig.WebConfig.ViewsPath)
	h.Render()
}
