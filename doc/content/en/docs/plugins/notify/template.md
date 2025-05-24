---
title: "Templates"
linkTitle: "templates"
date: 2021-10-20
description: >

---

The **Template** feature allows you to generate and distribute pre-defined messages. It offers the following advantages:

* Clean code
* Different templates for different types
* Convenient invocation

### JavaScript Properties:

----------------

### Generate Message from Template

```coffeescript
tpl = template.render(name, params)
```

| Value  | Description                                             |
|--------|---------------------------------------------------------|
| name   | Type: string, name of the template                      |
| params | Type: Object, template parameters (e.g., {'key':'val'}) |
| tpl    | Type: string, generated message ready to be sent        |

----------------

### Code Example:

```coffeescript
# telegram
# ##################################


sendMsg = (body)->
  tpl = template
    .render('name', {'key': 'val'})

  msg = notifr.newMessage();
  msg.type = 'telegram';
  msg.attributes = {
    'name': 'clavicus',
    'body': tpl
  };
  notifr.send(msg);
```
