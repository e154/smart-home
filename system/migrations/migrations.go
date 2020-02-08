// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

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
	"github.com/e154/smart-home/system/orm"
)

var (
	log = logging.MustGetLogger("migrations")
)

type Migrations struct {
	cfg    *orm.OrmConfig
	source migrate.MigrationSource
	db     *gorm.DB
}

func NewMigrations(cfg *orm.OrmConfig, db *gorm.DB, mConf *MigrationsConfig) *Migrations {

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

	log.Warningf("Purge database: %s", m.cfg.Name)

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
