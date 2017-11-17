package api

import (
	"fmt"
	"time"
	"io/ioutil"
	"path/filepath"
	"github.com/astaxie/beego/validation"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/e154/smart-home/api/routers"
	"github.com/e154/smart-home/api/filters"
	"github.com/e154/smart-home/api/core"
	"github.com/e154/smart-home/api/cron"
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/api/notifr"
	"github.com/e154/smart-home/api/telemetry"
	"github.com/e154/smart-home/api/variable"
)

func configuration() {

	// check if exist data dir
	data_dir := beego.AppConfig.String("data_dir")
	if _, err := ioutil.ReadDir(data_dir); err != nil {
		panic("data directory not found")
		return
	}

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

	// run migration
	Migration(db)

	if(beego.BConfig.RunMode == "dev") {
		if orm_debug, _ := beego.AppConfig.Bool("orm_debug"); orm_debug {
			orm.Debug = true
		}
	}

	log.Info("AppPath:", beego.AppPath)
	if(beego.BConfig.RunMode == "dev") {
		log.Info("Development mode enabled")
	} else {
		log.Info("Product mode enabled")
		beego.BConfig.ServerName = "smart-home"
	}

	file_storage_path := beego.AppConfig.String("file_storage_path")
	beego.SetStaticPath("/static", filepath.Join(data_dir, file_storage_path))

	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.EnableXSRF = false
	beego.BConfig.CopyRequestBody = true

	validation.SetDefaultMessage(map[string]string{
		"Required":     "Должно быть заполнено",
		"Min":          "Минимально допустимое значение %d",
		"Max":          "Максимально допустимое значение %d",
		"Range":        "Должно быть в диапазоне от %d до %d",
		"MinSize":      "Минимально допустимая длина %d",
		"MaxSize":      "Максимально допустимая длина %d",
		"Length":       "Длина должна быть равна %d",
		"Alpha":        "Должно состоять из букв",
		"Numeric":      "Должно состоять из цифр",
		"AlphaNumeric": "Должно состоять из букв или цифр",
		"Match":        "Должно совпадать с %s",
		"NoMatch":      "Не должно совпадать с %s",
		"AlphaDash":    "Должно состоять из букв, цифр или символов (-_)",
		"Email":        "Должно быть в правильном формате email",
		"IP":           "Должен быть правильный IP адрес",
		"Base64":       "Должно быть представлено в правильном формате base64",
		"Mobile":       "Должно быть правильным номером мобильного телефона",
		"Tel":          "Должно быть правильным номером телефона",
		"Phone":        "Должно быть правильным номером телефона или мобильного телефона",
		"ZipCode":      "Должно быть правильным почтовым индексом",
	})

	// register access filters
	filters.RegisterFilters()
}

func Initialize() {

	configuration()

	// init settings
	variable.Initialize()

	// routes
	routers.Initialize()

	// cron
	cron.Initialize()

	// notifr
	notifr.Initialize()

	// telemetry
	t := telemetry.Initialize()

	// core
	if err := core.Initialize(t); err != nil {
		log.Error(err.Error())
	}

	log.Info("Starting....")

	// rest api
	go beego.Run()
}
