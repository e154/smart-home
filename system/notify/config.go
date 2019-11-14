package notify

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

const (
	notifyVarName = "notify"
)

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

func NewNotifyConfig(adaptor *adaptors.Adaptors) *NotifyConfig {
	return &NotifyConfig{
		adaptor: adaptor,
	}
}

func (n *NotifyConfig) Get() {

	v, err := n.adaptor.Variable.GetByName(notifyVarName)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if err = json.Unmarshal([]byte(v.Value), n); err != nil {
		log.Error(err.Error())
	}
}

func (n *NotifyConfig) Update() (err error) {

	log.Infof("update settings")

	var b []byte
	if b, err = json.Marshal(n); err != nil {
		log.Error(err.Error())
		return
	}

	variable := &m.Variable{
		Name:     notifyVarName,
		Value:    string(b),
		Autoload: false,
	}

	if err = n.adaptor.Variable.Update(variable); err != nil {
		log.Error(err.Error())
	}

	return
}
