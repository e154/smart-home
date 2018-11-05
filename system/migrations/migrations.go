package migrations

import (
	"database/sql"
	"os"
	"github.com/op/go-logging"
	"github.com/rubenv/sql-migrate"
	. "github.com/e154/smart-home/system/migrations/assets"
	"github.com/jinzhu/gorm"
	"path"
	"fmt"
	"github.com/e154/smart-home/db"
)

var (
	log = logging.MustGetLogger("migrations")
)

type Migrations struct {
	cfg    *db.OrmConfig
	source migrate.MigrationSource
	db     *gorm.DB
}

func NewMigrations(cfg *db.OrmConfig, db *gorm.DB, mConf *MigrationsConfig) *Migrations {

	var source migrate.MigrationSource

	switch mConf.Source {
	case "assets", "":
		source = &migrate.AssetMigrationSource{
			Asset:    Asset,
			AssetDir: AssetDir,
			Dir:      mConf.Dir,
		}
	case "dir":
		source = &migrate.FileMigrationSource{
			Dir: path.Join(mConf.Dir),
		}
	default:
		panic(fmt.Sprintf("unknown source %s", mConf.Source))
	}

	return &Migrations{
		cfg:    cfg,
		source: source,
		db:     db,
	}
}

func (m Migrations) Connect() (sqlDb *sql.DB, err error) {
	sqlDb, err = sql.Open("postgres", m.cfg.String())
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	return
}

func (m Migrations) Up() (err error) {

	var sqlDb *sql.DB
	sqlDb, err = m.Connect()
	defer sqlDb.Close()

	var n int
	if n, err = migrate.Exec(sqlDb, "postgres", m.source, migrate.Up); err != nil {
		log.Error(err.Error())
	}

	log.Infof("Applied %d migrations!", n)

	return
}

func (m Migrations) Down() (err error) {

	var sqlDb *sql.DB
	sqlDb, err = m.Connect()
	defer sqlDb.Close()

	var n int
	if n, err = migrate.Exec(sqlDb, "postgres", m.source, migrate.Down); err != nil {
		log.Error(err.Error())
	}

	log.Infof("Applied %d migrations!", n)

	return
}

func (m Migrations) Purge() (err error) {

	log.Warning("Purge database")

	if err = m.db.Exec(`DROP SCHEMA IF EXISTS "public" CASCADE;`).Error; err != nil {
		log.Error(err.Error())
		return
	}
	if err = m.db.Exec(`CREATE SCHEMA "public";`).Error; err != nil {
		log.Error(err.Error())
		return
	}

	err = m.Up()

	return
}
