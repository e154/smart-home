// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"fmt"
	"net/http"
	"path"

	orm2 "github.com/e154/smart-home/internal/system/orm"
	"github.com/e154/smart-home/pkg/logger"

	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"

	"github.com/e154/smart-home/migrations"
)

var (
	log = logger.MustGetLogger("migrations")
)

// Migrations ...
type Migrations struct {
	cfg    *orm2.Config
	source migrate.MigrationSource
	orm    *orm2.Orm
	db     *gorm.DB
}

// NewMigrations ...
func NewMigrations(cfg *orm2.Config,
	db *gorm.DB,
	orm *orm2.Orm,
	mConf *Config) *Migrations {

	var source migrate.MigrationSource

	switch mConf.Source {
	case "embed":
		source = &migrate.HttpFileSystemMigrationSource{FileSystem: http.FS(migrations.F)}
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
		orm:    orm,
		db:     db,
	}
}

// Up ...
func (m Migrations) Up() (err error) {

	var n int
	if n, err = migrate.Exec(m.orm.DB(), "postgres", m.source, migrate.Up); err != nil {
		log.Error(err.Error())
	}

	log.Infof("Applied %d migrations!", n)

	return
}

// Down ...
func (m Migrations) Down() (err error) {

	var n int
	if n, err = migrate.Exec(m.orm.DB(), "postgres", m.source, migrate.Down); err != nil {
		log.Error(err.Error())
	}

	fmt.Printf("Down %d migrations!\n", n)

	return
}

// Purge ...
func (m Migrations) Purge() (err error) {

	fmt.Printf("Drop database: %s\n", m.cfg.Name)

	if err = m.db.Exec(`DROP SCHEMA IF EXISTS "public" CASCADE;`).Error; err != nil {
		log.Error(err.Error())
		return
	}
	if err = m.db.Exec(`CREATE SCHEMA "public";`).Error; err != nil {
		log.Error(err.Error())
		return
	}

	_ = m.orm.Check()

	err = m.Up()

	return
}
