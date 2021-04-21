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
	"math"
	"time"
)

// String ...
func String(v string) *string {
	return &v
}

// StringValue ...
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// Int ...
func Int(v int) *int {
	return &v
}

// IntValue ...
func IntValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

// Int64 ...
func Int64(v int64) *int64 {
	return &v
}

// Int64Value ...
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// Time ...
func Time(v time.Time) *time.Time {
	return &v
}

// TimeValue ...
func TimeValue(v *time.Time) time.Time {
	if v != nil {
		return *v
	}
	return time.Time{}
}

// ToEntityPrototypeType ...
func ToEntityPrototypeType(v EntityPrototypeType) *EntityPrototypeType {
	return &v
}

// ToEntityPrototypeTypeValue ...
func ToEntityPrototypeTypeValue(v *EntityPrototypeType) EntityPrototypeType {
	if v != nil {
		return *v
	}
	return ""
}

// Rounding ...
func Rounding(num float64, k uint) float64 {
	p := math.Pow(10, float64(k))
	return math.Floor(num*p) / p
}
