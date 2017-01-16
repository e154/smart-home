package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"../models"
	"encoding/json"
)

type BaseController struct {
	CommonController
}

func (b *BaseController) Prepare() {

	if userinfo := b.Ctx.Input.Session("userinfo"); userinfo != nil {
		user := userinfo.(*models.User)
		b.Data["current_user"], _ = json.Marshal(user)
	}

	if token := b.Ctx.Input.Session("token"); token != nil {
		b.Data["token"] = token.(string)
	}

	//b.Data["xsrf_token"] = template.HTML(b.XSRFToken())
	b.Data["domen"] = beego.AppConfig.String("domen")
	b.Data["server_url"] = fmt.Sprintf("%s:%s",beego.AppConfig.String("serveraddr"), beego.AppConfig.String("serverport"))
	b.Data["debug"] = beego.AppConfig.String("runmode") == "dev"
}