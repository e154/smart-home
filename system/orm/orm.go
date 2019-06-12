package orm

import (
	_ "github.com/lib/pq"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/op/go-logging"
)

type Orm struct {
	cfg *OrmConfig
	db  *gorm.DB
}

var (
	log = logging.MustGetLogger("orm")
)

func NewOrm(cfg *OrmConfig,
	graceful *graceful_service.GracefulService) (orm *Orm, db *gorm.DB) {

	log.Debugf("database connect %s", cfg.String())
	var err error
	db, err = gorm.Open("postgres", cfg.String())
	if err != nil {
		panic(err.Error())
	}

	db.LogMode(cfg.Logger)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	db.DB().SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.DB().SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifeTime) * time.Minute)

	orm = &Orm{
		cfg: cfg,
		db:  db,
	}

	graceful.Subscribe(orm)
	return
}

func (o *Orm) Shutdown() {
	if o.db != nil {
		log.Debugf("database shutdown")
		o.db.Close()
	}
}
