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

package version

import (
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

const (
	// EntityVersion ...
	EntityVersion = string("version")

	AttrVersion     = string("version")
	AttrRevision    = string("revision")
	AttrRevisionURL = string("revision_url")
	AttrGenerated   = string("generated")
	AttrDevelopers  = string("developers")
	AttrBuildNum    = string("build_num")
	AttrDockerImage = string("docker_image")
	AttrGoVersion   = string("go_version")

	// Name ...
	Name = "version"

	// EntityType ...
	EntityType = "version"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrVersion: {
			Name: AttrVersion,
			Type: common.AttributeString,
		},
		AttrRevision: {
			Name: AttrRevision,
			Type: common.AttributeString,
		},
		AttrRevisionURL: {
			Name: AttrRevisionURL,
			Type: common.AttributeString,
		},
		AttrGenerated: {
			Name: AttrGenerated,
			Type: common.AttributeString,
		},
		AttrDevelopers: {
			Name: AttrDevelopers,
			Type: common.AttributeString,
		},
		AttrBuildNum: {
			Name: AttrBuildNum,
			Type: common.AttributeString,
		},
		AttrDockerImage: {
			Name: AttrDockerImage,
			Type: common.AttributeString,
		},
		AttrGoVersion: {
			Name: AttrGoVersion,
			Type: common.AttributeString,
		},
	}
}
