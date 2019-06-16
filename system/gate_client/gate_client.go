package gate_client

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/uuid"
	"github.com/op/go-logging"
)

const (
	gateVarName = "gateClientParams"
)

var (
	log = logging.MustGetLogger("gate")
)

type GateClient struct {
	adaptors *adaptors.Adaptors
	settings *Settings
	wsClient *WsClient
}

func NewGateClient(adaptors *adaptors.Adaptors) (gate *GateClient) {
	gate = &GateClient{
		adaptors: adaptors,
		settings: &Settings{
			Id: uuid.NewV4(),
		},
		wsClient: NewWsClient(adaptors),
	}

	if err := gate.LoadSettings(); err != nil {
		log.Error(err.Error())
	}

	return
}

func (g *GateClient) Connect() {

	if !g.settings.Valid() {
		return
	}

	log.Info("Connect")

	g.wsClient.Connect(g.settings)

}

func (g *GateClient) GetToken() (token string, err error) {
	token, err = g.wsClient.GetToken()
	return
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
