
---
title: "HTML5 Уведомления"
linkTitle: "html5_notify"
date: 2021-10-20
description: >

---

HTML5 уведомление - это способ показать пользователю информацию на веб-странице без перезагрузки или перехода на другую страницу. Уведомления могут содержать текст, изображения, звук или видео и появляются в правом нижнем углу экрана. Уведомления на HTML5 используют API, который позволяет веб-приложениям отправлять уведомления браузеру с помощью JavaScript

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
msg.type = 'html5_notify';
msg.attributes = {
  'userIDS': '14',
  'title': 'Lorem Ipsum is simply',
  'body': 'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum'
};

```
|  значение  | описание  |
|-------------|---------|
| newMessage() |    метод   |
| msg |   type: Object (Message)  |

----------------

### пример кода

```coffeescript
# html5_notify
# ##################################

sendMsg =()->
  msg = notifr.newMessage();
  msg.type = 'html5_notify';
  msg.attributes = {
    'userIDS': '14',
    'title': 'Lorem Ipsum is simply',
    'body': 'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum'
  };
  notifr.send(msg);

sendMsg()
```

