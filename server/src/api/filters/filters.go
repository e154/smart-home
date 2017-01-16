package filters

import (
	"../log"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"../rbac"
)

func RegisterFilters() {
	log.Info("Filters initialize...")

	// CORS
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "UPDATE"},
		AllowCredentials: true,
	}))

	// register rbac access filter
	rbac.AccessFilter()
}
