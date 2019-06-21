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
		adaptors: adaptors,
		settings: &Settings{
			Id: uuid.NewV4(),
		},
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

func (g *GateClient) registerClient() {

	params := map[string]interface{}{"server_id": g.settings.Id}
	g.Send("register_server", params, func(payload interface{}) {
		fmt.Println(payload)
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
		return
	}

	for key, value := range re {

		switch key {

		default:
			for command, f := range g.subscribers {
				if key == command {
					f(value)
				}
			}
		}
	}
}

func (g *GateClient) onConnected() {
	g.registerClient()
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

	g.UnSubscribe(messageId)
}
