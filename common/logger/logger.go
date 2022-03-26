// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package logger

import (
	"go.uber.org/zap"
)

// Logger ...
type Logger struct {
	p string
}

// MustGetLogger ...
func MustGetLogger(p string) *Logger {
	return &Logger{
		p: p,
	}
}

// Error ...
func (l *Logger) Error(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Error(format)
}

// Errorf ...
func (l *Logger) Errorf(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Errorf(format, args...)
}

// Warn ...
func (l *Logger) Warn(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Warn(format)
}

// Warnf ...
func (l *Logger) Warnf(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Warnf(format, args...)
}

// Info ...
func (l *Logger) Info(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Info(format)
}

// Infof ...
func (l *Logger) Infof(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Infof(format, args...)
}

// Debug ...
func (l *Logger) Debug(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Debug(format)
}

// Debugf ...
func (l *Logger) Debugf(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Debugf(format, args...)
}

// Fatal ...
func (l *Logger) Fatal(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Fatal(format)
}

// Fatalf ...
func (l *Logger) Fatalf(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Fatalf(format, args...)
}
