package notify

type Telegram struct {
	Text    string   `json:"text"`
	Channel []string `json:"channel"`
}

func NewTelegram(text string, channel []string) *Telegram {
	return &Telegram{
		Text:    text,
		Channel: channel,
	}
}
