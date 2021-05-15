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
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/uuid"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
	"go.uber.org/fx"
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
	log = common.MustGetLogger("gate")
)

// GateClient ...
type GateClient struct {
	sync.Mutex
	metric          *metrics.MetricManager
	adaptors        *adaptors.Adaptors
	wsClient        *WsClient
	mobileApi       *gin.Engine
	alexaApi        *gin.Engine
	messagePool     chan stream.Message
	settingsLock    sync.Mutex
	settings        *Settings
	selfSubscrLock  sync.Mutex
	selfSubscribers map[uuid.UUID]func(msg stream.Message)
	subscrLock      sync.Mutex
	subscribers     map[string]func(client stream.IStreamClient, msg stream.Message)
	isStarted       *atomic.Bool
}

// NewGateClient ...
func NewGateClient(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors,
	metric *metrics.MetricManager) (gate *GateClient) {

	gate = &GateClient{
		adaptors:        adaptors,
		settings:        &Settings{},
		selfSubscribers: make(map[uuid.UUID]func(msg stream.Message)),
		subscribers:     make(map[string]func(client stream.IStreamClient, msg stream.Message)),
		metric:          metric,
		isStarted:       atomic.NewBool(false),
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			gate.Shutdown()
			return nil
		},
	})

	return
}

// Start ...
func (g *GateClient) Start() {
	if g.isStarted.Load() {
		return
	}
	g.isStarted.Store(true)

	g.wsClient = NewWsClient(g, g.metric)

	if err := g.loadSettings(); err != nil {
		log.Error(err.Error())
	}

	g.messagePool = make(chan stream.Message, 50)

	go func() {
		for v := range g.messagePool {
			g._onMessage(v)
		}
	}()

	log.Info("Start")
}

// Shutdown ...
func (g *GateClient) Shutdown() {
	if !g.isStarted.Load() {
		return
	}
	g.isStarted.Store(false)

	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()

	close(g.messagePool)
	g.wsClient.Close()

	log.Info("Shutdown")
}

// Close ...
func (g *GateClient) Close() {
	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()

	log.Info("Close")
	g.wsClient.Close()
}

// Restart ...
func (g *GateClient) Restart() {
	g.Close()

	go g.metric.Update(metrics.GateUpdate{
		AccessToken: g.settings.GateServerToken,
	})
}

