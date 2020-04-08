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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"time"
)

// Orm ...
type Orm struct {
	cfg *OrmConfig
	db  *gorm.DB
}

var (
	log = common.MustGetLogger("orm")
)

// NewOrm ...
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

// Shutdown ...
func (o *Orm) Shutdown() {
	if o.db != nil {
		log.Debug("database shutdown")
		o.db.Close()
	}
}
