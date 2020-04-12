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

package orm

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"time"
)

// Orm ...
type Orm struct {
	cfg            *Config
	db             *gorm.DB
	extCrypto      bool
	extTimescaledb bool
	version        string
	serverVersion  string
}

var (
	log = common.MustGetLogger("orm")
)

const (
	minimalDbVersion = ">= 9.6.0"
)

// NewOrm ...
func NewOrm(cfg *Config,
	graceful *graceful_service.GracefulService) (orm *Orm, db *gorm.DB, err error) {

	fmt.Printf("database connect %s\n", cfg.String())
	db, err = gorm.Open("postgres", cfg.String())
	if err != nil {
		return
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

	if err = orm.check(); err != nil {
		return
	}

	graceful.Subscribe(orm)
	return
}

// Shutdown ...
func (o *Orm) Shutdown() {
	if o.db != nil {
		log.Debug("database shutdown")
		o.db.Close()
	}
}

func (o *Orm) check() (err error) {

	if err = o.checkServerVersion(); err != nil {
		return
	}

	err = o.checkExtensions()

	return
}

func (o *Orm) checkServerVersion() (err error) {

	// get version
	row := o.db.Raw("select version()").Row()
	if err = row.Scan(&o.version); err != nil {
		return
	}

	fmt.Println(o.version)

	// get server version
	row = o.db.Raw("SHOW server_version").Row()
	if err = row.Scan(&o.serverVersion); err != nil {
		return
	}

	// check server version
	var v *semver.Constraints
	if v, err = semver.NewConstraint(minimalDbVersion); err != nil {
		return
	}

	var client *semver.Version
	if client, err = semver.NewVersion(o.serverVersion); err != nil {
		return
	}

	if ok := v.Check(client); !ok {
		err = fmt.Errorf("unsupported database version %s, expected %s", o.serverVersion, minimalDbVersion)
		return
	}
	return
}

func (o *Orm) checkAvailableExtensions(availableExtensions []Extension, extName string) (exist bool) {
	for _, ext := range availableExtensions {
		if ext.Extname == extName {
			exist = true
			return
		}
	}
	return
}

func (o *Orm) checkExtensions() (err error) {

	// check extensions
	availableExtensions := make([]Extension, 0)
	if err = o.db.Raw("select * from pg_available_extensions").Scan(&availableExtensions).Error; err != nil {
		return
	}

	// check extensions
	extensions := make([]Extension, 0)
	if err = o.db.Raw("select * from pg_extension").Scan(&extensions).Error; err != nil {
		return
	}

	for _, ext := range extensions {
		switch ext.Extname {
		case "pgcrypto":
			o.extCrypto = true
		case "timescaledb":
			o.extTimescaledb = true
		default:

		}
	}

	if !o.extCrypto {
		if o.checkAvailableExtensions(availableExtensions, "pgcrypto") {
			err = fmt.Errorf("extension 'pgcrypto' installed but not enabled, enable it \nCREATE EXTENSION IF NOT EXISTS pgcrypto CASCADE;")
			return
		}
		err = fmt.Errorf("please install pgcrypto extension for postgresql database (maybe need install postgresql-contrib)")
	}

	if !o.extTimescaledb {
		if o.checkAvailableExtensions(availableExtensions, "timescaledb") {
			err = fmt.Errorf("extension 'timescaledb' installed but not enabled, enable it \nCREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;")
			return
		}
		fmt.Println("please install timescaledb extension, website: https://docs.timescale.com/v1.1/getting-started/installation)")
	}

	return
}