// RegisterServer ...
func (g *GateClient) RegisterServer() {

	g.settingsLock.Lock()
	log.Info("Register server")
	if g.settings.GateServerToken != "" || !g.isStarted.Load() {
		g.settingsLock.Unlock()
		return
	}
	g.settingsLock.Unlock()

	payload := map[string]interface{}{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = g.Send("register_server", payload, ctx, func(msg stream.Message) {

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
	_ = g.Send("register_mobile", params, ctx, func(msg stream.Message) {

		if _, ok := msg.Payload["token"]; !ok {
			log.Errorf("no token in message payload")
			return
		}

		fmt.Println("mobile token ", msg.Payload["token"])
	})
}

func (g *GateClient) loadSettings() (err error) {
	log.Info("Load settings")

	var variable m.Variable
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

// GetSettings ...
func (g *GateClient) GetSettings() (Settings, error) {
	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()

	return *g.settings, nil
}

// UpdateSettings ...
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

	msg, err := stream.NewMessage(b)
	if err != nil {
		log.Error(err.Error())
		return
	}

	g.messagePool <- msg

}

func (g *GateClient) _onMessage(msg stream.Message) {

	//log.Debugf("message(%v)\n", msg)

	switch msg.Command {
	case MobileGateProxy:
		g.RequestFromProxy(msg, g.mobileApi)
		return
	case AlexaGateProxy:
		g.RequestFromProxy(msg, g.alexaApi)
		return
	}

	// check local selfSubscribers
	for command, f := range g.selfSubscribers {
		if msg.Id == command {
			f(msg)
			g.selfUnSubscribe(msg.Id)
			return
		}
	}

	// check subscriber on stream server
	if f, ok := g.subscribers[msg.Command]; ok {
		f(g.wsClient, msg)
	} else {
		log.Warnf("unknown command %v", msg.Command)
	}
}

func (g *GateClient) onConnected() {
	g.RegisterServer()
}

func (g *GateClient) onClosed() {

}

func (g *GateClient) selfSubscribe(id uuid.UUID, f func(msg stream.Message)) {
	g.selfSubscrLock.Lock()
	defer g.selfSubscrLock.Unlock()

	if g.selfSubscribers[id] != nil {
		delete(g.selfSubscribers, id)
	}
	g.selfSubscribers[id] = f
}

func (g *GateClient) selfUnSubscribe(id uuid.UUID) {
	g.selfSubscrLock.Lock()
	defer g.selfSubscrLock.Unlock()

	if g.selfSubscribers[id] != nil {
		delete(g.selfSubscribers, id)
	}
}

// Subscribe ...
func (g *GateClient) Subscribe(command string, f func(client stream.IStreamClient, msg stream.Message)) {
	g.subscrLock.Lock()
	defer g.subscrLock.Unlock()

	if g.subscribers[command] != nil {
		delete(g.subscribers, command)
	}
	g.subscribers[command] = f
}

// UnSubscribe ...
func (g *GateClient) UnSubscribe(command string) {
	g.subscrLock.Lock()
	defer g.subscrLock.Unlock()

	if g.subscribers[command] != nil {
		delete(g.subscribers, command)
	}
}

// Send ...
func (g *GateClient) Send(command string, payload map[string]interface{}, ctx context.Context, f func(msg stream.Message)) (err error) {

	if g.wsClient.Status() != GateStatusConnected {
		err = errors.New("gate not connected")
		return
	}

	done := make(chan struct{})

	message := stream.Message{
		Id:      uuid.NewV4(),
		Command: command,
		Payload: payload,
	}

	g.selfSubscribe(message.Id, func(msg stream.Message) {
		f(msg)
		done <- struct{}{}
	})
	defer g.selfUnSubscribe(message.Id)

	msg, _ := json.Marshal(message)
	if err := g.wsClient.selfWrite(websocket.TextMessage, msg); err != nil {
		log.Error(err.Error())
	}

	select {
	case <-time.After(2 * time.Second):
	case <-done:
	case <-ctx.Done():
	}

	return
}

// Broadcast ...
func (g *GateClient) Broadcast(message []byte) {
	if g.wsClient.Status() != GateStatusConnected {
		return
	}

	if err := g.wsClient.selfWrite(websocket.TextMessage, message); err != nil {
		log.Error(err.Error())
	}
}

// Status ...
func (g *GateClient) Status() string {

	if !g.settings.Enabled {
		return "disabled"
	}

	status := g.wsClient.Status()
	if status == "quit" {
		return "wait"
	}
	return status
}

// GetMobileList ...
func (g *GateClient) GetMobileList(ctx context.Context) (list *MobileList, err error) {

	list = &MobileList{
		TokenList: make([]string, 0),
	}

	payload := map[string]interface{}{}
	_ = g.Send("mobile_token_list", payload, ctx, func(msg stream.Message) {
		if err = msg.IsError(); err != nil {
			return
		}
		if err = common.Copy(&list, msg.Payload, common.JsonEngine); err != nil {
			return
		}
	})

	return
}

// DeleteMobile ...
func (g *GateClient) DeleteMobile(token string, ctx context.Context) (list *MobileList, err error) {

	payload := map[string]interface{}{
		"token": token,
	}
	_ = g.Send("remove_mobile", payload, ctx, func(msg stream.Message) {
		err = msg.IsError()
	})

	return
}

// AddMobile ...
func (g *GateClient) AddMobile(ctx context.Context) (list *MobileList, err error) {

	payload := map[string]interface{}{}
	_ = g.Send("register_mobile", payload, ctx, func(msg stream.Message) {
		err = msg.IsError()
	})

	return
}

// RequestFromProxy ...
func (g *GateClient) RequestFromProxy(message stream.Message, engine *gin.Engine) {

	if engine == nil {
		return
	}

	if g.wsClient.Status() != GateStatusConnected {
		return
	}

	//debug.Println(message.Obj["request"])

	if _, ok := message.Payload["request"]; !ok {
		log.Error("no request field from payload")
		return
	}

	requestParams := &StreamRequestModel{}
	if err := common.Copy(&requestParams, message.Payload["request"], common.JsonEngine); err != nil {
		log.Error(err.Error())
		return
	}

	payloadResponse := g.execRequest(requestParams, engine)

	response := stream.Message{
		Id:      uuid.NewV4(),
		Command: message.Id.String(),
		Payload: map[string]interface{}{
			"response": payloadResponse,
		},
	}

	msg, _ := json.Marshal(response)
	if err := g.wsClient.selfWrite(websocket.TextMessage, msg); err != nil {
		log.Error(err.Error())
	}
}

// SetMobileApiEngine ...
func (g *GateClient) SetMobileApiEngine(engine *gin.Engine) {
	g.mobileApi = engine
}

// SetAlexaApiEngine ...
func (g *GateClient) SetAlexaApiEngine(engine *gin.Engine) {
	g.alexaApi = engine
}

func (g *GateClient) execRequest(requestParams *StreamRequestModel, engine *gin.Engine) (response *StreamResponseModel) {

	if engine == nil {
		return
	}

	request, _ := http.NewRequest(requestParams.Method, requestParams.URI, bytes.NewBuffer(requestParams.Body))
	request.Header = requestParams.Header
	request.RequestURI = requestParams.URI
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	code := recorder.Code
	header := recorder.Header()
	body := recorder.Body.Bytes()
	response = &StreamResponseModel{
		Code:   code,
		Body:   body,
		Header: header,
	}

	return
}
