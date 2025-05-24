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

package mqtt

import (
	"context"
	"errors"

	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/retained"
	"github.com/DrmagicE/gmqtt/server"
)

var (
	// ErrInvalidTopicFilter ...
	ErrInvalidTopicFilter = errors.New("invalid topic filter")
	// ErrInvalidQos ...
	ErrInvalidQos = errors.New("invalid Qos")
	// ErrInvalidUtf8String ...
	ErrInvalidUtf8String = errors.New("invalid utf-8 string")
)

// MqttCli ...
type MqttCli interface {
	Publish(topic string, payload []byte) error
	Subscribe(topic string, handler MessageHandler) error
	Unsubscribe(topic string)
	UnsubscribeAll()
	OnMsgArrived(ctx context.Context, client server.Client, req *server.MsgArrivedRequest)
}

// MqttAuthenticator ...
type MqttAuthenticator interface {
	Authenticate(login string, pass interface{}) (err error)
	//DEPRECATED
	Register(fn func(login, password string) (err error)) (err error)
	//DEPRECATED
	Unregister(fn func(login, password string) (err error)) (err error)
}

// MqttServ ...
type MqttServ interface {
	Shutdown() error
	Start()
	Publish(topic string, payload []byte, qos uint8, retain bool) error
	NewClient(name string) MqttCli
	RemoveClient(name string)
	Authenticator() MqttAuthenticator
}

// GMqttServer ...
type GMqttServer interface {
	Run() error
	Stop(ctx context.Context) error
	Init(opts ...server.Options) error
	// SubscriptionStore returns the subscription.Store.
	SubscriptionService() server.SubscriptionService
	// RetainedStore returns the retained.Store.
	RetainedStore() retained.Store
	// Publisher returns the Publisher
	Publisher() server.Publisher
	// client return the ClientService
	ClientService() server.ClientService
	// GetConfig returns the config of the server
	GetConfig() config.Config
	// StatsManager returns StatsReader
	StatsManager() server.StatsReader
}

// Message ...
type Message struct {
	Dup      bool
	Qos      uint8
	Retained bool
	Topic    string
	PacketID uint16
	Payload  []byte
}

// MessageHandler ...
type MessageHandler func(MqttCli, Message)
