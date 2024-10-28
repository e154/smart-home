---
title: "Slack"
linkTitle: "slack"
date: 2021-10-20
description: >

---

To send notifications to Slack chat, you'll need to configure the following settings:

### Configuration:

- Token
- User name

### JavaScript Properties:

----------------

New Message:

Creates a message object.

```coffeescript
msg = notifr.newMessage();
msg.entity_id = 'slack.name';
msg.attributes = {
  'channel': '#ch',
  'text': 'some text'
};
```

| Value        | Description            |
|--------------|------------------------|
| newMessage() | Method                 |
| msg          | Type: Object (Message) |

----------------

### Code Example:

```coffeescript
# slack
# ##################################

sendMsg = (body)->
  msg = notifr.newMessage();
  msg.entity_id = 'slack.name';
  msg.attributes = {
    'channel': '#ch',
    'text': 'some text'
  };
  notifr.send(msg);
```

