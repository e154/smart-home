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

package orm

import (
	"context"
	"database/sql"
	"fmt"
	goLog "log"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/e154/smart-home/common/logger"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Orm ...
type Orm struct {
	cfg                 *Config
	db                  *gorm.DB
	extTimescaledb      bool
	availableExtensions []AvailableExtension
	version             string
	serverVersion       string
}

var (
	log = logger.MustGetLogger("orm")
)

const (
	minimalDbVersion = ">= 9.6.0"
)

// NewOrm ...
func NewOrm(lc fx.Lifecycle,
	cfg *Config) (orm *Orm, db *gorm.DB, err error) {

	orm = &Orm{
		cfg:                 cfg,
		availableExtensions: make([]AvailableExtension, 0),
	}

	if err = orm.Start(); err != nil {
		log.Error(err.Error())
		return
	}

	db = orm.db

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return orm.Shutdown()
		},
	})

	return
}

// Start ...
func (o *Orm) Start() (err error) {

	log.Infof("database connect %s", strings.ReplaceAll(o.cfg.String(), "password="+o.cfg.Password, "password=*****"))

	var logLevel = gormLogger.Silent
	if o.cfg.Debug {
		logLevel = gormLogger.Info
	}

	newLogger := gormLogger.New(
		goLog.New(os.Stdout, "\r\n", goLog.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	o.db, err = gorm.Open(postgres.Open(o.cfg.String()), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 newLogger,
	})
	if err != nil {
		// it for DI
		err = nil
		return
	}

	var db *sql.DB
	if db, err = o.db.DB(); err != nil {
		return
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	db.SetMaxIdleConns(o.cfg.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(o.cfg.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Duration(o.cfg.ConnMaxLifeTime) * time.Minute)

	// get server version
	row := o.db.Raw("SHOW server_version").Row()
	if err = row.Scan(&o.serverVersion); err != nil {
		return
	}

	log.Infof("database version %s", o.serverVersion)

	err = o.Check()

	return
}

// DB ...
func (o *Orm) DB() *sql.DB {
	db, _ := o.db.DB()
	return db
}

// Shutdown ...
func (o *Orm) Shutdown() (err error) {
	if o.db != nil {
		log.Info("database shutdown")
		var db *sql.DB
		if db, err = o.db.DB(); err != nil {
			return
		}
		err = db.Close()
	}
	return
}

func (o *Orm) Check() (err error) {

	if err = o.checkServerVersion(); err != nil {
		return
	}

	err = o.CheckExtensions()

	return
}

func (o *Orm) checkServerVersion() (err error) {

	// get version
	row := o.db.Raw("select version()").Row()
	if err = row.Scan(&o.version); err != nil {
		return
	}

	strArr := strings.Split(o.serverVersion, " ")
	if len(strArr) > 1 {
		o.serverVersion = strArr[0]
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

func (o *Orm) CheckAvailableExtension(extName string) (exist bool) {
	for _, ext := range o.availableExtensions {
		if ext.Name == extName {
			exist = true
			return
		}
	}
	return
}

func (o *Orm) CheckInstalledExtension(extName string) (exist bool) {
	for _, ext := range o.availableExtensions {
		if ext.Name == extName && ext.InstalledVersion != nil {
			exist = true
			return
		}
	}
	return
}

func (o *Orm) CheckExtensions() (err error) {

	// check extensions
	if err = o.db.Raw("select * from pg_available_extensions").Scan(&o.availableExtensions).Error; err != nil {
		return
	}

	if !o.CheckAvailableExtension("pgcrypto") {
		log.Warn("please install pgcrypto extension for postgresql database (maybe need install postgresql-contrib)\r")
	} else {
		if !o.CheckInstalledExtension("pgcrypto") {
			o.db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto CASCADE;`)
		}
	}

	if !o.CheckAvailableExtension("timescaledb") {
		log.Warn("please install timescaledb extension, website: https://docs.timescale.com/v1.1/getting-started/installation)\r")
	} else {
		if !o.CheckInstalledExtension("timescaledb") {
			o.db.Exec(`CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;`)
		}
		o.db.Exec(`SELECT create_hypertable('metric_bucket', 'time', migrate_data => true, if_not_exists => TRUE);`)
	}

	return
}

// ExtTimescaledbEnabled ...
func (o Orm) ExtTimescaledbEnabled() bool {
	return o.extTimescaledb
}

func (o *Orm) Ping() (latency float64, err error) {

	start := time.Now()

	if err = o.DB().Ping(); err != nil {
		return
	}

	diff := time.Since(start).Microseconds()
	latency = float64(diff) / 1000000.0

	return
}
