---
title: "Push уведомления"
linkTitle: "webpush"
date: 2022-06-08
description: >

---

Web push — это короткие сообщения, которые smart home отправляет в браузер, чтобы сообщить о какой-либо важной
информации. Они показываются даже тем, кто ушел с сайта, и не требуют указания почты или телефона.

### Настройка

* public key
* private key

### Атрибуты

* userIDS
* title
* body

### javascript свойства

----------------

### Новое сообщение

создает объект сообщения

```coffeescript
msg = notifr.newMessage();
msg.type = 'webpush';
msg.attributes = {
  'userIDS': '14',
  'title': 'Lorem Ipsum is simply',
  'body': 'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum'
};

```

| значение     | описание               |
|--------------|------------------------|
| newMessage() | метод                  |
| msg          | type: Object (Message) |

----------------

### пример кода

```coffeescript
# webpush
# ##################################
sendMsg = ()->
  msg = notifr.newMessage();
  msg.type = 'webpush';
  msg.attributes = {
    'userIDS': '14',
    'title': 'Lorem Ipsum is simply',
    'body': 'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum'
  };
  notifr.send(msg);

sendMsg()
```

