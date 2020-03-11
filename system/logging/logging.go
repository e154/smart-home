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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"sync"
)

type Logging struct {
	logger     *zap.Logger
	logDbSaver *LogDbSaver
	oldLogLock *sync.Mutex
	oldLog     m.Log
}

func NewLogger(logDbSaver *LogDbSaver) (logging *Logging) {

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
	config.EncodeLevel = CustomLevelEncoder
	config.EncodeName = CustomNameEncoder
	config.EncodeCaller = CustomCallerEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	logging = &Logging{
		logDbSaver: logDbSaver,
		oldLogLock: &sync.Mutex{},
	}

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Hooks(logging.selfSaver))


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

	b.oldLogLock.Lock()
	defer b.oldLogLock.Unlock()

	if b.oldLog.Body == e.Message && b.oldLog.Level == logLevel {
		return
	}

	record := m.Log{
		Level:     logLevel,
		Body:      e.Message,
		CreatedAt: e.Time,
	}

	b.oldLog = record

	b.logDbSaver.Save(record)

	return nil
}

func CustomLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := _levelToCapitalColorString[l]
	if !ok {
		s = _unknownLevelColor.Add(l.CapitalString())
	}
	enc.AppendString(s)
}

//TODO fix
func CustomNameEncoder(v string, enc zapcore.PrimitiveArrayEncoder) {
	var builder strings.Builder
	builder.WriteString(White.Add(v))
	builder.WriteString("                                      ")
	enc.AppendString(builder.String()[0:25])
}

func CustomCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(caller.TrimmedPath() + " >")
}
