package notify

type SMS struct {
	Phone    string            `json:"phone"`
	Text     string            `json:"text"`
	Template string            `json:"template"`
	Data     map[string]string `json:"data"`
}

func NewSMS() (sms *SMS) {
	return &SMS{
		Data: make(map[string]string),
	}
}
