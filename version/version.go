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

package version

import (
	"fmt"
)

var (
	VersionString = "?"
	RevisionString = "?"
	RevisionURLString = "?"
	GeneratedString = "?"
	DevelopersString = "?"
	BuildNumString = "?"
	DockerImageString = "?"
)

const VerboseVersionBanner string = `
 ___                _     _  _
/ __|_ __  __ _ _ _| |_  | || |___ _ __  ___
\__ \ '  \/ _' | '_|  _| | __ / _ \ '  \/ -_)
|___/_|_|_\__,_|_|  \__| |_||_\___/_|_|_\___|	%s

Usage: %s [option]

options:
-v -version - show build version
-backup	    - backup settings to file
-restore    - restore settings from backup archive
-reset      - cleanup and restore base settings
help	    - show this help text
`

const ShortVersionBanner = `
 ___                _     _  _
/ __|_ __  __ _ _ _| |_  | || |___ _ __  ___
\__ \ '  \/ _' | '_|  _| | __ / _ \ '  \/ -_)
|___/_|_|_\__,_|_|  \__| |_||_\___/_|_|_\___|	

%s
`

func GetHumanVersion() string {
	version := ""

	if DevelopersString != "" {
		version = fmt.Sprintf("Generated: %s\n", GeneratedString)
	}

	if RevisionString != "" {
		version += fmt.Sprintf("Revision: %s\n", RevisionString)
	}

	if RevisionURLString != "" {
		version += fmt.Sprintf("Revision url: %s\n", RevisionURLString)
	}

	if VersionString != "" {
		version += fmt.Sprintf("Version: %s\n", VersionString)
	}

	if DockerImageString != "" {
		version += fmt.Sprintf("Docker image: %s\n", DockerImageString)
	}

	if DevelopersString != "" {
		version += fmt.Sprintf("Developers: %s\n", DevelopersString)
	}

	if BuildNumString != "" {
		version += fmt.Sprintf("Build: %s\n", BuildNumString)
	}

	return version
}
