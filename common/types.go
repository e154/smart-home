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

import "github.com/gin-gonic/gin"

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

// FlowElementsPrototypeType ...
type FlowElementsPrototypeType string

const (
	// FlowElementsPrototypeDefault ...
	FlowElementsPrototypeDefault = FlowElementsPrototypeType("default")
	// FlowElementsPrototypeMessageHandler ...
	FlowElementsPrototypeMessageHandler = FlowElementsPrototypeType("MessageHandler")
	// FlowElementsPrototypeMessageEmitter ...
	FlowElementsPrototypeMessageEmitter = FlowElementsPrototypeType("MessageEmitter")
	// FlowElementsPrototypeTask ...
	FlowElementsPrototypeTask = FlowElementsPrototypeType("Task")
	// FlowElementsPrototypeGateway ...
	FlowElementsPrototypeGateway = FlowElementsPrototypeType("Gateway")
	// FlowElementsPrototypeFlow ...
	FlowElementsPrototypeFlow = FlowElementsPrototypeType("Flow")
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

// DeviceType ...
type DeviceType string

// PrototypeType ...
type PrototypeType string

const (
	// PrototypeTypeText ...
	PrototypeTypeText = PrototypeType("text")
	// PrototypeTypeImage ...
	PrototypeTypeImage = PrototypeType("image")
	// PrototypeTypeDevice ...
	PrototypeTypeDevice = PrototypeType("device")
	// PrototypeTypeEmpty ...
	PrototypeTypeEmpty = PrototypeType("")
)

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

// GinEngine ...
type GinEngine interface {
	GetEngine() *gin.Engine
}

// MapDeviceHistoryType ...
type MapDeviceHistoryType string

const (
	// MapDeviceHistoryState ...
	MapDeviceHistoryState = MapDeviceHistoryType("state")
	// MapDeviceHistoryOption ...
	MapDeviceHistoryOption = MapDeviceHistoryType("option")
)

type MetricType string

const (
	MetricTypeLine          = MetricType("line")
	MetricTypeBar           = MetricType("bar")
	MetricTypeDoughnut      = MetricType("doughnut")
	MetricTypeRadar         = MetricType("radar")
	MetricTypePie           = MetricType("pie")
	MetricTypeHorizontalBar = MetricType("horizontal bar")
)
