package main

import (
	"./serial"
	"./settings"
	"path/filepath"
	"os"
	"log"
	"fmt"
)

func run(dir string) {

	s := settings.SettingsPtr()
	s.Init(dir)

	fmt.Printf("start node v%s\n", s.AppVresion())

	serial.Init()
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	run(dir)
}

func init() {

}