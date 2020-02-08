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

type ScriptLang string

const (
	ScriptLangTs         = ScriptLang("ts")
	ScriptLangCoffee     = ScriptLang("coffeescript")
	ScriptLangJavascript = ScriptLang("javascript")
)

type FlowElementsPrototypeType string

const (
	FlowElementsPrototypeDefault        = FlowElementsPrototypeType("default")
	FlowElementsPrototypeMessageHandler = FlowElementsPrototypeType("MessageHandler")
	FlowElementsPrototypeMessageEmitter = FlowElementsPrototypeType("MessageEmitter")
	FlowElementsPrototypeTask           = FlowElementsPrototypeType("Task")
	FlowElementsPrototypeGateway        = FlowElementsPrototypeType("Gateway")
	FlowElementsPrototypeFlow           = FlowElementsPrototypeType("Flow")
)

type StatusType string

const (
	Enabled  = StatusType("enabled")
	Disabled = StatusType("disabled")
	Frozen   = StatusType("frozen")
)

type DeviceType string

type PrototypeType string

const (
	PrototypeTypeText   = PrototypeType("text")
	PrototypeTypeImage  = PrototypeType("image")
	PrototypeTypeDevice = PrototypeType("device")
	PrototypeTypeEmpty  = PrototypeType("")
)

type LogLevel string

const (
	LogLevelEmergency = LogLevel("Emergency")
	LogLevelAlert     = LogLevel("Alert")
	LogLevelCritical  = LogLevel("Critical")
	LogLevelError     = LogLevel("Error")
	LogLevelWarning   = LogLevel("Warning")
	LogLevelNotice    = LogLevel("Notice")
	LogLevelInfo      = LogLevel("Info")
	LogLevelDebug     = LogLevel("Debug")
)

type GinEngine interface {
	GetEngine() *gin.Engine
}
