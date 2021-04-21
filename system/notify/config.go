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

package notify

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

const (
	notifyVarName = "notify"
)

// NotifyConfig ...
type NotifyConfig struct {
	adaptor        *adaptors.Adaptors
	MbAccessKey    string `json:"mb_access_key"`
	MbName         string `json:"mb_name"`
	TWFrom         string `json:"tw_from"`
	TWSid          string `json:"tw_sid"`
	TWAuthToken    string `json:"tw_auth_token"`
	TelegramToken  string `json:"telegram_token"`
	TelegramChatId *int64 `json:"telegram_chat_id"`
	EmailAuth      string `json:"email_auth"`
	EmailPass      string `json:"email_pass"`
	EmailSmtp      string `json:"email_smtp"`
	EmailPort      int    `json:"email_port"`
	EmailSender    string `json:"email_sender"`
	SlackToken     string `json:"slack_token"`
	SlackUserName  string `json:"slack_user_name"`
}

// NewNotifyConfig ...
func NewNotifyConfig(adaptor *adaptors.Adaptors) *NotifyConfig {
	return &NotifyConfig{
		adaptor: adaptor,
	}
}

// Get ...
func (n *NotifyConfig) Get() {

	v, err := n.adaptor.Variable.GetByName(notifyVarName)
	if err != nil {
		log.Warnf("variable with name '%s', error: %s", notifyVarName, err.Error())
		return
	}

	if err = json.Unmarshal([]byte(v.Value), n); err != nil {
		log.Error(err.Error())
	}
}

// Update ...
func (n *NotifyConfig) Update() (err error) {

	log.Infof("update settings")

	var b []byte
	if b, err = json.Marshal(n); err != nil {
		log.Error(err.Error())
		return
	}

	variable := m.Variable{
		Name:     notifyVarName,
		Value:    string(b),
		Autoload: false,
	}

	if err = n.adaptor.Variable.Update(variable); err != nil {
		log.Error(err.Error())
	}

	return
}
