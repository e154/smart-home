---
weight: 38
title: notifr
groups:
    - javascript
---

## IC.Notifr() {#ic_notifr}

Отправка уведомлений

### .newSMS() {#ic_notifr_new_sms}

смс сообщение

```coffeescript
msg = IC.Notifr().newSMS()
msg.AddPhone("+1-222-333-44-55")
msg.Text = "hola"

# optional
#tpl = IC.Template()
#render = tpl.render(name, {'key':'val'})
#msg.setRender(render)

IC.Notifr().send(msg)
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, смс объект [sms](#sms)

### .newEmail() {#ic_notifr_new_email}

Email сообщение

```coffeescript
msg = IC.Notifr().newEmail()
msg.From = ""
msg.To = ""
msg.Subject = ""
msg.Body = ""

# optional
#tpl = IC.Template()
#render = tpl.render(name, {'key':'val'})
#msg.setRender(render)

IC.Notifr().send(msg)
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, email объект [email](#email)


### .newSlack(channel, text) {#ic_notifr_new_slack}

сообщение в slack чат

```coffeescript
msg = IC.Notifr().newSlack(@main, 'hola')

# optional
#tpl = IC.Template()
#render = tpl.render(name, {'key':'val'})
#msg.setRender(render)

IC.Notifr().send(msg)
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, slack объект [slack](#slack)

### .newTelegram(text) {#ic_notifr_new_telegram}

сообщение в телеграм канал

```coffeescript
msg = IC.Notifr().newTelegram(text)

# optional
#tpl = IC.Template()
#render = tpl.render(name, {'key':'val'})
#msg.setRender(render)

IC.Notifr().send(msg)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `text`      | type: string, текст сообщения

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, telegram объект [telegram](#telegram)

### .send(msg) {#ic_notifr_send}

сообщение в телеграм канал

```coffeescript
IC.Notifr().send(msg)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `msg`      | type: Object, объект сообщения (sms,email,slack,telegram)
