package main

import (
	"./settings"
	"./server"
	"./cache"
	"log"
	"time"
)

func main() {
	// settings
	st := settings.SettingsPtr()
	st.Init()

	// cache
	cache.Init(int64(st.Cachetime))

	log.Printf("Start node v%s\n", st.AppVresion())

	// tcp server
	sr := server.ServerPtr()
	if err := sr.Start(st.IP, st.Port); err != nil {
		log.Println(err.Error())
	}

	for ;; {
		time.Sleep(time.Second * 1)
	}
}