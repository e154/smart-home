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
	"github.com/e154/smart-home/system/metrics"
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

// WsClient ...
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
	connLock       *sync.Mutex
	conn           *websocket.Conn
	metric         *metrics.MetricManager
}

// NewWsClient ...
func NewWsClient(cb IWsCallback, metric *metrics.MetricManager) *WsClient {
	client := &WsClient{
		interrupt:  make(chan struct{}),
		quitWorker: make(chan struct{}),
		cb:         cb,
		setLock:    &sync.Mutex{},
		connLock:   &sync.Mutex{},
		metric:     metric,
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

// UpdateSettings ...
func (ws *WsClient) UpdateSettings(settings Settings) {
	ws.setLock.Lock()
	defer ws.setLock.Unlock()

	ws.settings = settings
	ws.reConnect = true
	ws.settingsLoaded = true
	ws.updateMetric()
}

func (ws *WsClient) connect() {

	ws.setLock.Lock()
	ws.selfUpdateStatus(GateStatusWait)

	if !ws.settings.Valid() {
		ws.setLock.Unlock()
		return
	}

	if !ws.settings.Enabled {
		ws.setLock.Unlock()
		return
	}
	ws.setLock.Unlock()

	var err error
	ticker := time.NewTicker(pingPeriod)
	defer func() {

		ws.setLock.Lock()
		ws.selfUpdateStatus(GateStatusNotConnected)
		ws.reConnect = false
		ws.setLock.Unlock()

		ticker.Stop()

		go ws.cb.onClosed()

		if ws.conn != nil {
			_ = ws.conn.Close()
		}
	}()

	var uri *url.URL
	if uri, err = url.Parse(ws.settings.Address); err != nil {
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
	if ws.settings.GateServerToken != "" {
		requestHeader.Add("X-API-Key", ws.settings.GateServerToken)
		//log.Debugf("X-API-Key: %v", ws.settings.GateServerToken)
	}

	ws.setLock.Lock()
	if ws.conn, _, err = websocket.DefaultDialer.Dial(uri.String(), requestHeader); err != nil {
		ws.setLock.Unlock()
		return
	}

	log.Infof("endpoint %v connected ...", uri.String())
	ws.selfUpdateStatus(GateStatusConnected)
	ws.setLock.Unlock()

	loseChan := make(chan struct{})
	var messageType int
	var message []byte

	//ws.conn.SetCloseHandler(func(code int, text string) error {
	//	log.Warn("connection closed")
	//
	//	loseChan <- struct{}{}
	//	return nil
	//})

	go ws.cb.onConnected()

	if err = ws.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}

	ws.conn.SetPongHandler(func(string) error {
		_ = ws.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		defer close(loseChan)
		for {

			if messageType, message, err = ws.conn.ReadMessage(); err != nil {
				//log.Error(err.Error())
				loseChan <- struct{}{}
				break
			}
			switch messageType {
			case websocket.TextMessage:
				//fmt.Printf("recv: %s\n", string(message))
				go ws.cb.onMessage(message)
			default:
				log.Warnf("unknown message type(%v)", messageType)
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			ws.setLock.Lock()
			if ws.reConnect {
				//_ = ws.write(websocket.CloseMessage, []byte{})
				log.Info("reconnect...")
				ws.setLock.Unlock()
				return
			}
			ws.setLock.Unlock()

			if err := ws.selfWrite(websocket.PingMessage, []byte{}); err != nil {
				log.Error(err.Error())
				return
			}
		case <-ws.interrupt:
			log.Info("Disconnected...")
			return
		case <-loseChan:
			return
		}
	}
}

// Close ...
func (ws *WsClient) Close() {

	ws.setLock.Lock()
	defer ws.setLock.Unlock()

	if ws.status == GateStatusQuit {
		return
	}

	if ws.status == GateStatusConnected {
		ws.interrupt <- struct{}{}
	}

	ws.selfUpdateStatus(GateStatusQuit)
}

func (ws *WsClient) selfWrite(opCode int, payload []byte) (err error) {

	ws.setLock.Lock()
	if ws.status != GateStatusConnected {
		ws.setLock.Unlock()
		return
	}
	ws.setLock.Unlock()

	ws.connLock.Lock()
	defer ws.connLock.Unlock()

	if err = ws.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return
	}

	err = ws.conn.WriteMessage(opCode, payload)

	return
}

// Write ...
func (ws *WsClient) Write(payload []byte) (err error) {
	err = ws.selfWrite(websocket.TextMessage, payload)
	return
}

// Status ...
func (ws *WsClient) Status() string {
	ws.setLock.Lock()
	defer ws.setLock.Unlock()

	return ws.status
}

// Notify ...
func (ws *WsClient) Notify(t, b string) {

	msg := &stream.Message{
		Type:    stream.Notify,
		Forward: Request,
		Payload: map[string]interface{}{
			"type": t,
			"body": b,
		},
	}

	ws.selfWrite(websocket.TextMessage, msg.Pack())
}

func (ws *WsClient) selfUpdateStatus(status string) {
	ws.status = status
	ws.updateMetric()
}

func (ws *WsClient) updateMetric() {

	//var status = ws.status
	//
	//if ws.status == "quit" {
	//	status = "wait"
	//}
	//
	//if !ws.settings.Enabled {
	//	status = "disabled"
	//}
	//
	//go ws.metric.Update(metrics.GateUpdate{
	//	Status:      status,
	//	AccessToken: ws.settings.GateServerToken,
	//})
}
