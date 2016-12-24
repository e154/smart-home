package main

import (
	"time"
	"./settings"
	"log"
	"fmt"
	"encoding/hex"
	"flag"
	"net/rpc"
	r "./lib/rpc"
)

var (
	st	*settings.Settings
	client	*rpc.Client
	ip 	string
	port 	int
	baud	int
	com	string
	device	string
	counter int
	last int
	errors 	[]int
)

func testNode(command []byte) {

	args := r.Request{
		Baud: 19200,
		Result: true,
		Command: command,
		Device: st.Device,
		Line: "",
		StopBits: 2,
		Timeout: st.Timeout,
	}

	log.Println("send -> ", args)

	version := ""
	err := client.Call("Node.Version", "", &version)
	if err != nil {
		log.Println("error: ", err)
	}

	log.Printf("Node version %s\r\n", version)

	for i := 0; i< st.Iterations; i++ {

		result := &r.Result{}

		err := client.Call("Modbus.Send", args, result)

		if err == nil {
			if len(result.Result) == 0 {
				counter = i - last
				last = i
				errors = append(errors, counter)
			}
		}

		log.Printf("counter: %d, data %v, %d",i, result, len(errors))

		if err != nil {
			log.Println("error: ", err)
		}

		//time.Sleep(time.Millisecond * 50)
	}
}

func main() {

	// settings
	st = settings.SettingsPtr()
	st.Init()

	errors = []int{}

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
	client, err = rpc.Dial("tcp",fmt.Sprintf("%s:%d", st.IP, st.Port))
	if err != nil {
		log.Println(err.Error())
		return
	}

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