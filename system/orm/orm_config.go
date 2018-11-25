package orm

import (
	"fmt"
	"github.com/e154/smart-home/system/config"
)

type OrmConfig struct {
	Alias           string
	Name            string
	User            string
	Password        string
	Host            string
	Port            string
	Debug           bool
	Logger          bool
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifeTime int
}

func (c OrmConfig) String() string {

	// parseTime https://github.com/go-sql-driver/mysql#parsetime
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable", c.Name, c.User, c.Password, c.Host, c.Port)
}

func NewOrmConfig(cfg *config.AppConfig) *OrmConfig {
	return &OrmConfig{
		Alias:           "default",
		Name:            cfg.PgName,
		User:            cfg.PgUser,
		Password:        cfg.PgPass,
		Host:            cfg.PgHost,
		Port:            cfg.PgPort,
		Debug:           cfg.PgDebug,
		Logger:          cfg.PgLogger,
		MaxIdleConns:    cfg.PgMaxIdleConns,
		MaxOpenConns:    cfg.PgMaxOpenConns,
		ConnMaxLifeTime: cfg.PgConnMaxLifeTime,
	}
}
