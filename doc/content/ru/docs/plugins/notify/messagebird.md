
---
title: "Messagebird"
linkTitle: "messagebird"
date: 2021-10-20
description: >

---

Отправка sms уведомлений через провайдера услуг оповещения [messagebird](https://messagebird.com). Для возможности 
отправки сообщений потребуется регистрация в [messagebird](https://messagebird.com), и положительный баланс.

### Настройка
* Access Key
* Name

### Атрибуты
* Payment
* Type
* Amount

### javascript свойства

----------------

### Новое сообщение

создает объект сообщения

```coffeescript
msg = notifr.newMessage();
msg.entity_id = 'messagebird.name';
msg.attributes = {
  'phone': '+79990000001',
  'body': 'some text'
};

```
|  значение  | описание  |
|-------------|---------|
| newMessage() |    метод   |
| msg |   type: Object (Message)  |

----------------

### пример кода

```coffeescript
# messagebird
# ##################################

sendMsg =(body)->
  msg = notifr.newMessage();
  msg.entity_id = 'messagebird.name';
  msg.attributes = {
    'phone': '+79990000001',
    'body': 'some text'
  };
  notifr.send(msg);
```

