
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

The system reserves three commands:

* **/start** - Subscribe to notifications
* **/quit** - Unsubscribe from notifications
* **/help** - Display available commands

**Action** - It should be named in lowercase, without the "/" sign. A custom command will automatically be added to the list of available commands in the **/help** section. The custom command should be called in uppercase **/ACTION** -> **/action**.

### JavaScript Properties:

----------------

New Message:

Creates a message object.

```coffeescript
msg = notifr.newMessage();
msg.entity_id = 'telegram.name';
msg.attributes = {
  'name': 'clavicus',
  'body': 'some text msg'
};
```

| Value | Description |
|-------|-------------|
| newMessage() | Method |
| msg | Type: Object (Message) |

----------------

telegramAction Function:

```coffeescript
telegramAction = (entityId, actionName)->
```

| Value | Description |
|-------|-------------|
| entityId | Type: string, ID of the entity sending the message |
| actionName | Type: string, name of the action in uppercase, without the '/' character |

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
  msg.entity_id = 'telegram.name';
  msg.attributes = {
    'name': 'clavicus',
    'body': body
  };
  notifr.send(msg);
```
