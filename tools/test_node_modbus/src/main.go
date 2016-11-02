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
	"./lib/pack"
	"flag"
)

var (
	st	*settings.Settings
	conn	net.Conn
	ip 	string
	port 	int
	baud	int
	com	string
	device	string
)

func testNode(command []byte) {

	j, _ := json.Marshal(&map[string]interface {}{
		"line": "modbus",
		"device": st.Device,
		"baud": st.Baud,
		"timeout": time.Millisecond * st.Timeout,
		"result": true,
		"command": command,
	})

	j = append(j, '\n')

	log.Printf("send -> %s\r\n", string(j))


	for i := 0; i< st.Iterations; i++ {
		_, err := io.Copy(conn, bytes.NewBuffer(j))
		if err != nil {
			log.Printf("error %s\r\n", err.Error())
			return
		}

		for {
			res, err := bufio.NewReader(conn).ReadBytes('\n')
			if err != nil {
				log.Printf("error %s\r\n", err.Error())
			}

			r := &pack.Result{}
			if err := json.Unmarshal(res, r); err != nil {
				log.Printf("error %s\r\n", err.Error())
			}

			if r.Error != "" {
				i = st.Iterations
				log.Printf("error %s\r\n", r.Error)
			}

			log.Println("data", r)

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

	command, err := hex.DecodeString(st.Command)
	if err != nil {
		log.Printf("error %s\r\n",err.Error())
		return
	}

	log.Println("##########################")
	log.Println("# modbus test")
	log.Println("##########################")
	log.Printf("command: %d\r\n", command)
	log.Printf("connect %s:%d\r\n", st.IP, st.Port)

	// connect to node
	conn, err = net.Dial("tcp",fmt.Sprintf("%s:%d", st.IP, st.Port))
	if err != nil {
		beego.Debug(err.Error())
		return
	}

	defer conn.Close()

	testNode(command)

	log.Println("##########################")
}

func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "node service address")
	flag.IntVar(&port, "p", 3001, "node service port")
	flag.IntVar(&baud, "b", 19200, "serial port baud")
	flag.StringVar(&device, "d", "", "serial port device")
	flag.StringVar(&com, "c", "010300000005", "serial port command")
}