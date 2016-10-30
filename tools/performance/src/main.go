package main

import (
	"time"
	"./settings"
	"./serial"
	"log"
	"fmt"
	"encoding/hex"
)

var (
	cache 	string
	st	*settings.Settings
)


func exec(device string, command []byte) (res []byte, err error) {

	serial_port := &serial.Serial{
		Dev:device,
		Baud:st.Baud,
		ReadTimeout: time.Millisecond * st.Timeout,
		StopBits: st.StopBits,
	}

	if _, err = serial_port.Open(); err != nil {
		cache = ""
		log.Printf("error: %s - %s\r\n",device, err.Error())
		return
	}

	modbus := &serial.Modbus{Serial: serial_port}
	res, err = modbus.Send(command)
	if err != nil {
		log.Printf("error: %s - %s\r\n",device, err.Error())
		return
	}

	//cache update
	cache = device

	return
}

func performance(command []byte) (res []byte, err error) {

	if cache != "" {
		res, err = exec(cache, command)
		if err == nil {
			return
		}
	} else {
		log.Println("############ search ##############")
		devices := serial.FindSerials()
		for _, device := range devices {
			res, err = exec(device, command)
			if err == nil {
				return
			}
		}
	}

	return
}

func main() {

	// settings
	st = settings.SettingsPtr()
	st.Init()

	command, err := hex.DecodeString(st.Command)
	if err != nil {
		fmt.Printf("%s\r\n",err.Error())
		return
	}

	log.Println("##########################")
	log.Println("# performance test")
	log.Println("##########################")
	log.Println("command:")
	log.Println(command)

	var found, notfound, lrc int

	for i := 0; i<st.Iterations; i++ {
		res, err := performance(command)
		if err != nil {
			if err.Error() == "ILLEGAL_LRC" {
				lrc++
			}
			notfound++
		} else {
			found++
		}


		log.Printf("found %d, not found %d, lrc %d, res %d\r\n", found, notfound, lrc, res)

		//time.Sleep(time.Millisecond * 50)
	}

	log.Println("##########################")
	log.Printf("#found %d\r\n", found)
	log.Printf("#not found %d\r\n", notfound)
	log.Printf("#illegal lrc %d\r\n", lrc)
	log.Println("##########################")
}
