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
	"strconv"
	"strings"

	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/encryptor"
	m "github.com/e154/smart-home/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AttributeFromApi ...
func AttributeFromApi(apiAttr map[string]*api.Attribute) (attributes m.Attributes) {
	return attributeFromApi(apiAttr)
}

func attributeFromApi(apiAttr map[string]*api.Attribute) (attributes m.Attributes) {
	attributes = make(m.Attributes)
	for k, v := range apiAttr {
		attr := &m.Attribute{
			Name: v.Name,
		}
		switch v.Type {
		case api.Types_INT:
			attr.Value = v.GetInt()
			attr.Type = common.AttributeInt
		case api.Types_STRING:
			attr.Value = v.GetString_()
			attr.Type = common.AttributeString
		case api.Types_BOOL:
			attr.Value = v.GetBool()
			attr.Type = common.AttributeBool
		case api.Types_FLOAT:
			attr.Value = v.GetFloat()
			attr.Type = common.AttributeFloat
		case api.Types_IMAGE:
			attr.Value = v.GetImageUrl()
			attr.Type = common.AttributeImage
		case api.Types_ARRAY:
			//	attr.Value = v.GetArray()
			attr.Type = common.AttributeArray
		case api.Types_MAP:
			attr.Type = common.AttributeMap
			attr.Value = AttributeFromApi(v.GetMap())
		case api.Types_TIME:
			attr.Value = v.GetTime().AsTime()
			attr.Type = common.AttributeTime
		case api.Types_POINT:
			point := []interface{}{0.0, 0.0}
			str := v.GetPoint()
			str = strings.ReplaceAll(str, "[", "")
			str = strings.ReplaceAll(str, "]", "")
			str = strings.ReplaceAll(str, " ", "")
			arr := strings.Split(str, ",")
			point[0], _ = strconv.ParseFloat(arr[0], 64)
			point[1], _ = strconv.ParseFloat(arr[1], 64)
			attr.Value = point
			attr.Type = common.AttributePoint
		case api.Types_ENCRYPTED:
			value, err := encryptor.Encrypt(v.GetEncrypted())
			if err == nil {
				attr.Value = value
			}
			attr.Type = common.AttributeEncrypted
		}
		attributes[k] = attr
	}
	return
}

// AttributeToApi ...
func AttributeToApi(attributes m.Attributes) (apiAttr map[string]*api.Attribute) {
	apiAttr = make(map[string]*api.Attribute)
	for k, v := range attributes {
		apiAttr[k] = &api.Attribute{
			Name: v.Name,
		}
		switch v.Type {
		case "int":
			apiAttr[k].Type = api.Types_INT
			apiAttr[k].Int = common.Int64(v.Int64())
		case "string":
			apiAttr[k].Type = api.Types_STRING
			apiAttr[k].String_ = common.String(v.String())
		case "bool":
			apiAttr[k].Type = api.Types_BOOL
			apiAttr[k].Bool = common.Bool(v.Bool())
		case "float":
			apiAttr[k].Type = api.Types_FLOAT
			apiAttr[k].Float = common.Float32(float32(v.Float64()))
		case "array":
			apiAttr[k].Type = api.Types_ARRAY
		case "map":
			apiAttr[k].Type = api.Types_MAP
		case "time":
			apiAttr[k].Type = api.Types_TIME
			apiAttr[k].Time = timestamppb.New(v.Time())
		case "image":
			apiAttr[k].Type = api.Types_IMAGE
			apiAttr[k].ImageUrl = common.String(v.String())
		case "point":
			apiAttr[k].Type = api.Types_POINT
			apiAttr[k].Point = fmt.Sprintf("[%f, %f]", v.Point().Lon, v.Point().Lat)
		case "encrypted":
			apiAttr[k].Type = api.Types_ENCRYPTED
			apiAttr[k].Encrypted = common.String(v.Decrypt())
		}
	}
	return
}
