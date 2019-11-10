package telegram

type TelegramConfig struct {
	Token  string `json:"token"`
	ChatId *int64 `json:"chat_id"`
}

func NewTelegramConfig(token string, chatId *int64) *TelegramConfig {
	return &TelegramConfig{
		Token:  token,
		ChatId: chatId,
	}
}
