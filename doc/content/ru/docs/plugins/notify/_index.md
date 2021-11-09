
---
title: "Notify"
linkTitle: "notify"
date: 2021-10-20
description: >

---

Базовое расширение для обеспечения уведомлений, таких как _sms_, _push_, _email_, _telegram_, etc.




### javascript свойства
----------------

### Новое сообщение

создает объект сообщения

```coffeescript
  msg = notifr.newMessage();
```
|  значение  | описание  |
|-------------|---------|
| newMessage() |    метод   |
| msg |   type: Object [Message](#объект-message)  |


----------------

### отправка сообщения


```coffeescript
  notifr.send(msg);
```
|  значение  | описание  |
|-------------|---------|
| send() |    метод   |
| msg |   type: Object (Message)  |


### объект Message

```coffeescript
  message = {
  from: "",
  type: "",
  attributes: {
    "key": "vlue"
  }
}
``` 

|  значение  | описание  |
|-------------|---------|
| from |    type: string   |
| type |    type: string   |
| attributes |   type: map[string]string  |
