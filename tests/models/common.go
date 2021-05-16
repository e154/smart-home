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

package models

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

func NetEntityAttr() m.EntityAttributes {
	return m.EntityAttributes{
		"s": {
			Name: "s",
			Type: common.EntityAttributeString,
		},
		"i": {
			Name: "i",
			Type: common.EntityAttributeInt,
		},
		"f": {
			Name: "f",
			Type: common.EntityAttributeFloat,
		},
		"b": {
			Name: "b",
			Type: common.EntityAttributeBool,
		},
		"m": {
			Name: "m",
			Type: common.EntityAttributeMap,
			Value: m.EntityAttributes{
				"s2": {
					Name: "s2",
					Type: common.EntityAttributeString,
				},
				"i2": {
					Name: "i2",
					Type: common.EntityAttributeInt,
				},
				"f2": {
					Name: "f2",
					Type: common.EntityAttributeFloat,
				},
				"b2": {
					Name: "b2",
					Type: common.EntityAttributeBool,
				},
				"m2": {
					Name: "m2",
					Type: common.EntityAttributeMap,
					Value: m.EntityAttributes{
						"s3": {
							Name: "s3",
							Type: common.EntityAttributeString,
						},
						"i3": {
							Name: "i3",
							Type: common.EntityAttributeInt,
						},
						"f3": {
							Name: "f3",
							Type: common.EntityAttributeFloat,
						},
						"b3": {
							Name: "b3",
							Type: common.EntityAttributeBool,
						},
					},
				},
			},
		},
	}

}
