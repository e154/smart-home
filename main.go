package main

import (
	"os"
	"log"
)

var (
	stdlog, errlog *log.Logger
)

func main() {

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