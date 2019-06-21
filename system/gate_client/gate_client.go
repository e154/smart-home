package gate_client

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/uuid"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"reflect"
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
	subscribers map[string]func(payload interface{})
}

func NewGateClient(adaptors *adaptors.Adaptors,
	graceful *graceful_service.GracefulService, ) (gate *GateClient) {
	gate = &GateClient{
		adaptors:    adaptors,
		settings:    &Settings{},
		subscribers: make(map[string]func(payload interface{})),
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

	params := map[string]interface{}{}
	g.Send("register_server", params, func(payload interface{}) {

		v, ok := reflect.ValueOf(payload).Interface().(map[string]interface{})
		if !ok {
			log.Error("bad reflect casting")
			return
		}

		g.settings.GateServerToken = v["token"].(string)

		_ = g.SaveSettings()
	})
}

func (g *GateClient) registerMobile() {

	params := map[string]interface{}{}
	g.Send("register_mobile", params, func(payload interface{}) {

		v, ok := reflect.ValueOf(payload).Interface().(map[string]interface{})
		if !ok {
			log.Error("bad reflect casting")
			return
		}

		fmt.Println("mobile token ", v["token"])
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

func (g *GateClient) onMessage(message []byte) {

	fmt.Println(string(message))

	re := map[string]interface{}{}
	if err := json.Unmarshal(message, &re); err != nil {
		log.Error(err.Error())
		return
	}

	var id string
	for key, v := range re {

		if key == "id" {
			id = v.(string)

			for command, f := range g.subscribers {
				if id == command {
					f(re)
					g.UnSubscribe(command)
				}
			}
		}
	}
}

func (g *GateClient) onConnected() {
	g.registerServer()
	//g.registerMobile()
}

func (g *GateClient) onClosed() {

}

func (g *GateClient) Subscribe(command string, f func(payload interface{})) {
	if g.subscribers[command] != nil {
		delete(g.subscribers, command)
	}
	g.subscribers[command] = f
}

func (g *GateClient) UnSubscribe(command string) {
	if g.subscribers[command] != nil {
		delete(g.subscribers, command)
	}
}

func (g *GateClient) Send(command string, params map[string]interface{}, f func(payload interface{})) {

	messageId := uuid.NewV4().String()
	g.Subscribe(messageId, f)

	params["id"] = messageId

	payload, _ := json.Marshal(&map[string]interface{}{command: params})
	if err := g.wsClient.write(websocket.TextMessage, payload); err != nil {
		log.Error(err.Error())
	}
}
