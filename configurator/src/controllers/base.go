package controllers

import (
	"github.com/astaxie/beego"
	"html/template"
)

type BaseController struct {
	CommonController
}

func (b *BaseController) Prepare() {
	// scripts
	// -------------------------------
	b.Data["Scripts"] = []string{
		"/static/js/lib.min.js",
		"/static/js/app.min.js",
		"/static/js/templates.min.js",
	}

	// styles
	// -------------------------------
	b.Data["Styles"] = []string{
		"/static/css/lib.min.css",
		"/static/css/app.min.css",
	}

	b.Data["xsrf_token"] = template.HTML(b.XSRFToken())

	// get user role
	// -------------------------------
	//...

	b.Data["domen"] = beego.AppConfig.String("domen")
}