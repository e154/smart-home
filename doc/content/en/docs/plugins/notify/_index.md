
---
title: "Отправка сообщений"
linkTitle: "notify"
date: 2021-10-20
description: >

---

The system has a plugin called "notifr" that provides notification capabilities. It supports methods for sending notifications through various channels, including SMS, push notifications, email, Telegram, and HTML5 notifications.

Here are the supported notification methods provided by the "notifr" plugin:

1. SMS: Allows sending SMS notifications.
2. Push: Enables sending push notifications to mobile devices.
3. Email: Supports sending email notifications.
4. Telegram: Provides the ability to send notifications through Telegram messenger.
5. HTML5 notifications: Allows displaying notifications using HTML5 capabilities.

With the "notifr" plugin, you can utilize these methods to send notifications to users through different channels based on their preferences or the nature of the notification.

Please note that the specific implementation and usage details of each method may vary based on the configuration and setup of the "notifr" plugin in your system.

### javascript свойства
----------------

### New message

Create message object

```coffeescript
  msg = notifr.newMessage();
```
|  Property  | Description  |
|-------------|---------|
| newMessage() |    метод   |
| msg |   type: Object [Message](#объект-message)  |


----------------

### Send message


```coffeescript
  notifr.send(msg);
```
|  Property  | Description  |
|-------------|---------|
| send() |    метод   |
| msg |   type: Object (Message)  |


### Object Message

```coffeescript
  message = {
  from: "",
  type: "",
  attributes: {
    "key": "vlue"
  }
}
``` 

|  Property  | Description  |
|-------------|---------|
| from |    type: string   |
| type |    type: string   |
| attributes |   type: map[string]string  |
