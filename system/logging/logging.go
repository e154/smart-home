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
	"github.com/e154/smart-home/system/logging_db"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logging struct {
	logger     *zap.Logger
	logDbSaver *logging_db.LogDbSaver
	oldLog     *m.Log
}

func NewLogger(logDbSaver *logging_db.LogDbSaver) (logging *Logging) {

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = nil
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	logging = &Logging{
		logDbSaver: logDbSaver,
	}

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core, zap.AddCaller(), zap.Hooks(logging.selfSaver))

	zap.ReplaceGlobals(logger)

	logging.logger = logger

	return
}

func (b *Logging) selfSaver(e zapcore.Entry) (err error) {

	var logLevel common.LogLevel
	switch e.Level {
	case zapcore.ErrorLevel:
		logLevel = "Error"
	case zapcore.WarnLevel:
		logLevel = "Warning"
	case zapcore.InfoLevel:
		logLevel = "Info"
	case zapcore.DebugLevel:
		logLevel = "Debug"
	}

	if b.oldLog != nil {
		if b.oldLog.Body == e.Message && b.oldLog.Level == logLevel {
			return
		}
	}

	record := &m.Log{
		Level:     logLevel,
		Body:      e.Message,
		CreatedAt: e.Time,
	}

	b.oldLog = record

	b.logDbSaver.Save(*record)

	return nil
}

func (b *Logging) Infof(log string, fields ...zapcore.Field) {
	zap.L().Info(log, fields...)
}

func (b *Logging) Info(log string, fields ...zapcore.Field) {
	zap.L().Info(log, fields...)
}

func (b *Logging) Debug(log string, fields ...zapcore.Field) {
	zap.L().Debug(log, fields...)
}

func (b *Logging) Debugf(log string, fields ...zapcore.Field) {
	zap.L().Debug(log, fields...)
}

func (b *Logging) Warn(log string, fields ...zapcore.Field) {
	zap.L().Warn(log, fields...)
}

func (b *Logging) Warnf(log string, fields ...zapcore.Field) {
	zap.L().Warn(log, fields...)
}

func (b *Logging) Error(log string, fields ...zapcore.Field) {
	zap.L().Error(log, fields...)
}

func (b *Logging) Errorf(log string, fields ...zapcore.Field) {
	zap.L().Error(log, fields...)
}

func (b *Logging) Panic(log string, fields ...zapcore.Field) {
	zap.L().Panic(log, fields...)
}

func Info(log string, fields ...zapcore.Field) {
	zap.L().Info(log, fields...)
}

func Debug(log string, fields ...zapcore.Field) {
	zap.L().Debug(log, fields...)
}

func Warn(log string, fields ...zapcore.Field) {
	zap.L().Warn(log, fields...)
}

func Error(log string, fields ...zapcore.Field) {
	zap.L().Error(log, fields...)
}

func Panic(log string, fields ...zapcore.Field) {
	zap.L().Panic(log, fields...)
}


func Infof(log string, fields ...zapcore.Field) {
	zap.L().Info(log, fields...)
}

func Debugf(log string, fields ...zapcore.Field) {
	zap.L().Debug(log, fields...)
}

func Warnf(log string, fields ...zapcore.Field) {
	zap.L().Warn(log, fields...)
}

func Errorf(log string, fields ...zapcore.Field) {
	zap.L().Error(log, fields...)
}

func Panicf(log string, fields ...zapcore.Field) {
	zap.L().Panic(log, fields...)
}

func MustGetLogger(p string) *zap.Logger {
	return zap.L()
}
