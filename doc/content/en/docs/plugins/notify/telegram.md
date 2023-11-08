
---
title: "Telegram"
linkTitle: "telegram"
date: 2021-10-20
description: >

---

To work with the Telegram API and handle interactive commands from clients or use event-based notifications, the system provides the following configuration:

### Configuration:
* Token

### Commands:
* Custom set of commands

### Attributes:
* Custom set of attributes

### Actions:

There are two reserved commands in the system:

* **/start** - Subscribe to notifications
* **/quit** - Unsubscribe from notifications

**Action** - It should be named in lowercase, without the "/" sign. A custom command will automatically be added to the list of available commands in the **/help** section. The custom command should be called in uppercase **/ACTION** -> **/action**.

### JavaScript Properties:

----------------

### New Message:

Creates a message object.

```coffeescript
msg = notifr.newMessage();
msg.entity_id = 'telegram.name';
msg.attributes = {
  'body': 'some text msg',
  'chat_id': 123456,
  'keys': ['foo', 'bar']
};

```
|  значение  | описание  |
|-------------|---------|
| newMessage() |    метод   |
| msg |   type: Object (Message)  |

attributes:

|  значение  | описание  |
|-------------|---------|
| body |  Type: string,  message body   |
| chat_id | Type: int64,  user's id  |
| keys |  Type: []string, keyboard   |

----------------

### telegramAction Function:

```coffeescript
telegramAction = (entityId, actionName, args)->
```

| Value | Description |
|-------|-------------|
| entityId | Type: string, ID of the entity sending the message |
| actionName | Type: string, name of the action in uppercase, without the '/' character |
| args  | attributes |

----------------

### Code Example:

```coffeescript
# telegram
# ##################################
telegramSendReport =->
  entities = ['device.l3n1','device.l3n2','device.l3n3','device.l3n4']
  for entityId, i in entities
    entity = GetEntity(entityId)
    attr = EntityGetAttributes(entityId)
    sendMsg(format(entityId, entity.state.name, attr))
  
telegramAction = (entityId, actionName)->
switch actionName
    when 'CHECK' then telegramSendReport()

sendMsg =(body)->
  msg = notifr.newMessage();
  msg.entity_id = 'telegram.testbot';
  msg.attributes = {
    'body': body
  };
  notifr.send(msg);
```
