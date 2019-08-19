package gate_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/debug"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/uuid"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

const (
	gateVarName = "gateClientParams"
)

var (
	log                 = logging.MustGetLogger("gate")
	ErrGateNotConnected = fmt.Errorf("gate not connected")
)

type GateClient struct {
	adaptors    *adaptors.Adaptors
	settings    *Settings
	wsClient    *WsClient
	subscribers map[uuid.UUID]func(msg Message)
}

func NewGateClient(adaptors *adaptors.Adaptors,
	graceful *graceful_service.GracefulService, ) (gate *GateClient) {
	gate = &GateClient{
		adaptors:    adaptors,
		settings:    &Settings{},
		subscribers: make(map[uuid.UUID]func(msg Message)),
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
	g.wsClient.Close()
}

func (g *GateClient) Connect() {

	if !g.settings.Valid() {
		return
	}

	g.wsClient.Connect(g.settings)
}

func (g *GateClient) registerServer() {

	if g.settings.GateServerToken != "" {
		return
	}

	payload := map[string]interface{}{}
	g.Send("register_server", payload, func(msg Message) {

		if _, ok := msg.Payload["token"]; !ok {
			log.Errorf("no token in message payload")
			return
		}

		g.settings.GateServerToken = msg.Payload["token"].(string)

		_ = g.SaveSettings()
	})
}

func (g *GateClient) registerMobile() {

	params := map[string]interface{}{}
	g.Send("register_mobile", params, func(msg Message) {

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

	//fmt.Printf("message(%v)\n", string(b))

	msg, err := NewMessage(b)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if msg.Command == "mobile_gate_proxy" {
		g.RequestFromMobileProxy(msg)
		return
	}

	for command, f := range g.subscribers {
		if msg.Id == command {
			f(msg)
			g.UnSubscribe(msg.Id)
		}
	}
}

func (g *GateClient) onConnected() {
	g.registerServer()
	//g.registerMobile()
}

func (g *GateClient) onClosed() {

}

func (g *GateClient) Subscribe(id uuid.UUID, f func(msg Message)) {
	if g.subscribers[id] != nil {
		delete(g.subscribers, id)
	}
	g.subscribers[id] = f
}

func (g *GateClient) UnSubscribe(id uuid.UUID) {
	if g.subscribers[id] != nil {
		delete(g.subscribers, id)
	}
}

func (g *GateClient) Send(command string, payload map[string]interface{}, f func(msg Message)) (err error) {

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
	<- done

	return
}

func (g *GateClient) Status() string {
	return g.wsClient.status
}

func (g *GateClient) GetSettings() (*Settings, error) {
	return g.settings, nil
}

func (g *GateClient) UpdateSettings(settings *Settings) (err error) {
	g.settings = settings
	if err = g.SaveSettings(); err != nil {
		return
	}
	if !g.settings.Enabled {
		g.Close()
	}
	return
}

func (g *GateClient) GetMobileList() (list *MobileList, err error) {

	list = &MobileList{
		TokenList: make([]string, 0),
	}

	payload := map[string]interface{}{}
	g.Send("mobile_token_list", payload, func(msg Message) {
		if err = msg.IsError(); err != nil {
			return
		}
		if err = common.Copy(&list, msg.Payload, common.JsonEngine); err != nil {
			return
		}
	})

	return
}

func (g *GateClient) DeleteMobile(token string) (list *MobileList, err error) {

	payload := map[string]interface{}{
		"token": token,
	}
	g.Send("remove_mobile", payload, func(msg Message) {
		err = msg.IsError()
	})

	return
}

func (g *GateClient) AddMobile() (list *MobileList, err error) {

	payload := map[string]interface{}{}
	g.Send("register_mobile", payload, func(msg Message) {
		err = msg.IsError()
	})

	return
}

func (g *GateClient) RequestFromMobileProxy(message Message) {

	if g.wsClient.status != GateStatusConnected {
		return
	}

	debug.Println(message)

	response := Message{
		Id:      uuid.NewV4(),
		Command: message.Id.String(),
		Payload: map[string]interface{}{},
	}

	msg, _ := json.Marshal(response)
	if err := g.wsClient.write(websocket.TextMessage, msg); err != nil {
		log.Error(err.Error())
	}
}