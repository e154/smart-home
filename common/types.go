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

package common

import (
	"github.com/gin-gonic/gin"
	"strings"
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

// MapElementPrototypeId ...
type MapElementPrototypeId interface{}

// MapElementPrototypeType ...
type MapElementPrototypeType string

const (
	// MapElementPrototypeText ...
	MapElementPrototypeText = MapElementPrototypeType("text")
	// MapElementPrototypeImage ...
	MapElementPrototypeImage = MapElementPrototypeType("image")
	// MapElementPrototypeEntity ...
	MapElementPrototypeEntity = MapElementPrototypeType("entity")
	// MapElementPrototypeEmpty ...
	MapElementPrototypeEmpty = MapElementPrototypeType("")
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

// EntityHistoryType ...
type EntityHistoryType string

const (
	// EntityHistoryState ...
	EntityHistoryState = EntityHistoryType("state")
	// EntityHistoryOption ...
	EntityHistoryOption = EntityHistoryType("option")
)

// MetricType
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

// Icon ...
type Icon string

// NewIcon ...
func NewIcon(v string) *Icon {
	s := Icon(v)
	return &s
}

// String ...
func (i *Icon) String() string {
	if i == nil {
		return ""
	}
	return string(*i)
}

// EntityType ...
type EntityType string

// EntityId ...
type EntityId string

func NewEntityId(s string) *EntityId {
	e := EntityId(s)
	return &e
}

func (e EntityId) Name() string {
	arr := strings.Split(string(e), ".")
	if len(arr) > 1 {
		return arr[1]
	}
	return string(e)
}

func (e EntityId) Type() EntityType {
	arr := strings.Split(string(e), ".")
	if len(arr) > 1 {
		return EntityType(arr[0])
	}
	return EntityType(e)
}

func (e *EntityId) String() string {
	if e == nil {
		return ""
	} else {
		return string(*e)
	}
}

func (e EntityType) String() string {
	return string(e)
}

type AttributeType string

const (
	AttributeString = AttributeType("string")
	AttributeInt    = AttributeType("int")
	AttributeTime   = AttributeType("time")
	AttributeBool   = AttributeType("bool")
	AttributeFloat  = AttributeType("float")
	//DEPRECATED
	AttributeArray = AttributeType("array")
	AttributeMap   = AttributeType("map")
)

type ConditionType string

const (
	ConditionOr  = ConditionType("or")
	ConditionAnd = ConditionType("and")
)
