
---
title: "Twilio"
linkTitle: "twilio"
date: 2021-10-20
description: >

---

Sending SMS notifications through the notification service provider [Twilio](https://www.twilio.com/messaging). To be able to send messages, registration with [Twilio](https://www.twilio.com/messaging) and a positive balance are required.

### Configuration:
* Token
* Sid
* From

### Attributes:
* Amount
* Sid
* Currency

### JavaScript Properties:

----------------

### New Message

Creates a message object

```coffeescript
msg = notifr.newMessage();
msg.entity_id = 'twilio.name';
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
# twilio
# ##################################

sendMsg =(body)->
  msg = notifr.newMessage();
  msg.entity_id = 'twilio.name';
  msg.attributes = {
    'phone': '+79990000001',
    'body': 'some text'
  };
  notifr.send(msg);
```

