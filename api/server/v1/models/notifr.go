package models

// swagger:model
type UpdateNotifrConfig struct {
	MbAccessKey    string `json:"mb_access_key"`
	MbName         string `json:"mb_name"`
	TWFrom         string `json:"tw_from"`
	TWSid          string `json:"tw_sid"`
	TWAuthToken    string `json:"tw_auth_token"`
	TelegramToken  string `json:"telegram_token"`
	TelegramChatId int64  `json:"telegram_chat_id"`
	EmailAuth      string `json:"email_auth"`
	EmailPass      string `json:"email_pass"`
	EmailSmtp      string `json:"email_smtp"`
	EmailPort      int    `json:"email_port"`
	EmailSender    string `json:"email_sender"`
	SlackToken     string `json:"slack_token"`
	SlackUserName  string `json:"slack_user_name"`
}

// swagger:model
type NotifrConfig struct {
	MbAccessKey    string `json:"mb_access_key"`
	MbName         string `json:"mb_name"`
	TWFrom         string `json:"tw_from"`
	TWSid          string `json:"tw_sid"`
	TWAuthToken    string `json:"tw_auth_token"`
	TelegramToken  string `json:"telegram_token"`
	TelegramChatId int64  `json:"telegram_chat_id"`
	EmailAuth      string `json:"email_auth"`
	EmailPass      string `json:"email_pass"`
	EmailSmtp      string `json:"email_smtp"`
	EmailPort      int    `json:"email_port"`
	EmailSender    string `json:"email_sender"`
	SlackToken     string `json:"slack_token"`
	SlackUserName  string `json:"slack_user_name"`
}

// swagger:model
type NewNotifrMessage struct {
	Type         string            `json:"type"`
	BodyType     string            `json:"body_type"`
	EmailFrom    *string           `json:"email_from"`
	EmailSubject *string           `json:"email_subject"`
	EmailBody    *string           `json:"email_body"`
	Template     *string           `json:"template"`
	SmsText      *string           `json:"sms_text"`
	SlackText    *string           `json:"slack_text"`
	Params       map[string]string `json:"params"`
	Address      string            `json:"address"`
}
