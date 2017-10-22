package main

import (
	"os"
)

var (
	VersionString = "?"
	RevisionString = "?"
	RevisionURLString = "?"
	GeneratedString = "?"
	DevelopersString = "?"
	BuildNumString = "?"
)

const verboseVersionBanner string = `
 ___                _     _  _
/ __|_ __  __ _ _ _| |_  | || |___ _ __  ___
\__ \ '  \/ _' | '_|  _| | __ / _ \ '  \/ -_)
|___/_|_|_\__,_|_|  \__| |_||_\___/_|_|_\___|	%s

Usage: %s [option]

options:
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
