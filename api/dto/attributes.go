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

package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// AttributeFromApi ...
func AttributeFromApi(apiAttr map[string]*api.Attribute) (attributes m.Attributes) {
	attributes = make(m.Attributes)
	for k, v := range apiAttr {
		attr := &m.Attribute{
			Name: v.Name,
			Type: common.AttributeType(v.Type),
		}
		switch v.Type {
		case api.Types_INT:
			attr.Value = v.GetInt()
		case api.Types_STRING:
			attr.Value = v.String()
		case api.Types_BOOL:
			attr.Value = v.GetBool()
		case api.Types_FLOAT:
			attr.Value = v.GetFloat()
		case api.Types_ARRAY:
		//	attr.Value = v.GetArray()
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
		}
	}
	return
}
