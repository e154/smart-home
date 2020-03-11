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

package common

import (
	"go.uber.org/zap"
)

type Logger struct {
	p string
}

func MustGetLogger(p string) *Logger {
	return &Logger{
		p: p,
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Error(format)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Errorf(format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Warn(format)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Warnf(format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Info(format)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Infof(format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Debug(format)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Debugf(format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Fatal(format)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	zap.L().Named(l.p).Sugar().Fatalf(format, args...)
}
