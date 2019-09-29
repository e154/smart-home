package gate_client

import (
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
	"github.com/jinzhu/copier"
	"github.com/op/go-logging"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

const (
	gateVarName = "gateClientParams"
)

var (
	log                 = logging.MustGetLogger("gate")
	ErrGateNotConnected = fmt.Errorf("gate not connected")
)

type GateClient struct {
	sync.Mutex
	adaptors    *adaptors.Adaptors
	settings    *Settings
	wsClient    *WsClient
	subscribers map[uuid.UUID]func(msg Message)
	engine      *gin.Engine
	stream      *stream.StreamService
}

func NewGateClient(adaptors *adaptors.Adaptors,
	graceful *graceful_service.GracefulService,
	stream *stream.StreamService) (gate *GateClient) {
	gate = &GateClient{
		adaptors:    adaptors,
		settings:    &Settings{},
		subscribers: make(map[uuid.UUID]func(msg Message)),
		stream:      stream,
	}

	gate.wsClient = NewWsClient(adaptors, gate)

	graceful.Subscribe(gate)

	if err := gate.LoadSettings(); err != nil {
		log.Error(err.Error())
	}

	return
}

func (g *GateClient) Shutdown() {
	g.wsClient.Close()
}

func (g *GateClient) Close() {
	log.Info("Close")
	g.wsClient.Close()
}

func (g *GateClient) Connect() {
	log.Info("Connect")
	if !g.settings.Valid() {
		return
	}

	g.wsClient.Connect(g.settings)
}

func (g *GateClient) Restart() {
	g.Close()
	g.Connect()
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
	log.Info("Register server")
	if g.settings.GateServerToken != "" {
		return
	}

	payload := map[string]interface{}{}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = g.Send("register_server", payload, ctx, func(msg Message) {

		if _, ok := msg.Payload["token"]; !ok {
			log.Errorf("no token in message payload")
			return
		}

		g.settings.GateServerToken = msg.Payload["token"].(string)

		g.Restart()

		_ = g.SaveSettings()
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

func (g *GateClient) LoadSettings() (err error) {
	log.Info("Load settings")

	var variable *m.Variable
	if variable, err = g.adaptors.Variable.GetByName(gateVarName); err != nil {
		if err = g.SaveSettings(); err != nil {
			log.Error(err.Error())
		}
		return
	}

	if err = variable.GetObj(g.settings); err != nil {
		log.Error(err.Error())
	}

	if g.settings.Address == "" {
		log.Info("no gate address")
	}

	return
}

func (g *GateClient) SaveSettings() (err error) {
	log.Info("Save settings")

	variable := m.NewVariable(gateVarName)
	if err = variable.SetObj(g.settings); err != nil {
		return
	}

	err = g.adaptors.Variable.Update(variable)

	return
}

func (g *GateClient) onMessage(b []byte) {

	//log.Debugf("message(%v)\n", string(b))

	msg, err := NewMessage(b)
	if err != nil {
		log.Error(err.Error())
		return
	}

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
		log.Debugf("client was success")
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

	if g.wsClient.status != GateStatusConnected {
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

func (g *GateClient) Status() string {
	return g.wsClient.status
}

func (g *GateClient) GetSettings() (*Settings, error) {
	return g.settings, nil
}

func (g *GateClient) UpdateSettings(settings *Settings) (err error) {
	if err = copier.Copy(&g.settings, &settings); err != nil {
		return
	}
	if err = g.SaveSettings(); err != nil {
		return
	}
	if !g.settings.Enabled {
		g.Close()
	}
	return
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

	if g.wsClient.status != GateStatusConnected {
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

	request, _ := http.NewRequest(requestParams.Method, requestParams.URI, nil)
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
