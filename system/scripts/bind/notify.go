package bind

import (
	"github.com/e154/smart-home/system/notify"
)

type NotifyBind struct {
	notify *notify.Notify
}

func NewNotifyBind(notify *notify.Notify) *NotifyBind {
	return &NotifyBind{notify: notify}
}

func (b *NotifyBind) NewSms() *notify.SMS {
	return notify.NewSMS()
}

func (b *NotifyBind) Send(msg interface{}) {
	b.notify.Send(msg)
}
