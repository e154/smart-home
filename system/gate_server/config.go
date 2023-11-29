// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package gate_server

import (
	"strconv"
	"time"
)

// Config configures an Server
type Config struct {
	Host        string
	Port        int
	Timeout     int
	IdleTimeout int
	SecretKey   string
}

// GetAddr returns the address to specify a HTTP server address
func (c Config) GetAddr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

// GetTimeout returns the time.Duration converted to millisecond
func (c Config) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Millisecond
}
