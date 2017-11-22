package telemetry_map

import (
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/stream"
)

func NewMap() *Map {
	return &Map{}
}

type Map struct {
	DeviceStates	map[int64]*models.DeviceState	`json:"device_states"`
}

func (m *Map) Run() {

}

func (m *Map) Stop() {

}

func (m *Map) Broadcast(pack string) {

}

func (m *Map) BroadcastOne(pack string, id int64) {

}

// only on request: 'get.map.states'
//
func streamGetMapStates(client *stream.Client, value interface{}) {


}