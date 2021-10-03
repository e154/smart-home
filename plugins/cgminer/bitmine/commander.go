// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package bitmine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Transport struct {
	server  string
	timeout time.Duration
	device  DeviceType
}

func NewTransport(host string, port int, timeout int64) ITransport {
	return &Transport{
		server:  fmt.Sprintf("%s:%d", host, port),
		timeout: time.Duration(timeout) * time.Second,
	}
}

func (tr Transport) RunCommand(command, argument string) (res []byte, err error) {
	var conn net.Conn
	if conn, err = net.DialTimeout("tcp", tr.server, tr.timeout); err != nil {
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(tr.timeout))

	request := &commandRequest{
		Command:   command,
		Parameter: argument,
	}

	requestBody, _ := json.Marshal(request)
	if _, err = conn.Write(requestBody); err != nil {
		return
	}
	var result []byte
	if result, err = bufio.NewReader(conn).ReadBytes(0x00); err != nil {
		return
	}
	res = bytes.TrimRight(result, "\x00")
	return
}
