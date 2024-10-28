---
title: "Шаблон сообщения"
linkTitle: "шаблоны"
date: 2021-10-20
description: >

---

**Template** (Шаблон) - инструмент для заполнения и рассылки готовых сообщений.

преимущества:

* чистый код
* разные шаблоны для разных типов
* удобный вызов

### javascript свойства

----------------

### Генерация сообщения из шаблона

```coffeescript
  tpl = template
  .render(name, params)
```

| значение | описание                                                    |
|----------|-------------------------------------------------------------|
| name     | type: string, название шаблона                              |
| params   | type: Object, параметры шаблона, пример: {'key':'val'}      |
| tpl      | type: string, сгенерированное сообщение, готовое к отправке |

### пример кода

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
