package main

import (
	"time"
	"./settings"
	"log"
	"fmt"
	"encoding/hex"
	"net"
	"github.com/astaxie/beego"
	"encoding/json"
	"io"
	"bytes"
	"bufio"
	"flag"
)

var (
	st	*settings.Settings
	l	net.Conn
	ip 	string
	port 	int
	baud	int
	com	string
	device	string
	command []byte
)

func LRC(data []byte) byte {

	var ucLRC uint8 = 0

	var b byte
	for _, b = range data {
		ucLRC += b
	}

	return uint8(-int8(ucLRC))
}

func char2bin(b byte) byte {

	if( b >= '0' ) && ( b <= '9' ) {
		return byte( b - '0' )
	} else if( ( b >= 'A' ) && ( b <= 'F' ) ) {
		return byte( b - 'A' + 0x0A )
	}

	return 0xFF
}

func testNode(command []byte) {

	j, _ := json.Marshal(&map[string]interface{}{
		"line": "serial",
		"device": st.Device,
		"baud": st.Baud,
		"timeout": time.Millisecond * st.Timeout,
		"result": true,
		"command": "",
	})

	j = append(j, '\r', '\n')

	if _, err := io.Copy(l, bytes.NewBuffer(j)); err != nil {
		log.Println("error: ", err.Error())
		return
	}

	command = []byte{0x01,0x03,0x00,0x00,0x00,0x05}
	command = append(command, LRC(command))

	log.Printf("send -> %X\r\n", command)
	command = append(command, '\r', '\n')

	for i:=0;i<1;i++ {
		if _, err := io.Copy(l, bytes.NewBuffer(command)); err != nil {
			log.Println("error: ", err.Error())
			return
		}

		for {
			data, err := bufio.NewReader(l).ReadBytes('\n')
			if err != nil {
				log.Println("error: ", err.Error())
				return
			}

			log.Println("data: ", string(data))

			break
		}
	}
}

func main() {

	// settings
	st = settings.SettingsPtr()
	st.Init()

	// args
	flag.Parse()
	if ip != st.IP {st.IP = ip}
	if com != st.Command {st.Command = com}
	if port != st.Port {st.Port = port}
	if baud != st.Baud {st.Baud = baud}
	if device != st.Device {st.Device = device}

	var err error
	command, err = hex.DecodeString(st.Command)
	if err != nil {
		log.Printf("error %s\r\n",err.Error())
		return
	}

	log.Println("##########################")
	log.Println("# modbus test")
	log.Println("##########################")
	log.Printf("connect %s:%d\r\n", st.IP, st.Port)

	// connect to node
	l, err = net.Dial("tcp",fmt.Sprintf("%s:%d", st.IP, st.Port))
	if err != nil {
		beego.Debug(err.Error())
		return
	}

	defer l.Close()

	testNode(command)

	log.Println("##########################")
}

func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "node service address")
	flag.IntVar(&port, "p", 3001, "node service port")
	flag.IntVar(&baud, "b", 19200, "serial port baud")
	flag.StringVar(&device, "d", "/dev/ttyUSB0", "serial port device")
	flag.StringVar(&com, "c", "010300000005", "serial port command")
}