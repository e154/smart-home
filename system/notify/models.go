package notify

import "github.com/e154/smart-home/models"

type NotifyConfig struct {
	MbAccessKey   string `json:"mb_access_key"`
	MbName        string `json:"mb_name"`
	TWFrom        string `json:"tw_from"`
	TWSid         string `json:"tw_sid"`
	TWAuthToken   string `json:"tw_auth_token"`
	TelegramToken string `json:"telegram_token"`
	EmailAuth     string `json:"email_auth"`
	EmailPass     string `json:"email_pass"`
	EmailSmtp     string `json:"email_smtp"`
	EmailPort     int    `json:"email_port"`
	EmailSender   string `json:"email_sender"`
	SlackToken    string `json:"slack_token"`
	SlackUserName string `json:"slack_user_name"`
}

type NotifyStat struct {
	MbBalance float32 `json:"mb_balance,omitempty"`
	TwBalance float32 `json:"tw_balance,omitempty"`
	Workers   int     `json:"workers"`
}

type IMessage interface {
	Save() (addresses []string, message *models.Message)
}
