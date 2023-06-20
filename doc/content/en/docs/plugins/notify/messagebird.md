
---
title: "Messagebird"
linkTitle: "messagebird"
date: 2021-10-20
description: >

---

To send SMS notifications through the MessageBird notification service provider, you'll need to register with MessageBird and have a positive balance.

### Configuration:
- Access Key
- Name

### Attributes:
- Payment
- Type
- Amount

### JavaScript Properties:

----------------

New Message:

Creates a message object.

```coffeescript
msg = notifr.newMessage();
msg.type = 'messagebird';
msg.attributes = {
  'phone': '+79990000001',
  'body': 'some text'
};
```

| Value | Description |
|-------|-------------|
| newMessage() | Method |
| msg | Type: Object (Message) |

----------------

### Code Example:

```coffeescript
# messagebird
# ##################################

sendMsg =(body)->
  msg = notifr.newMessage();
  msg.type = 'messagebird';
  msg.attributes = {
    'phone': '+79990000001',
    'body': 'some text'
  };
  notifr.send(msg);
```

