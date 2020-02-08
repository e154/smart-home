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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/uuid"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"time"
)

const (
	gateVarName = "gateClientParams"
)

var (
	log = logging.MustGetLogger("gate")
)

type GateClient struct {
	sync.Mutex
	adaptors        *adaptors.Adaptors
	wsClient        *WsClient
	subscribers     map[uuid.UUID]func(msg Message)
	engine          *gin.Engine
	stream          *stream.StreamService
	messagePoolQuit chan struct{}
	messagePool     chan Message
	settingsLock    sync.Mutex
	settings        *Settings
	quit            bool
}

func NewGateClient(adaptors *adaptors.Adaptors,
	graceful *graceful_service.GracefulService,
	stream *stream.StreamService) (gate *GateClient) {
	gate = &GateClient{
		adaptors:        adaptors,
		settings:        &Settings{},
		subscribers:     make(map[uuid.UUID]func(msg Message)),
		stream:          stream,
		messagePoolQuit: make(chan struct{}),
		messagePool:     make(chan Message),
	}

	gate.wsClient = NewWsClient(gate)

	graceful.Subscribe(gate)

	if err := gate.loadSettings(); err != nil {
		log.Error(err.Error())
	}

	stream.GateClient(gate)

	go func() {
		for {
			select {
			case <-gate.messagePoolQuit:
				return
			case v := <-gate.messagePool:
				gate.settingsLock.Lock()
				if gate.quit {
					gate.settingsLock.Unlock()
					return
				}
				gate.settingsLock.Unlock()
				gate._onMessage(v)
			}
		}
	}()

	return
}

func (g *GateClient) Shutdown() {
	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()

	log.Info("Stopping")

	g.quit = true
	g.messagePoolQuit <- struct{}{}
	g.wsClient.Close()
}

func (g *GateClient) Close() {
	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()

	log.Info("Close")
	g.wsClient.Close()
}

func (g *GateClient) Restart() {
	g.Close()
	g.BroadcastAccessToken()
}

func (g *GateClient) BroadcastAccessToken() {
	log.Info("Broadcast access token")

	broadcastMsg := &stream.Message{
		Command: "gate.access_token",
		Type:    stream.Broadcast,
		Forward: stream.Request,
		Payload: map[string]interface{}{
			"accessToken": g.settings.GateServerToken,
		},
	}
	g.stream.Broadcast(broadcastMsg.Pack())

}

