package main

import (
	"./settings"
	"./serial"
	"fmt"

)

func main() {
	s := settings.SettingsPtr()
	s.Init()

	fmt.Printf("start node v%s\n", s.AppVresion())

	//fmt.Println(serial.LRC([]byte{1,3,0,0,0,5}))

	serial.Run()
}