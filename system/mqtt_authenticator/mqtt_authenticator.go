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

package mqtt_authenticator

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/uuid"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("mqtt_authenticator")
)

var ErrBadLoginOrPassword = fmt.Errorf("bad login or password")
var ErrPrincipalDisabled = fmt.Errorf("principal disabled")

type Authenticator struct {
	adaptors *adaptors.Adaptors
	name     string
	login    string
	password string
}

func NewAuthenticator(adaptors *adaptors.Adaptors) *Authenticator {
	return &Authenticator{
		adaptors: adaptors,
		name:     "base",
		login:    "local",
		password: uuid.NewV4().String(),
	}
}

func (a *Authenticator) Authenticate(login string, pass interface{}) (err error) {

	log.Infof("login: %v, pass: %v", login, pass)

	password, ok := pass.(string)
	if !ok || password == "" {
		err = ErrBadLoginOrPassword
	}

	if login == a.login && pass == a.password {
		return
	}

	// nodes
	if err = a.checkNode(login, password); err == nil {
		return
	}

	// zigbee2mqtt
	if err = a.checkZigbee2matt(login, password); err == nil {
		return
	}

	return
}

func (a Authenticator) checkZigbee2matt(login, password string) (err error) {

	var bridge *m.Zigbee2mqtt
	if bridge, err = a.adaptors.Zigbee2mqtt.GetByLogin(login); err != nil {
		return
	}

	//if bridge.Status == "disabled" {
	//	err = ErrPrincipalDisabled
	//	return
	//}

	if ok := common.CheckPasswordHash(password, bridge.EncryptedPassword); !ok {
		err = ErrBadLoginOrPassword
	}

	return
}

func (a Authenticator) checkNode(login, password string) (err error) {

	var node *m.Node
	if node, err = a.adaptors.Node.GetByLogin(login); err != nil {
		return
	}

	if node.Status == "disabled" {
		err = ErrPrincipalDisabled
		return
	}

	if ok := common.CheckPasswordHash(password, node.EncryptedPassword); !ok {
		err = ErrBadLoginOrPassword
	}

	return
}

func (a Authenticator) Name() string {
	return a.name
}

func (a Authenticator) Password() string {
	return a.password
}

func (a Authenticator) Login() string {
	return a.login
}
