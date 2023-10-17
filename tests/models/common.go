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

package models

import (
	"context"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// NetAttr ...
func NetAttr() m.Attributes {
	return m.Attributes{
		"s": {
			Name: "s",
			Type: common.AttributeString,
		},
		"d": {
			Name: "d",
			Type: common.AttributeEncrypted,
		},
		"p": {
			Name: "p",
			Type: common.AttributePoint,
		},
		"i": {
			Name: "i",
			Type: common.AttributeInt,
		},
		"f": {
			Name: "f",
			Type: common.AttributeFloat,
		},
		"b": {
			Name: "b",
			Type: common.AttributeBool,
		},
		"m": {
			Name: "m",
			Type: common.AttributeMap,
			Value: m.Attributes{
				"s2": {
					Name: "s2",
					Type: common.AttributeString,
				},
				"i2": {
					Name: "i2",
					Type: common.AttributeInt,
				},
				"f2": {
					Name: "f2",
					Type: common.AttributeFloat,
				},
				"b2": {
					Name: "b2",
					Type: common.AttributeBool,
				},
				"m2": {
					Name: "m2",
					Type: common.AttributeMap,
					Value: m.Attributes{
						"s3": {
							Name: "s3",
							Type: common.AttributeString,
						},
						"i3": {
							Name: "i3",
							Type: common.AttributeInt,
						},
						"f3": {
							Name: "f3",
							Type: common.AttributeFloat,
						},
						"b3": {
							Name: "b3",
							Type: common.AttributeBool,
						},
					},
				},
			},
		},
	}

}

func NetSettings() m.Attributes {
	return m.Attributes{
		"s": {
			Name:  "s",
			Type:  common.AttributeString,
			Value: "s",
		},
	}

}

// AddPlugin ...
func AddPlugin(adaptors *adaptors.Adaptors, name string, opts ...m.AttributeValue) (err error) {
	plugin := &m.Plugin{
		Name:    name,
		Version: "0.0.1",
		Enabled: true,
		System:  true,
	}
	if len(opts) > 0 {
		plugin.Settings = opts[0]
	}
	err = adaptors.Plugin.CreateOrUpdate(context.Background(), plugin)
	return
}
