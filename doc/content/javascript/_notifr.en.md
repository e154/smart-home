---
weight: 38
title: notifr
groups:
    - javascript
---

## Notifr {#notifr}

Отправка уведомлений

### .NewSMS() {#notifr_new_sms}

смс сообщение

```coffeescript
msg = Notifr.NewSMS()
msg.AddPhone("+1-222-333-44-55")
msg.Text = "hola"

# optional
#render = Template.Render(name, {'key':'val'})
#msg.SetRender(render)

Notifr.Send(msg)
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, смс объект [sms](#sms)

### .NewEmail() {#notifr_new_email}

Email сообщение

```coffeescript
msg = Notifr.NewEmail()
msg.From = ""
msg.To = ""
msg.Subject = ""
msg.Body = ""

# optional
#render = Template.Render(name, {'key':'val'})
#msg.SetRender(render)

Notifr.Send(msg)
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, email объект [email](#email)


### .NewSlack(channel, text) {#notifr_new_slack}

сообщение в slack чат

```coffeescript
msg = Notifr.NewSlack(@main, 'hola')

# optional
#render = Template.Render(name, {'key':'val'})
#msg.SetRender(render)

Notifr.Send(msg)
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, slack объект [slack](#slack)

### .NewTelegram(text) {#notifr_new_telegram}

сообщение в телеграм канал

```coffeescript
msg = Notifr.NewTelegram(text)

# optional
#render = Template.Render(name, {'key':'val'})
#msg.SetRender(render)

Notifr.Send(msg)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `text`      | type: string, текст сообщения

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, telegram объект [telegram](#telegram)

### .Send(msg) {#notifr_send}

сообщение в телеграм канал

```coffeescript
Notifr.Send(msg)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, объект сообщения (sms,email,slack,telegram)
