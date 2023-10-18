
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

В системе зарезервированно три команды:

* **/start** - подписаться на уведомления
* **/quit** - отписаться от уведомлений
* **/help** - вывод доступных команд

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
  'body': 'some text msg'
};

```
|  значение  | описание  |
|-------------|---------|
| newMessage() |    метод   |
| msg |   type: Object (Message)  |

----------------

### функция telegramAction

```coffeescript
telegramAction = (entityId, actionName)->
```
| значение   | описание               |
|-------------|-------------------|
| entityId    | type: string, id сущности отправляющего сообщение |
| actionName  | type: string, название действия, без символа '/' в верхнем регистре |


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
  msg.entity_id = 'telegram.name';
  msg.attributes = {
    'body': body
  };
  notifr.send(msg);
```
