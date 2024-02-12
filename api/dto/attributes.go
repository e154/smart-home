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

package dto

import (
	"fmt"
	stub "github.com/e154/smart-home/api/stub"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/encryptor"
	m "github.com/e154/smart-home/models"
	"strconv"
	"strings"
)

// AttributeFromApi ...
func AttributeFromApi(apiAttr map[string]stub.ApiAttribute) (attributes m.Attributes) {
	if apiAttr == nil {
		return
	}
	return attributeFromApi(apiAttr)
}

func attributeFromApi(apiAttr map[string]stub.ApiAttribute) (attributes m.Attributes) {
	attributes = make(m.Attributes)
	for k, v := range apiAttr {
		attr := &m.Attribute{
			Name: v.Name,
		}
		switch v.Type {
		case stub.INT:
			if v.Int != nil {
				attr.Value = *v.Int
			}
			attr.Type = common.AttributeInt
		case stub.STRING:
			if v.String != nil {
				attr.Value = *v.String
			}
			attr.Type = common.AttributeString
		case stub.BOOL:
			if v.Bool != nil {
				attr.Value = *v.Bool
			}
			attr.Type = common.AttributeBool
		case stub.FLOAT:
			if v.Float != nil {
				attr.Value = *v.Float
			}
			attr.Type = common.AttributeFloat
		case stub.IMAGE:
			if v.ImageUrl != nil {
				attr.Value = *v.ImageUrl
			}
			attr.Type = common.AttributeImage
		case stub.ICON:
			if v.Icon != nil {
				attr.Value = *v.Icon
			}
			attr.Type = common.AttributeIcon
		case stub.ARRAY:
			//	attr.Value = v.Array
			attr.Type = common.AttributeArray
		case stub.MAP:
			attr.Type = common.AttributeMap
			//attr.Value = AttributeFromApi(v.Map)
		case stub.TIME:
			if v.Time != nil {
				attr.Value = *v.Time
			}
			attr.Type = common.AttributeTime
		case stub.POINT:
			if v.Point != nil {
				point := []interface{}{0.0, 0.0}
				str := *v.Point
				str = strings.ReplaceAll(str, "[", "")
				str = strings.ReplaceAll(str, "]", "")
				str = strings.ReplaceAll(str, " ", "")
				arr := strings.Split(str, ",")
				if len(arr) == 2 {
					point[0], _ = strconv.ParseFloat(arr[0], 64)
					point[1], _ = strconv.ParseFloat(arr[1], 64)
				}
				attr.Value = point
			}
			attr.Type = common.AttributePoint
		case stub.ENCRYPTED:
			if v.Encrypted != nil {
				value, err := encryptor.Encrypt(*v.Encrypted)
				if err == nil {
					attr.Value = value
				}
			}
			attr.Type = common.AttributeEncrypted
		}
		attributes[k] = attr
	}
	return
}

// AttributeToApi ...
func AttributeToApi(attributes m.Attributes) (apiAttr map[string]stub.ApiAttribute) {
	apiAttr = make(map[string]stub.ApiAttribute)
	var attr stub.ApiAttribute
	for k, v := range attributes {
		attr = stub.ApiAttribute{
			Name: v.Name,
		}
		switch v.Type {
		case "int":
			attr.Type = stub.INT
			attr.Int = common.Int64(v.Int64())
		case "string":
			attr.Type = stub.STRING
			attr.String = common.String(v.String())
		case "bool":
			attr.Type = stub.BOOL
			attr.Bool = common.Bool(v.Bool())
		case "float":
			attr.Type = stub.FLOAT
			attr.Float = common.Float32(float32(v.Float64()))
		case "array":
			attr.Type = stub.ARRAY
		case "map":
			attr.Type = stub.MAP
		case "time":
			attr.Type = stub.TIME
			attr.Time = common.Time(v.Time())
		case "image":
			attr.Type = stub.IMAGE
			attr.ImageUrl = common.String(v.String())
		case "icon":
			attr.Type = stub.ICON
			attr.Icon = common.String(v.String())
		case "point":
			attr.Type = stub.POINT
			attr.Point = common.String(fmt.Sprintf("[%f, %f]", v.Point().Lon, v.Point().Lat))
		case "encrypted":
			attr.Type = stub.ENCRYPTED
			attr.Encrypted = common.String(v.Decrypt())
		}
		apiAttr[k] = attr
	}
	return
}
