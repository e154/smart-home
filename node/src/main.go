package main

import (
	"./settings"
	"./server"
	"log"
)

func main() {
	st := settings.SettingsPtr()
	st.Init()

	log.Printf("Start node v%s\n", st.AppVresion())

	//263
	sr := server.ServerPtr()
	if err := sr.Start(st.IP, st.Port); err != nil {
		log.Println(err.Error())
	}

	for {}
}