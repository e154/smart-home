package notify

// Javascript Binding
//
// IC.Notifr()
//	 .newSMS()
//	 .newEmail()
//	 .newSlack(channel, text)
//	 .newTelegram(text)
//	 .send(msg)
//
type NotifyBind struct {
	notify *Notify
}

func (b *NotifyBind) NewSMS() *SMS {
	return NewSMS()
}

func (b *NotifyBind) NewEmail() *Email {
	return NewEmail()
}

func (b *NotifyBind) NewSlack(channel, text string) *SlackMessage {
	return NewSlackMessage(channel, text)
}

func (b *NotifyBind) NewTelegram(text string) *Telegram {
	return NewTelegram(text)
}

func (b *NotifyBind) Send(msg interface{}) {
	b.notify.Send(msg)
}
