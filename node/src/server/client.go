package server

import (
	"net"
	"bufio"
	"io"
	"time"
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
		return 0, errors.New("Client conn is nil")
	}

	b = append(b, '\n')

	return  io.Copy(c.Conn, bytes.NewBuffer(b))
}

func (c *Client) listener(l net.Conn) {

	re := &pack.Package{}

	for ;; {
		data, err := bufio.NewReader(l).ReadBytes('\n')
		if err != nil {
			c.Err(err, TCP_READ_LINE_ERROR)
			return
		}

		if err := json.Unmarshal(data, &re); err != nil {
			c.Err(err, TCP_UNMARSHAL_ERROR)
			continue
		}

		switch re.Line {
		case "modbus":
			c.modBusExec(re)

		case "get_device_list":
			devices := serial.FindSerials()
			j, _ := json.Marshal(&map[string]interface {}{"line": "get_device_list", "result": devices})
			c.Send(j)

		case "serial":
			c.serialExec(re, l)

		default:
			c.serialExec(re, l)
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

	res, err, errcode := ModBusProxy(conn, p.Command)
	if err != nil {
		c.Err(err, errcode)
		return
	}

	j, _ := json.Marshal(&pack.Result{
		Result:res,
		Command:p.Command,
		Device: p.Device,
		Line: p.Line,
		Time: time.Now(),
	})

	if _, err = c.Send(j); err != nil {
		return
	}
}

func (c *Client) serialExec(p *pack.Package, l net.Conn) (err error, errcode int) {

	s := &serial.Serial{
		Dev: p.Device,
		Baud: p.Baud,
		ReadTimeout: time.Second * p.Timeout,
		StopBits: p.StopBits,
	}

	if _, err = s.Open(); err != nil {
		c.Err(err, SERIAL_PORT_ERROR)
		return
	}

	done := make(chan bool, 2)

	go serialToNet(l, s.Port, done)
	go netToSerial(l, s.Port, done)

	for {
		<-done
		s.Close()
	}
}

func (c *Client) Err(err error, errcode int) {
	j, _ := json.Marshal(&pack.Result{Error: err.Error(), ErrorCode: errcode, Time: time.Now()})
	c.Send(j)
}

