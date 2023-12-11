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

package models

import (
	"fmt"
)

// GateConfig ...
type GateConfig struct {
	ApiHttpPort      int    `json:"api_http_port" env:"API_HTTP_PORT"`
	ApiDebug         bool   `json:"api_debug" env:"API_DEBUG"`
	ApiGzip          bool   `json:"api_gzip" env:"API_GZIP"`
	Domain           string `json:"domain" env:"DOMAIN"`
	Pprof            bool   `json:"pprof" env:"PPROF"`
	Https            bool   `json:"https" env:"HTTPS"`
	ProxyTimeout     int    `json:"proxy_timeout" env:"PROXY_TIMEOUT"`
	ProxyIdleTimeout int    `json:"proxy_idle_timeout" env:"PROXY_IDLE_TIMEOUT"`
	ProxySecretKey   string `json:"proxy_secret_key" env:"PROXY_SECRET_KEY"`
}

func (c *GateConfig) ApiScheme() (scheme string) {
	scheme = "http"
	if c.Https {
		scheme = "https"
	}
	return
}

func (c *GateConfig) ApiFullAddress() (scheme string) {
	return fmt.Sprintf("%s://%s:%d", c.ApiScheme(), c.Domain, c.ApiHttpPort)
}