func (g *GateClient) RegisterServer() {

	g.settingsLock.Lock()
	log.Info("Register server")
	if g.settings.GateServerToken != "" || g.quit {
		g.settingsLock.Unlock()
		return
	}
	g.settingsLock.Unlock()

	payload := map[string]interface{}{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = g.Send("register_server", payload, ctx, func(msg Message) {

		if _, ok := msg.Payload["token"]; !ok {
			log.Errorf("no token in message payload")
			return
		}

		g.settingsLock.Lock()
		g.settings.GateServerToken = msg.Payload["token"].(string)
		settings := *g.settings
		g.settingsLock.Unlock()

		g.wsClient.UpdateSettings(settings)

		g.Restart()

		_ = g.saveSettings()
	})
}

func (g *GateClient) registerMobile(ctx *gin.Context) {

	params := map[string]interface{}{}
	_ = g.Send("register_mobile", params, ctx, func(msg Message) {

		if _, ok := msg.Payload["token"]; !ok {
			log.Errorf("no token in message payload")
			return
		}

		fmt.Println("mobile token ", msg.Payload["token"])
	})
}

func (g *GateClient) loadSettings() (err error) {
	log.Info("Load settings")

	var variable *m.Variable
	if variable, err = g.adaptors.Variable.GetByName(gateVarName); err != nil {
		if err = g.saveSettings(); err != nil {
			log.Error(err.Error())
		}
		return
	}

	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()

	if err = variable.GetObj(g.settings); err != nil {
		log.Error(err.Error())
	}

	if g.settings.Address == "" {
		log.Info("no gate address")
	}

	g.wsClient.UpdateSettings(*g.settings)

	return
}

func (g *GateClient) saveSettings() (err error) {
	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()

	log.Info("Save settings")

	variable := m.NewVariable(gateVarName)
	if err = variable.SetObj(g.settings); err != nil {
		return
	}

	err = g.adaptors.Variable.Update(variable)

	return
}

func (g *GateClient) GetSettings() (Settings, error) {
	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()

	return *g.settings, nil
}

func (g *GateClient) UpdateSettings(settings Settings) (err error) {

	g.settingsLock.Lock()
	if g.settings.Equal(settings) {
		g.settingsLock.Unlock()
		return
	}
	g.settingsLock.Unlock()

	var uri *url.URL
	if uri, err = url.Parse(settings.Address); err != nil {
		return
	}

	log.Infof("update gate settings, address: %v, enabled: %v", settings.Address, settings.Enabled)

	settings.Address = uri.String()

	g.settingsLock.Lock()
	g.settings.GateServerToken = settings.GateServerToken
	g.settings.Address = settings.Address
	g.settings.Enabled = settings.Enabled
	g.settingsLock.Unlock()

	g.wsClient.UpdateSettings(settings)

	if err = g.saveSettings(); err != nil {
		return
	}

	return
}

func (g *GateClient) onMessage(b []byte) {

	//log.Debugf("message(%v)\n", string(b))

	msg, err := NewMessage(b)
	if err != nil {
		log.Error(err.Error())
		return
	}

	g.messagePool <- msg

}

func (g *GateClient) _onMessage(msg Message) {

	if msg.Command == "mobile_gate_proxy" {
		g.RequestFromMobileProxy(msg)
		return
	}

	// check local subscribers
	for command, f := range g.subscribers {
		if msg.Id == command {
			f(msg)
			g.UnSubscribe(msg.Id)
			return
		}
	}

	// check subscriber on stream server
	if f := g.stream.Hub.Subscriber(msg.Command); f != nil {

		streamMsg := stream.Message{}
		_ = common.Copy(&streamMsg, &msg)
		streamClient := &stream.Client{
			ConnType: stream.WEBSOCK,
			Connect:  g.wsClient.conn,
			Send:     make(chan []byte),
		}

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			for {
				select {
				case message, ok := <-streamClient.Send:

					if !ok {
						_ = streamClient.Write(websocket.CloseMessage, []byte{})
						return
					}
					if err := streamClient.Write(websocket.TextMessage, message); err != nil {
						return
					}

					goto END
				}
			}
		END:
			wg.Done()
		}()

		go func() {
			f(streamClient, streamMsg)
			wg.Done()
		}()

		// check channels
		wg.Wait()
		close(streamClient.Send)
		//log.Debugf("client was success")
	}
}

func (g *GateClient) onConnected() {
	g.RegisterServer()
}

func (g *GateClient) onClosed() {

}

func (g *GateClient) Subscribe(id uuid.UUID, f func(msg Message)) {
	g.Lock()
	if g.subscribers[id] != nil {
		delete(g.subscribers, id)
	}
	g.subscribers[id] = f
	g.Unlock()
}

func (g *GateClient) UnSubscribe(id uuid.UUID) {
	g.Lock()
	if g.subscribers[id] != nil {
		delete(g.subscribers, id)
	}
	g.Unlock()
}

func (g *GateClient) Send(command string, payload map[string]interface{}, ctx context.Context, f func(msg Message)) (err error) {

	if g.wsClient.Status() != GateStatusConnected {
		err = errors.New("gate not connected")
		return
	}

	done := make(chan struct{})

	message := Message{
		Id:      uuid.NewV4(),
		Command: command,
		Payload: payload,
	}

	g.Subscribe(message.Id, func(msg Message) {
		f(msg)
		done <- struct{}{}
	})
	defer g.UnSubscribe(message.Id)

	msg, _ := json.Marshal(message)
	if err := g.wsClient.write(websocket.TextMessage, msg); err != nil {
		log.Error(err.Error())
	}

	select {
	case <-time.After(2 * time.Second):
	case <-done:
	case <-ctx.Done():
	}

	return
}

