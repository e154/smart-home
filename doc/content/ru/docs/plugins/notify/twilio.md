---
title: "Twilio"
linkTitle: "twilio"
date: 2021-10-20
description: >

---

Отправка sms уведомлений через провайдера услуг оповещения [twilio](https://www.twilio.com/messaging). Для возможности
отправки сообщений потребуется регистрация в [twilio](https://www.twilio.com/messaging), и положительный баланс.

### Настройка

* Token
* Sid
* From

### Атрибуты

* Amount
* Sid
* Currency

### javascript свойства

----------------

### Новое сообщение

создает объект сообщения

```coffeescript
msg = notifr.newMessage();
msg.entity_id = 'twilio.name';
msg.attributes = {
  'phone': '+79990000001',
  'body': 'some text'
};

```

| значение     | описание               |
|--------------|------------------------|
| newMessage() | метод                  |
| msg          | type: Object (Message) |

----------------

### пример кода

```coffeescript
# twilio
# ##################################

sendMsg = (body)->
  msg = notifr.newMessage();
  msg.entity_id = 'twilio.name';
  msg.attributes = {
    'phone': '+79990000001',
    'body': 'some text'
  };
  notifr.send(msg);
```

