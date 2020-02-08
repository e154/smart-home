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

package logging

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/config"
	"github.com/op/go-logging"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type LogBackend struct {
	dbSaver *LogDbSaver
	oldLog  *m.Log
	logging bool
	mx      *sync.Mutex
	L       *logrus.Logger
}

func NewLogBackend(logger *logrus.Logger, dbSaver *LogDbSaver, conf *config.AppConfig) (back *LogBackend) {
	back = &LogBackend{
		dbSaver: dbSaver,
		logging: conf.Logging,
		mx:      &sync.Mutex{},
		L:       logger,
	}
	back.L.Out = os.Stdout
	return
}

func (b *LogBackend) Log(level logging.Level, calldepth int, rec *logging.Record) (err error) {

	if !b.logging {
		return
	}

	b.mx.Lock()
	defer b.mx.Unlock()

	var logLevel common.LogLevel
	s := rec.Formatted(calldepth + 1)
	switch level {
	case logging.CRITICAL:
		b.L.Level = logrus.FatalLevel
		b.L.Fatal(s)
		logLevel = "Critical"
	case logging.ERROR:
		b.L.Level = logrus.ErrorLevel
		b.L.Error(s)
		logLevel = "Error"
	case logging.WARNING:
		b.L.Level = logrus.WarnLevel
		b.L.Warning(s)
		logLevel = "Warning"
	case logging.INFO, logging.NOTICE:
		b.L.Level = logrus.InfoLevel
		b.L.Info(s)
		logLevel = "Info"
	case logging.DEBUG:
		b.L.Level = logrus.DebugLevel
		b.L.Debug(s)
		logLevel = "Debug"
	}

	if b.oldLog != nil {
		if b.oldLog.Body == rec.Message() && b.oldLog.Level == logLevel {
			return
		}
	}

	record := &m.Log{
		Level:     logLevel,
		Body:      rec.Message(),
		CreatedAt: rec.Time,
	}

	b.oldLog = record

	b.dbSaver.Save(*record)

	return nil
}
