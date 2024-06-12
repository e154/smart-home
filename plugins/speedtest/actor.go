// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package speedtest

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/showwin/speedtest-go/speedtest"
	"go.uber.org/atomic"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	isStarted *atomic.Bool
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	actor := &Actor{
		BaseActor: supervisor.NewBaseActor(entity, service),
		isStarted: atomic.NewBool(false),
	}

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	if actor.States == nil {
		actor.States = NewStates()
	}

	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) runTest() {
	if e.isStarted.Load() {
		return
	}
	e.isStarted.Store(true)
	defer e.isStarted.Store(false)

	var speedtestClient = speedtest.New()

	// Use a proxy for the speedtest. eg: socks://127.0.0.1:7890
	// speedtest.WithUserConfig(&speedtest.UserConfig{Proxy: "socks://127.0.0.1:7890"})(speedtestClient)

	// Select a network card as the data interface.
	// speedtest.WithUserConfig(&speedtest.UserConfig{Source: "192.168.1.101"})(speedtestClient)

	// Get user's network information
	// user, _ := speedtestClient.FetchUserInfo()

	// Get a list of servers near a specified location
	// user.SetLocationByCity("Tokyo")
	// user.SetLocation("Osaka", 34.6952, 135.5006)

	// Search server using serverID.
	// eg: fetch server with ID 28910.
	// speedtest.ErrServerNotFound will be returned if the server cannot be found.
	// server, err := speedtest.FetchServerByID("28910")

	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	e.updateState(nil)

	for _, s := range targets {
		// Please make sure your host can access this test server,
		// otherwise you will get an error.
		// It is recommended to replace a server at this time
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()
		e.updateState(s)
		s.Context.Reset() // reset counter
	}

}

func (e *Actor) updateState(state *speedtest.Server) {

	var attributeValues = make(m.AttributeValue)
	if state != nil {
		attributeValues[AttrDLSpeed] = state.DLSpeed.Mbps()
		attributeValues[AttrULSpeed] = state.ULSpeed.Mbps()
		attributeValues[AttrLatency] = state.Latency.String()
		attributeValues[AttrName] = state.Name
		attributeValues[AttrCountry] = state.Country
		attributeValues[AttrSponsor] = state.Sponsor
		attributeValues[AttrDistance] = state.Distance
		attributeValues[AttrJitter] = state.Jitter.String()
		attributeValues[AttrPoint] = []interface{}{state.Lon, state.Lat}

		s := e.States[StateCompleted]
		e.SetActorState(common.String(s.Name))
	} else {
		s := e.States[StateInProcess]
		e.SetActorState(common.String(s.Name))
	}

	e.DeserializeAttr(attributeValues)

	e.SaveState(state == nil, state != nil)
}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	if event.ActionName != ActionCheck {
		return
	}

	go e.runTest()
}
