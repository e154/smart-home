
---
title: "Push уведомления"
linkTitle: "webpush"
date: 2022-06-08
description: >

---

Web push notifications are short messages that Smart Home sends to the browser to convey important information. They are displayed even to users who have left the website and do not require email or phone number input.

### Configuration:
* Public Key
* Private Key

### Attributes:
* userIDS
* title
* body

### JavaScript Properties:

----------------

### New Message

Creates a message object

```coffeescript
msg = notifr.newMessage();
msg.type = 'webpush';
msg.attributes = {
  'userIDS': '14',
  'title': 'Lorem Ipsum is simply',
  'body': 'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum'
};

```
| Value | Description |
|-------|-------------|
| newMessage() | Method |
| msg | Type: Object (Message) |

----------------

### Code Example:

```coffeescript
# webpush
# ##################################
sendMsg =()->
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

