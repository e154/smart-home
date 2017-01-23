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

func init() {
	_ = os.Setenv("VERSION", VersionString)
	_ = os.Setenv("REVISION", RevisionString)
	_ = os.Setenv("REVISION_URL", RevisionURLString)
	_ = os.Setenv("GENERATED", GeneratedString)
	_ = os.Setenv("DEVELOPERS", DevelopersString)
	_ = os.Setenv("BUILD_NUMBER", BuildNumString)
}
