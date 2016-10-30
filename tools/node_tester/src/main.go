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
)

var (
	st	*settings.Settings
	conn	net.Conn
)

func testNode(command []byte) {

	j, _ := json.Marshal(&map[string]interface {}{
		"line": "modbus",
		"device": nil,
		"baud": nil,
		"timeout": time.Millisecond * st.Timeout,
		"result": true,
		"command": command,
	})

	j = append(j, '\n')

	log.Printf("send -> %s\r\n", string(j))


	for i := 0; i< st.Iterations; i++ {
		_, err := io.Copy(conn, bytes.NewBuffer(j))
		if err != nil {
			log.Println(999)
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
			log.Println("data", r)

			break
		}
	}
}

func main() {

	// settings
	st = settings.SettingsPtr()
	st.Init()

	command, err := hex.DecodeString(st.Command)
	if err != nil {
		log.Printf("error %s\r\n",err.Error())
		return
	}

	log.Println("##########################")
	log.Println("# node test")
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
