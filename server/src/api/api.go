package api

import (
	"fmt"
	"time"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"./routers"
	"./filters"
)

func Initialize() {
	// site base
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_host := beego.AppConfig.String("db_host")
	db_name := beego.AppConfig.String("db_name")
	db_port := beego.AppConfig.String("db_port")
	// parseTime https://github.com/go-sql-driver/mysql#parsetime
	db := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", db_user, db_pass, db_host, db_port, db_name)
	orm.RegisterDataBase("default", "mysql", db, 30, 30)
	// Timezone http://beego.me/docs/mvc/model/orm.md#timezone-config
	orm.DefaultTimeLoc, _ = time.LoadLocation("Asia/Novosibirsk")

	routers.Initialize()

	beego.Info("AppPath:", beego.AppPath)
	if(beego.BConfig.RunMode == "dev") {
		beego.Info("Develment mode enabled")
		// orm debug mode
		if orm_debug, _ := beego.AppConfig.Bool("orm_debug"); orm_debug {
			orm.Debug = true
		}

		//beego.BConfig.WebConfig.DirectoryIndex = true
		//beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		//beego.SetStaticPath("/admin/static", "static_source/admin")
	} else {
		beego.Info("Product mode enabled")
		//beego.SetStaticPath("/admin/static", "www/admin")
	}

	// register access filters
	filters.RegisterFilters()
}
