package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"../models"
	"encoding/json"
	"os"
)

type BaseController struct {
	CommonController
}

func (b *BaseController) Prepare() {

	if userinfo := b.Ctx.Input.Session("userinfo"); userinfo != nil {
		user := userinfo.(*models.User)
		current_user, _ := json.Marshal(user)
		b.Data["current_user"] = string(current_user)
	}

	if token := b.Ctx.Input.Session("token"); token != nil {
		b.Data["token"] = token.(string)
	}

	b.Data["version"] = map[string]string{
		"version": os.Getenv("VERSION"),
		"revision": os.Getenv("REVISION"),
		"revision_url": os.Getenv("REVISION_URL"),
		"generated": os.Getenv("GENERATED"),
		"developers": os.Getenv("DEVELOPERS"),
		"build_number": os.Getenv("BUILD_NUMBER"),
	}

	//b.Data["xsrf_token"] = template.HTML(b.XSRFToken())
	b.Data["domen"] = beego.AppConfig.String("domen")
	b.Data["server_url"] = fmt.Sprintf("%s:%s",beego.AppConfig.String("serveraddr"), beego.AppConfig.String("serverport"))
	b.Data["debug"] = beego.AppConfig.String("runmode") == "dev"
}