// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package console

import (
	"github.com/e154/smart-home/common/logger"
)

var (
	log = logger.MustGetLogger("js")

	defaultStdPrinter Printer = &LogPrinter{}
)

// LogPrinter implements the console.Printer interface
// that prints to the stdout or stderr.
type LogPrinter struct {
}

// Log prints s to the stdout.
func (p LogPrinter) Log(s ...interface{}) {
	log.Infof("%v", s...)
}

// Warn prints s to the stderr.
func (p LogPrinter) Warn(s ...interface{}) {
	log.Warnf("%v", s...)
}

// Error prints s to the stderr.
func (p LogPrinter) Error(s ...interface{}) {
	log.Errorf("%v", s...)
}

// Debug prints s to the stderr.
func (p LogPrinter) Debug(s ...interface{}) {
	log.Debug("%v", s...)
}
