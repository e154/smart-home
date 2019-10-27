package telegram

type TelegramConfig struct {
	Token string `json:"token"`
}

func NewTelegramConfig(token string) *TelegramConfig {
	return &TelegramConfig{
		Token: token,
	}
}
