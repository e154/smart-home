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
	"github.com/e154/smart-home/system/config"
)

// Config ...
type Config struct {
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

// String ...
func (c Config) String() string {

	// parseTime https://github.com/go-sql-driver/mysql#parsetime
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable", c.Name, c.User, c.Password, c.Host, c.Port)
}

// NewConfig ...
func NewConfig(cfg *config.AppConfig) *Config {
	return &Config{
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
