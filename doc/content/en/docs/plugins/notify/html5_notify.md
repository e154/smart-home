---
title: "HTML5 Уведомления"
linkTitle: "html5_notify"
date: 2021-10-20
description: >

---

HTML5 notifications are a way to display information to users on a web page without reloading or navigating to another
page. Notifications can contain text, images, sound, or video and appear in the bottom-right corner of the screen. HTML5
notifications utilize an API that allows web applications to send notifications to the browser using JavaScript.

Attributes:

- userIDS: The user IDs associated with the notification.
- title: The title of the notification.
- body: The body or content of the notification.

JavaScript Properties:

----------------

### New Message:

Creates a message object.

```coffeescript
msg = notifr.newMessage();
msg.type = 'html5_notify';
msg.attributes = {
  'userIDS': '14',
  'title': 'Lorem Ipsum is simply',
  'body': 'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum'
};
```

| Value        | Description            |
|--------------|------------------------|
| newMessage() | Method                 |
| msg          | Type: Object (Message) |

----------------

### Code Example:

```coffeescript
# html5_notify
# ##################################

sendMsg = ()->
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

