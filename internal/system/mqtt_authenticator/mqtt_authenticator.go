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

package mqtt_authenticator

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/e154/smart-home/internal/system/cache"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/logger"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/mqtt"
)

var (
	log = logger.MustGetLogger("mqtt_authenticator")
)

var _ mqtt.MqttAuthenticator = (*Authenticator)(nil)

// Authenticator ...
type Authenticator struct {
	adaptors  *adaptors.Adaptors
	cache     cache.Cache
	handlerMu *sync.Mutex
	handlers  []reflect.Value
}

// NewAuthenticator ...
func NewAuthenticator(adaptors *adaptors.Adaptors) mqtt.MqttAuthenticator {
	bm, _ := cache.NewCache("memory", `{"interval":60}`)
	return &Authenticator{
		adaptors:  adaptors,
		cache:     bm,
		handlerMu: &sync.Mutex{},
	}
}

// Authenticate ...
func (a *Authenticator) Authenticate(login string, pass interface{}) (err error) {

	password, ok := pass.(string)
	if !ok || password == "" {
		err = apperr.ErrBadLoginOrPassword
	}

	var value interface{}
	if value, err = a.cache.Get(context.Background(), login); value != nil {
		if password == pass.(string) {
			return
		}
	}

	defer func() {
		if err == nil {
			_ = a.cache.Put(context.Background(), login, pass, 60*time.Second)
		}
	}()

	for _, v := range a.handlers {
		result := v.Call([]reflect.Value{reflect.ValueOf(login), reflect.ValueOf(pass)})
		if result[0].Interface() != nil {
			if err, ok = result[0].Interface().(error); !ok {
				err = nil
				return
			}
		} else {
			err = nil
			return
		}
	}

	var user *m.User
	if user, err = a.adaptors.User.GetByNickname(context.Background(), login); err != nil {
		if user, err = a.adaptors.User.GetByEmail(context.Background(), login); err != nil {
			err = fmt.Errorf("%s: %w", fmt.Sprintf("email %s", login), apperr.ErrUnauthorized)
			log.Warnf("failed login attempt: \"%v\", pass: \"%v\"", login, pass)
			return
		}
	}

	if !user.CheckPass(password) {
		err = apperr.ErrPassNotValid
		return
	} else if user.Status == "blocked" {
		err = apperr.ErrAccountIsBlocked
		return
	}

	return
}

// Register ...
func (a *Authenticator) Register(fn func(login, password string) (err error)) (err error) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		err = fmt.Errorf("%s: %w", fmt.Sprintf("%s is not a reflect.Func", reflect.TypeOf(fn)), apperr.ErrInternal)
	}

	a.handlerMu.Lock()
	defer a.handlerMu.Unlock()

	rv := reflect.ValueOf(fn)

	for _, v := range a.handlers {
		if v == rv || v.Pointer() == rv.Pointer() {
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

	var indexesToDelete []int

	for i, v := range a.handlers {
		if v == rv || v.Pointer() == rv.Pointer() {
			indexesToDelete = append(indexesToDelete, i)
		}
	}

	for i := len(indexesToDelete) - 1; i >= 0; i-- {
		index := indexesToDelete[i]
		a.handlers = append(a.handlers[:index], a.handlers[index+1:]...)
	}

	log.Infof("unregister ...")

	return
}