func (g *GateClient) Broadcast(message []byte) {
	if g.wsClient.Status() != GateStatusConnected {
		return
	}

	if err := g.wsClient.write(websocket.TextMessage, message); err != nil {
		log.Error(err.Error())
	}
}

func (g *GateClient) Status() string {

	status := g.wsClient.Status()
	if !g.settings.Enabled {
		return "disabled"
	}

	if status == "quit" {
		return "wait"
	}
	return status
}

func (g *GateClient) GetMobileList(ctx context.Context) (list *MobileList, err error) {

	list = &MobileList{
		TokenList: make([]string, 0),
	}

	payload := map[string]interface{}{}
	_ = g.Send("mobile_token_list", payload, ctx, func(msg Message) {
		if err = msg.IsError(); err != nil {
			return
		}
		if err = common.Copy(&list, msg.Payload, common.JsonEngine); err != nil {
			return
		}
	})

	return
}

func (g *GateClient) DeleteMobile(token string, ctx context.Context) (list *MobileList, err error) {

	payload := map[string]interface{}{
		"token": token,
	}
	_ = g.Send("remove_mobile", payload, ctx, func(msg Message) {
		err = msg.IsError()
	})

	return
}

func (g *GateClient) AddMobile(ctx context.Context) (list *MobileList, err error) {

	payload := map[string]interface{}{}
	_ = g.Send("register_mobile", payload, ctx, func(msg Message) {
		err = msg.IsError()
	})

	return
}

func (g *GateClient) RequestFromMobileProxy(message Message) {

	if g.wsClient.Status() != GateStatusConnected {
		return
	}

	//debug.Println(message.Payload["request"])

	if _, ok := message.Payload["request"]; !ok {
		log.Error("no request field from payload")
		return
	}

	requestParams := &StreamRequestModel{}
	if err := common.Copy(&requestParams, message.Payload["request"], common.JsonEngine); err != nil {
		log.Error(err.Error())
		return
	}

	payloadResponse := g.execRequest(requestParams)

	response := Message{
		Id:      uuid.NewV4(),
		Command: message.Id.String(),
		Payload: map[string]interface{}{
			"response": payloadResponse,
		},
	}

	msg, _ := json.Marshal(response)
	if err := g.wsClient.write(websocket.TextMessage, msg); err != nil {
		log.Error(err.Error())
	}
}

func (g *GateClient) SetEngine(engine *gin.Engine) {
	g.engine = engine
}

func (g *GateClient) execRequest(requestParams *StreamRequestModel) (response *StreamResponseModel) {

	if g.engine == nil {
		return
	}

	request, _ := http.NewRequest(requestParams.Method, requestParams.URI, bytes.NewBuffer(requestParams.Body))
	request.Header = requestParams.Header
	request.RequestURI = requestParams.URI

	//fmt.Println("----------")
	//fmt.Println("Request")
	//fmt.Println("----------")
	//debug.Println(request.RequestURI)
	//debug.Println(request.Header)

	recorder := httptest.NewRecorder()
	g.engine.ServeHTTP(recorder, request)

	//fmt.Println("----------")
	//fmt.Println("response")
	//fmt.Println("----------")
	//fmt.Println(recorder.Code)
	//fmt.Println(recorder.Header())
	//fmt.Println(recorder.Body)

	code := recorder.Code
	header := recorder.Header()
	body := recorder.Body.Bytes()

	//if err != nil {
	//	log.Error(err.Error())
	//	code = 500
	//	body = []byte(err.Error())
	//	header = http.Header{}
	//	header.Set("Content-Type", "text/plain")
	//}

	response = &StreamResponseModel{
		Code:   code,
		Body:   body,
		Header: header,
	}

	return
}
