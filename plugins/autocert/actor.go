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

package autocert

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	actionPool chan events.EventCallEntityAction
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		actionPool: make(chan events.EventCallEntityAction, 1000),
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return
}

func (e *Actor) Destroy() {

}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {
	switch msg.ActionName {
	case "RequestCertificate":
		e.requestCertificate()
	}
}

func (e *Actor) requestCertificate() {

	var domains = strings.Split(strings.TrimSpace(e.Setts[AttrDomains].String()), " ")
	if len(domains) == 0 {
		log.Error("domain list is nil")
		return
	}

	var ownerEmail = strings.Split(strings.TrimSpace(e.Setts[AttrEmails].String()), " ")
	if len(ownerEmail) == 0 {
		log.Error("email list is nil")
		return
	}
	for i, email := range ownerEmail {
		if !strings.Contains(email, "mailto:") {
			ownerEmail[i] = fmt.Sprintf("mailto:%s", email)
		}
	}

	var cloudflareAPIToken = ""
	if val, ok := e.Setts[AttrCloudflareAPIToken]; ok {
		cloudflareAPIToken = strings.TrimSpace(val.Decrypt())
	}

	if cloudflareAPIToken == "" {
		log.Error("Cloudflare API Token is nil")
		return
	}

	prod := e.Setts[AttrProduction].Bool()

	autocert := NewAutocert(domains, ownerEmail, cloudflareAPIToken, prod)
	err := autocert.RequestCertificate(context.Background())
	if err != nil {
		log.Error(err.Error())
		e.SetActorState(common.String(StateError))
		e.SaveState(false, true)
		return
	}

	//fmt.Println("private", string(autocert.PrivateKey()))
	//fmt.Println("public", string(autocert.PublicKey()))

	var attributeValues = make(m.AttributeValue)
	attributeValues[AttrPrivateKey] = string(autocert.PrivateKey())
	attributeValues[AttrPublicKey] = string(autocert.PublicKey())
	e.DeserializeAttr(attributeValues)

	e.SetActorState(common.String(StateSuccessfully))
	e.SaveState(false, true)
}

func (e *Actor) checkCertificate() {

	log.Infof("check ssl certificate %s", e.Id)

	// get last state
	state := e.GetCurrentState()

	// if last state is nil
	if state == nil || state.LastUpdated == nil {
		e.requestCertificate()
		return
	}

	// if last state error
	if state.State.Name == StateError {
		e.requestCertificate()
		return
	}

	// if more than 85 days have passed
	now := time.Now()
	if now.Sub(*state.LastUpdated).Hours()/24 > 85 {
		e.requestCertificate()
		return
	}
}
