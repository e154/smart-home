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
	"github.com/e154/smart-home/system/cache"
	"reflect"
	"sync"
	"time"
)

var (
	log = common.MustGetLogger("mqtt_authenticator")
)

// ErrBadLoginOrPassword ...
var ErrBadLoginOrPassword = fmt.Errorf("bad login or password")

// ErrPrincipalDisabled ...
var ErrPrincipalDisabled = fmt.Errorf("principal disabled")

// MqttAuthenticator ...
type MqttAuthenticator interface {
	Authenticate(login string, pass interface{}) (err error)
	Register(fn func(login, password string) (err error)) (err error)
	Unregister(fn func(login, password string) (err error)) (err error)
}

// Authenticator ...
type Authenticator struct {
	adaptors  *adaptors.Adaptors
	cache     cache.Cache
	handlerMu *sync.Mutex
	handlers  []reflect.Value
}

// NewAuthenticator ...
func NewAuthenticator(adaptors *adaptors.Adaptors) MqttAuthenticator {
	bm, _ := cache.NewCache("memory", fmt.Sprintf(`{"interval":%d}`, time.Second*60))
	return &Authenticator{
		adaptors:  adaptors,
		cache:     bm,
		handlerMu: &sync.Mutex{},
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

	for _, v := range a.handlers {
		result := v.Call([]reflect.Value{reflect.ValueOf(login), reflect.ValueOf(pass)})
		if result[0].Interface() != nil {
			err = result[0].Interface().(error)
		} else {
			continue
		}
		break
	}

	a.cache.Put(login, pass, 60*time.Second)

	return
}

// Register ...
func (a *Authenticator) Register(fn func(login, password string) (err error)) (err error) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		err = fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	a.handlerMu.Lock()
	defer a.handlerMu.Unlock()

	rv := reflect.ValueOf(fn)

	for _, v := range a.handlers {
		if v == rv {
			return
		}
	}

	a.handlers = append(a.handlers, rv)

	log.Infof("register ...")

	return
}

// Unregister ...
func (a *Authenticator) Unregister(fn func(login, password string) (err error)) (err error) {
	a.handlerMu.Lock()
	defer a.handlerMu.Unlock()

	rv := reflect.ValueOf(fn)

	for i, v := range a.handlers {
		if v == rv {
			a.handlers = append(a.handlers[:i], a.handlers[i+1:]...)
		}
	}

	log.Infof("unregister ...")

	return
}
