package server

import (
	"net"
	"bufio"
	"io"
	"time"
	"log"
	"encoding/json"
	"bytes"
	"errors"
	"../serial"
	"../settings"
	"../lib/pack"
)

type Client struct {
	Conn net.Conn
}

func (c *Client) Send(b []byte) (int64, error) {
	if c.Conn == nil {
		log.Print("Client conn is nil")
		return 0, errors.New("Client conn is nil")
	}

	return  io.Copy(c.Conn, bytes.NewBuffer(b))
}

func (c *Client) listener(conn net.Conn) {

	re := &pack.Package{}

	for ;; {
		data, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			log.Println(5)
			c.Err(err)
			return
		}

		if err := json.Unmarshal(data, &re); err != nil {
			log.Println(6)
			c.Err(err)
			continue
		}

		switch re.Line {
		case "modbus":
			c.modBusExec(re)

		case "get_device_list":
			devices := serial.FindSerials()
			j, _ := json.Marshal(&map[string]interface {}{"line": "get_device_list", "result": devices, "time": time.Now()})
			c.Send(j)

		default:
			c.rawExec(re)
		}
	}
}

func (c *Client) modBusExec(p *pack.Package) {

	st := settings.SettingsPtr()

	conn := &serial.Serial{
		Dev: "",
		Baud: st.Baud,
		ReadTimeout: time.Second * st.Timeout,
		StopBits: st.StopBits,
	}

	if p.Device != "" {
		conn.Dev = p.Device
	}

	if p.Baud != 0 {
		conn.Baud = p.Baud
	}

	if p.Timeout != 0 {
		conn.ReadTimeout = p.Timeout
	}

	res, err := ModBusProxy(conn, p.Command)
	if err != nil {
		log.Println(3)
		c.Err(err)
		return
	}

	j, _ := json.Marshal(&pack.Result{
		Result:res,
		Command:p.Command,
		Device: p.Device,
		Line: p.Line,
		Time: time.Now(),
	})
	j = append(j, '\n')

	if _, err = c.Send(j); err != nil {
		log.Println(4)
		c.Err(err)
		return
	}
}

func (c *Client) rawExec(p *pack.Package) {

}

func (c *Client) Err(err error) {
	log.Print("error", err.Error())
	j, _ := json.Marshal(&pack.Error{Error: err.Error(), Time:time.Now()})
	c.Send(j)
}

