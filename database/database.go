package database

import (
	"os"
	"time"
	"fmt"
	"database/sql"
	"github.com/e154/smart-home/api/log"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/e154/smart-home/api/models"
	migrate "github.com/rubenv/sql-migrate"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func Initialize(testMode bool) string {

	config := GetDbConfig(testMode)

	orm.RegisterDataBase("default", "mysql", config, 30, 30)
	// Timezone http://beego.me/docs/mvc/model/orm.md#timezone-config
	orm.DefaultTimeLoc, _ = time.LoadLocation("Asia/Novosibirsk")

	if(beego.BConfig.RunMode == "dev" && !testMode) {
		if orm_debug, _ := beego.AppConfig.Bool("orm_debug"); orm_debug {
			orm.Debug = true
		}
	}

	orm.RegisterModel(new(models.Connection))
	orm.RegisterModel(new(models.Dashboard))
	orm.RegisterModel(new(models.Device))
	orm.RegisterModel(new(models.DeviceAction))
	orm.RegisterModel(new(models.DeviceState))
	orm.RegisterModel(new(models.EmailItem))
	orm.RegisterModel(new(models.Flow))
	orm.RegisterModel(new(models.FlowElement))
	orm.RegisterModel(new(models.Image))
	orm.RegisterModel(new(models.Log))
	orm.RegisterModel(new(models.Map))
	orm.RegisterModel(new(models.MapDevice))
	orm.RegisterModel(new(models.MapDeviceAction))
	orm.RegisterModel(new(models.MapDeviceState))
	orm.RegisterModel(new(models.MapElement))
	orm.RegisterModel(new(models.MapImage))
	orm.RegisterModel(new(models.MapLayer))
	orm.RegisterModel(new(models.MapText))
	orm.RegisterModel(new(models.Message))
	orm.RegisterModel(new(models.MessageDeliverie))
	orm.RegisterModel(new(models.Node))
	orm.RegisterModel(new(models.Permission))
	orm.RegisterModel(new(models.Role))
	orm.RegisterModel(new(models.Script))
	orm.RegisterModel(new(models.User))
	orm.RegisterModel(new(models.UserMeta))
	orm.RegisterModel(new(models.Variable))
	orm.RegisterModel(new(models.Worker))
	orm.RegisterModel(new(models.Workflow))
	orm.RegisterModel(new(models.WorkflowScenario))
	orm.RegisterModel(new(models.WorkflowScenarioScript))
	orm.RegisterModel(new(models.WorkflowScript))

	return config
}

func Migration(mConn string) {

	db, err := sql.Open("mysql", mConn)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// OR: Read migrations from a folder:
	//migrations := &migrate.FileMigrationSource{
	//	Dir: "database/migrations",
	//}

	// OR: Use migrations from bindata:
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "database/migrations",
	}

	if _, err := migrate.Exec(db, "mysql", migrations, migrate.Up); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	//log.Infof("Applied %d migrations!", n)
}

func GetDbConfig(testMode bool) string {

	var db_user, db_pass, db_host, db_name, db_port string

	if testMode {
		var err error
		config, err := beego.AppConfig.GetSection("test")
		if err != nil {
			panic(err.Error())
		}

		db_user = config["db_user"]
		db_pass = config["db_pass"]
		db_host = config["db_host"]
		db_name = config["db_name"]
		db_port = config["db_port"]
	} else {

		db_user = beego.AppConfig.String("db_user")
		db_pass = beego.AppConfig.String("db_pass")
		db_host = beego.AppConfig.String("db_host")
		db_name = beego.AppConfig.String("db_name")
		db_port = beego.AppConfig.String("db_port")
	}

	// parseTime https://github.com/go-sql-driver/mysql#parsetime
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", db_user, db_pass, db_host, db_port, db_name)
}

func DropDb(db_name string) {

	beego.Info("Drop database", db_name)

	commands := []string{
		`SET FOREIGN_KEY_CHECKS = 0;`,
		`SET @tables = NULL;`,
		`SELECT GROUP_CONCAT(table_schema, '.', table_name) INTO @tables FROM information_schema.tables WHERE table_schema = '`+db_name+`';`,
		`SET @tables = CONCAT('DROP TABLE IF EXISTS ', @tables);`,
		`PREPARE stmt FROM @tables;`,
		`EXECUTE stmt;`,
		`DEALLOCATE PREPARE stmt;`,
		`SET FOREIGN_KEY_CHECKS = 1;`,
	}

	o := orm.NewOrm()
	for _, command := range commands {
		if _, err := o.Raw(command).Exec(); err != nil {
			beego.Error(err.Error())
		}
	}
}