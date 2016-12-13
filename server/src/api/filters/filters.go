package filters

import (
	"../log"
)

func RegisterFilters() {
	log.Println("Filters initialize...")

	 //https://github.com/astaxie/beego/issues/1294
	//var filer = func(ctx *context.Context) {
	//	ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	//	ctx.ResponseWriter.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	//}

	//beego.InsertFilter("*", beego.BeforeRouter, filer)

	//beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
	//	AllowAllOrigins: true,
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "UPDATE"},
	//	AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
	//	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
	//	AllowCredentials: true,
	//}))

	// register rbac filter
	//rbac.AccessRegister()
}
