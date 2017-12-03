package main

import (
	"os"
	"log"
	"strings"
	"runtime/trace"
)

var (
	stdlog, errlog *log.Logger
)

func main() {

	cwd := os.Getenv("PWD")
	rootDir := strings.Split(cwd, "tests/")
	os.Args[0] = rootDir[0]

	args := os.Args
	switch len(args) {
	case 1:
		stdlog.Printf(shortVersionBanner, "")
		ServiceInitialize()
	case 2:
		switch args[1] {
		case "trace":
			stdlog.Println("Trace mode enabled")
			f, err := os.Create("trace.out")
			if err != nil {
				panic(err)
			}

			defer f.Close()

			if err = trace.Start(f); err != nil {
				panic(err)
			}

			defer trace.Stop()

			ServiceInitialize()

		case "install", "remove", "start", "stop", "status":
			ServiceInitialize()

		default:
			stdlog.Printf(verboseVersionBanner, "", args[0])
		}

	default:
		stdlog.Printf(verboseVersionBanner, "", args[0])
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}