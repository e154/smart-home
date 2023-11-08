
---
title: "Telegram"
linkTitle: "telegram"
date: 2021-10-20
description: >

---

Позволяет работать с telegram api, обеспечивая интерактивную обработку комманд от клиентов либо использовать
событийное оповещение. Система не имеет ограничение на количество ботов, один entity - один бот.

### Настройка
* Token

### Команды:
* произвольный набор

### Атрибуты
* произвольный набор

### Действия

В системе зарезервированно две команды:

* **/start** - подписаться на уведомления
* **/quit** - отписаться от уведомлений

**Action** (действие) - должно наименоваться в нижнем регистре, без знаков "/". 
Кастомная команда в системе автоматически добавится в список доступных команд раздела **/help**.
Вызов кастомной команды с клиента следует производить в верхнем регистре **/ACTION** -> **/action**


### javascript свойства


----------------

### Новое сообщение

создает объект сообщения

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
| body |  Type: string,  тело сообщения   |
| chat_id | Type: int64,  id пользователя  |
| keys |  Type: []string, клавиатура   |

----------------

### функция telegramAction

```coffeescript
telegramAction = (entityId, actionName, args)->
```
| значение   | описание               |
|-------------|-------------------|
| entityId    | type: string, id сущности отправляющего сообщение |
| actionName  | type: string, название действия, без символа '/' в верхнем регистре |
| args  | attributes |

----------------

### пример кода

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
