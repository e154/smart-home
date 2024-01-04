
---
title: "EMAIL"
linkTitle: "email"
date: 2024-01-04
description: >
  
---

The Email plugin provides the ability to integrate email sending functionality into the Smart Home system. This plugin allows users to configure email parameters and send notifications or messages from devices within the smart home.

#### Device Settings

- **`email` (type: String)**: The email through which messages are sent.

- **`auth` (type: String)**: Authentication method for accessing the mail server (e.g., "PLAIN" or "LOGIN").

- **`pass` (type: Encrypted)**: Encrypted password for authentication when accessing the mail server.

- **`smtp` (type: String)**: Outgoing mail server (SMTP) address.

- **`port` (type: Int)**: Outgoing mail server port (e.g., 587 for TLS or 465 for SSL).

- **`sender` (type: String)**: Sender's email address.

#### Message Attributes

- **`addresses` (type: String)**: A list of recipient email addresses, separated by commas.

- **`subject` (type: String)**: The subject of the email message.

- **`body` (type: String)**: The body text of the email message.

#### Example in Coffeescript

```coffeescript
msg = notifr.newMessage();
msg.entity_id = 'email.google';
msg.attributes = {
  'body': 'some text msg',
  'addresses': 'john.smith@example.com,jane.doe@example.com',
  'subject': 'Important Notification',
};
```
