// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"strings"

	"github.com/gin-gonic/gin"
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

// EntityId ...
type EntityId string

// NewEntityId ...
func NewEntityId(s string) *EntityId {
	e := EntityId(s)
	return &e
}

// NewEntityIdFromPtr ...
func NewEntityIdFromPtr(s *string) *EntityId {
	if s == nil {
		return nil
	}
	e := EntityId(*s)
	return &e
}

// Name ...
func (e EntityId) Name() string {
	arr := strings.Split(string(e), ".")
	if len(arr) > 1 {
		return arr[1]
	}
	return string(e)
}

// PluginName ...
func (e EntityId) PluginName() string {
	arr := strings.Split(string(e), ".")
	if len(arr) > 1 {
		return arr[0]
	}
	return string(e)
}

// String ...
func (e *EntityId) String() string {
	if e == nil {
		return ""
	} else {
		return string(*e)
	}
}

// StringPtr ...
func (e *EntityId) StringPtr() *string {
	if e == nil {
		return nil
	}
	r := e.String()
	return &r
}

// Ptr ...
func (e EntityId) Ptr() *EntityId {
	return &e
}

// AttributeType ...
type AttributeType string

const (
	// AttributeString ...
	AttributeString = AttributeType("string")
	// AttributeInt ...
	AttributeInt = AttributeType("int")
	// AttributeTime ...
	AttributeTime = AttributeType("time")
	// AttributeBool ...
	AttributeBool = AttributeType("bool")
	// AttributeFloat ...
	AttributeFloat = AttributeType("float")
	// AttributeImage ...
	AttributeImage = AttributeType("image")
	// AttributePoint ...
	AttributePoint = AttributeType("point")
	// AttributeEncrypted ...
	AttributeEncrypted = AttributeType("encrypted")
	//DEPRECATED
	AttributeArray = AttributeType("array")
	//DEPRECATED
	AttributeMap = AttributeType("map")
)

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

// PageParams ...
type PageParams struct {
	Limit   int64  `json:"limit" validate:"required,gte=1,lte=1000"`
	Offset  int64  `json:"offset" validate:"required,gte=0,lte=1000"`
	Order   string `json:"order" validate:"required,oneof=created_at"`
	SortBy  string `json:"sort_by" validate:"required,oneof=desc asc"`
	PageReq int64
	SortReq string
}

// SearchParams ...
type SearchParams struct {
	Query  string `json:"query" validate:"required,min=1,max;255"`
	Limit  int64  `json:"limit" validate:"required,gte=1,lte=1000"`
	Offset int64  `json:"offset" validate:"required,gte=0,lte=1000"`
}
