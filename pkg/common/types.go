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

package common

import "github.com/gin-gonic/gin"

// ConditionType ...
type ConditionType string

const (
	// ConditionOr ...
	ConditionOr = ConditionType("or")
	// ConditionAnd ...
	ConditionAnd = ConditionType("and")
)

// RunMode ...
type RunMode string

const (
	// DebugMode ...
	DebugMode = RunMode("debug")
	// ReleaseMode ...
	ReleaseMode = RunMode("release")
	// DemoMode
	DemoMode = RunMode("demo")
)

// MetricType ...
type MetricType string

const (
	// MetricTypeLine ...
	MetricTypeLine = MetricType("line")
	// MetricTypeBar ...
	MetricTypeBar = MetricType("bar")
	// MetricTypeDoughnut ...
	MetricTypeDoughnut = MetricType("doughnut")
	// MetricTypeRadar ...
	MetricTypeRadar = MetricType("radar")
	// MetricTypePie ...
	MetricTypePie = MetricType("pie")
	// MetricTypeHorizontalBar ...
	MetricTypeHorizontalBar = MetricType("horizontal bar")
)

type MetricRange string

const (
	MetricRange6H  = MetricRange("6h")
	MetricRange12H = MetricRange("12h")
	MetricRange24H = MetricRange("24h")
	MetricRange7d  = MetricRange("7d")
	MetricRange30d = MetricRange("30d")
	MetricRange1m  = MetricRange("1m")
)

func (m MetricRange) String() string {
	return string(m)
}

func (m MetricRange) Ptr() *MetricRange {
	return &m
}

// LogLevel ...
type LogLevel string

const (
	// LogLevelEmergency ...
	LogLevelEmergency = LogLevel("Emergency")
	// LogLevelAlert ...
	LogLevelAlert = LogLevel("Alert")
	// LogLevelCritical ...
	LogLevelCritical = LogLevel("Critical")
	// LogLevelError ...
	LogLevelError = LogLevel("Error")
	// LogLevelWarning ...
	LogLevelWarning = LogLevel("Warning")
	// LogLevelNotice ...
	LogLevelNotice = LogLevel("Notice")
	// LogLevelInfo ...
	LogLevelInfo = LogLevel("Info")
	// LogLevelDebug ...
	LogLevelDebug = LogLevel("Debug")
)

// StatusType ...
type StatusType string

const (
	// Enabled ...
	Enabled = StatusType("enabled")
	// Disabled ...
	Disabled = StatusType("disabled")
	// Frozen ...
	Frozen = StatusType("frozen")
)

// ScriptLang ...
type ScriptLang string

const (
	// ScriptLangTs ...
	ScriptLangTs = ScriptLang("ts")
	// ScriptLangCoffee ...
	ScriptLangCoffee = ScriptLang("coffeescript")
	// ScriptLangJavascript ...
	ScriptLangJavascript = ScriptLang("javascript")
)

// GinEngine ...
type GinEngine interface {
	GetEngine() *gin.Engine
}
