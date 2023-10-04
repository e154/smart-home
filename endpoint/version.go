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

package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/version"
)

// VersionEndpoint ...
type VersionEndpoint struct {
	*CommonEndpoint
}

// NewVersionEndpoint ...
func NewVersionEndpoint(common *CommonEndpoint) *VersionEndpoint {
	return &VersionEndpoint{
		CommonEndpoint: common,
	}
}

// ServerVersion ...
func (v *VersionEndpoint) ServerVersion() (ver m.Version) {

	ver = m.Version{
		Version:     version.VersionString,
		Revision:    version.RevisionString,
		RevisionURL: version.RevisionURLString,
		Generated:   version.GeneratedString,
		Developers:  version.DevelopersString,
		BuildNum:    version.BuildNumString,
		DockerImage: version.DockerImageString,
	}

	return
}
