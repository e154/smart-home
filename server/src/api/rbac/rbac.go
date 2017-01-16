package rbac

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
	"fmt"
)

func AccessFilter() {
	beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {

		// for dev
		if beego.BConfig.RunMode == "dev" {
			//return
		}

		requestURI := strings.ToLower(ctx.Request.RequestURI)
		method := strings.ToLower(ctx.Request.Method)
		token := ctx.Input.Header("Authenticate")

		//TODO remove
		fmt.Println("requestURI",requestURI)
		fmt.Println("method",method)
		fmt.Println("token",token)
	})
}

func Initialize() {}