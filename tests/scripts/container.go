package scripts

import (
	"github.com/e154/smart-home/system/dig"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/db"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/api/server/v1/controllers"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/config"
)

func BuildContainer() (container *dig.Container) {

	container = dig.New()
	container.Provide(server.NewServer)
	container.Provide(server.NewServerConfig)
	container.Provide(controllers.NewControllersV1)
	//container.Provide(config.ReadConfig)
	container.Provide(graceful_service.NewGracefulService)
	container.Provide(graceful_service.NewGracefulServicePool)
	container.Provide(graceful_service.NewGracefulServiceConfig)
	container.Provide(db.NewOrm)
	container.Provide(db.NewOrmConfig)
	container.Provide(core.NewCore)
	container.Provide(migrations.NewMigrations)
	container.Provide(migrations.NewMigrationsConfig)
	container.Provide(adaptors.NewAdaptors)
	container.Provide(logging.NewLogrus)
	container.Provide(scripts.NewScriptService)
	container.Provide(func() (conf *config.AppConfig, err error) {
		conf, err = config.ReadConfig()
		conf.PgName = "smart_home_test"
		return
	})

	return
}
