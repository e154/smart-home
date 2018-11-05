package main

import (
	"os"
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

const verboseVersionBanner string = `
 ___                _     _  _
/ __|_ __  __ _ _ _| |_  | || |___ _ __  ___
\__ \ '  \/ _' | '_|  _| | __ / _ \ '  \/ -_)
|___/_|_|_\__,_|_|  \__| |_||_\___/_|_|_\___|	%s

Usage: %s [option]

options:
trace	- trace mode for debug
start	- [sudo] Start the service
stop	- [sudo] Stop the service
install	- [sudo] Install the service into the system
remove	- [sudo] Remove the service and all corresponding files from the system
status	- [sudo] Check the service status
help	- show this help text

`

const shortVersionBanner = `
 ___                _     _  _
/ __|_ __  __ _ _ _| |_  | || |___ _ __  ___
\__ \ '  \/ _' | '_|  _| | __ / _ \ '  \/ -_)
|___/_|_|_\__,_|_|  \__| |_||_\___/_|_|_\___|	%s

`

func init() {
	_ = os.Setenv("VERSION", VersionString)
	_ = os.Setenv("REVISION", RevisionString)
	_ = os.Setenv("REVISION_URL", RevisionURLString)
	_ = os.Setenv("GENERATED", GeneratedString)
	_ = os.Setenv("DEVELOPERS", DevelopersString)
	_ = os.Setenv("BUILD_NUMBER", BuildNumString)
}

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