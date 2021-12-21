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
	"encoding/json"
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
		_ = json.Unmarshal(v.Value.Value, attr.Value)
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
			Type: string(v.Type),
		}
	}
	return
}
