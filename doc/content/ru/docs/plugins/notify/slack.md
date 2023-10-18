
---
title: "Slack"
linkTitle: "slack"
date: 2021-10-20
description: >

---

Отправка уведомлений в чат slack

### Настройка
* Token
* User name

### javascript свойства

----------------

### Новое сообщение

создает объект сообщения

```coffeescript
msg = notifr.newMessage();
msg.entity_id = 'slack.name';
msg.attributes = {
  'channel': '#ch',
  'text': 'some text'
};

```
|  значение  | описание  |
|-------------|---------|
| newMessage() |    метод   |
| msg |   type: Object (Message)  |

----------------

### пример кода

```coffeescript
# slack
# ##################################

sendMsg =(body)->
  msg = notifr.newMessage();
    msg.entity_id = 'slack.name';
    msg.attributes = {
      'channel': '#ch',
      'text': 'some text'
  };
  notifr.send(msg);
```

