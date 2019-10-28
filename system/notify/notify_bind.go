package notify

type NotifyBind struct {
	notify *Notify
}

func (b *NotifyBind) NewSms() *SMS {
	return NewSMS()
}

func (b *NotifyBind) NewEmail() *Email {
	return NewEmail()
}

func (b *NotifyBind) Send(msg interface{}) {
	b.notify.Send(msg)
}
