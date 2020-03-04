// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package gate_client

import (
	"fmt"
	"github.com/e154/smart-home/system/stream"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type WsClient struct {
	settings       Settings
	settingsLoaded bool
	delta          time.Duration
	interrupt      chan struct{}
	quitWorker     chan struct{}
	cb             IWsCallback
	reConnect      bool
	setLock        *sync.Mutex
	status         string
	counter        int
	connLock       *sync.Mutex
	conn           *websocket.Conn
}

func NewWsClient(
	cb IWsCallback) *WsClient {
	client := &WsClient{
		interrupt:  make(chan struct{}),
		quitWorker: make(chan struct{}),
		cb:         cb,
		setLock:    &sync.Mutex{},
		connLock:   &sync.Mutex{},
	}

	go func() {
		for {
			client.setLock.Lock()
			status := client.status
			settingsLoaded := client.settingsLoaded
			client.setLock.Unlock()

			if status == GateStatusQuit {
				return
			}

			time.Sleep(time.Second)

			switch status {
			case GateStatusWait, GateStatusNotConnected:
			case "":
				if !settingsLoaded {
					continue
				}
			default:
				log.Infof("unknown status %v", status)
				continue
			}

			client.connect()
		}
	}()
	return client
}

func (client *WsClient) UpdateSettings(settings Settings) {
	client.setLock.Lock()
	defer client.setLock.Unlock()

	client.settings = settings
	client.reConnect = true
	client.settingsLoaded = true
}

func (client *WsClient) connect() {

	client.setLock.Lock()
	client.status = GateStatusWait

	if !client.settings.Valid() {
		client.setLock.Unlock()
		return
	}

	if !client.settings.Enabled {
		client.setLock.Unlock()
		return
	}
	client.setLock.Unlock()

	var err error
	ticker := time.NewTicker(pingPeriod)
	client.counter++
	defer func() {
		client.counter--

		client.setLock.Lock()
		client.status = GateStatusNotConnected
		client.reConnect = false
		client.setLock.Unlock()

		ticker.Stop()

		go client.cb.onClosed()

		if client.conn != nil {
			_ = client.conn.Close()
		}

		//if err != nil {
		//	if strings.Contains(err.Error(), "connection refused") {
		//		return
		//	}
		//	if strings.Contains(err.Error(), "bad handshake") {
		//		return
		//	}
		//	if strings.Contains(err.Error(), "use of closed network connection") {
		//		return
		//	}
		//	log.Debug(err.Error())
		//}
	}()

	var uri *url.URL
	if uri, err = url.Parse(client.settings.Address); err != nil {
		return
	}

	uri.Path = "ws"

	if uri.Scheme == "http" {
		uri.Scheme = "ws"
	} else {
		uri.Scheme = "wss"
	}

	requestHeader := http.Header{
		"X-Client-Type": {ClientTypeServer},
	}
	if client.settings.GateServerToken != "" {
		requestHeader.Add("X-API-Key", client.settings.GateServerToken)
		//log.Debugf("X-API-Key: %v", client.settings.GateServerToken)
	}

	client.setLock.Lock()
	if client.conn, _, err = websocket.DefaultDialer.Dial(uri.String(), requestHeader); err != nil {
		client.setLock.Unlock()
		return
	}

	log.Info("endpoint %v connected ...", uri.String())
	client.status = GateStatusConnected
	client.setLock.Unlock()

	loseChan := make(chan struct{})
	var messageType int
	var message []byte

	//client.conn.SetCloseHandler(func(code int, text string) error {
	//	log.Warning("connection closed")
	//
	//	loseChan <- struct{}{}
	//	return nil
	//})

	go client.cb.onConnected()

	if err = client.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}

	client.conn.SetPongHandler(func(string) error {
		_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		defer close(loseChan)
		for {

			if messageType, message, err = client.conn.ReadMessage(); err != nil {
				//log.Error(err.Error())
				loseChan <- struct{}{}
				break
			}
			switch messageType {
			case websocket.TextMessage:
				//fmt.Printf("recv: %s\n", string(message))
				go client.cb.onMessage(message)
			default:
				log.Warningf("unknown message type(%v)", messageType)
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			client.setLock.Lock()
			if client.reConnect {
				//_ = client.write(websocket.CloseMessage, []byte{})
				fmt.Println("reconnect...")
				client.setLock.Unlock()
				return
			}
			client.setLock.Unlock()

			if err := client.selfWrite(websocket.PingMessage, []byte{}); err != nil {
				log.Error(err.Error())
				return
			}
		case <-client.interrupt:
			log.Info("Disconnected...")
			return
		case <-loseChan:
			return
		}
	}
}

func (client *WsClient) Close() {

	client.setLock.Lock()
	defer client.setLock.Unlock()

	if client.status == GateStatusQuit {
		return
	}
	if client.status != GateStatusQuit {
		client.status = GateStatusQuit
		if client.status == GateStatusConnected {
			client.interrupt <- struct{}{}
		}
	}
}

func (client *WsClient) selfWrite(opCode int, payload []byte) (err error) {

	client.setLock.Lock()
	if client.status != GateStatusConnected {
		client.setLock.Unlock()
		return
	}
	client.setLock.Unlock()

	client.connLock.Lock()
	defer client.connLock.Unlock()

	if err = client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return
	}

	err = client.conn.WriteMessage(opCode, payload)

	return
}

func (client *WsClient) Write(payload []byte) (err error) {
	err = client.selfWrite(websocket.TextMessage, payload)
	return
}

func (client *WsClient) Status() string {
	client.setLock.Lock()
	defer client.setLock.Unlock()

	return client.status
}

func (client *WsClient) Notify(t, b string) {

	msg := &stream.Message{
		Type:    stream.Notify,
		Forward: Request,
		Payload: map[string]interface{}{
			"type": t,
			"body": b,
		},
	}

	client.selfWrite(websocket.TextMessage, msg.Pack())
}
