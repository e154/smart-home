// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"github.com/e154/smart-home/system/cache"
	"time"
)

var (
	log = common.MustGetLogger("mqtt_authenticator")
)

// ErrBadLoginOrPassword ...
var ErrBadLoginOrPassword = fmt.Errorf("bad login or password")

// ErrPrincipalDisabled ...
var ErrPrincipalDisabled = fmt.Errorf("principal disabled")

type MqttAuthenticator interface {
	Authenticate(login string, pass interface{}) (err error)
}

// Authenticator ...
type Authenticator struct {
	adaptors *adaptors.Adaptors
	cache    cache.Cache
}

// NewAuthenticator ...
func NewAuthenticator(adaptors *adaptors.Adaptors) MqttAuthenticator {
	bm, _ := cache.NewCache("memory", fmt.Sprintf(`{"interval":%d}`, time.Second*60))
	return &Authenticator{
		adaptors: adaptors,
		cache:    bm,
	}
}

// Authenticate ...
func (a *Authenticator) Authenticate(login string, pass interface{}) (err error) {

	log.Infof("login: \"%v\", pass: \"%v\"", login, pass)

	password, ok := pass.(string)
	if !ok || password == "" {
		err = ErrBadLoginOrPassword
	}

	if a.cache.IsExist(login) {
		if password == a.cache.Get(login) {
			return
		}
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

	if bridge.EncryptedPassword == "" && password == "" {
		return
	}

	if ok := common.CheckPasswordHash(password, bridge.EncryptedPassword); !ok {
		err = ErrBadLoginOrPassword
		return
	}

	a.cache.Put(login, password, 60*time.Second)

	return
}
