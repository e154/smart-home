package main

import (
	"os"
	"log"
	"strings"
)

var (
	stdlog, errlog *log.Logger
)

func main() {

	cwd := os.Getenv("PWD")
	rootDir := strings.Split(cwd, "tests/")
	os.Args[0] = rootDir[0]

	// just start
	args := os.Args
	if len(args) == 1 {
		stdlog.Printf(shortVersionBanner, "")
		ServiceInitialize()
		return
	}

	switch args[1] {
	case "install", "remove", "start", "stop", "status":
		ServiceInitialize()
	default:
		stdlog.Printf(verboseVersionBanner, "", args[0])
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}