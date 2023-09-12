
---
title: "Bitmine"
linkTitle: "bitmine"
date: 2021-10-20
description: >

---

Работа с устройствами ASIC производителя Bitmine Antminer.


поддерживаемые усторойства:
* S9
* S7
* L3
* L3
* D3
* T9

### javascript свойства

----------------

### Miner Базовая сущность 

```coffeescript
# константа с уникальным id устройства
const ENTITY_ID
````

```coffeescript
result = Miner
    .stats()
    .devs()
    .summary()
    .pools()
    .addPool(url string)
    .version()
    .enable(poolId int64)
    .disable(poolId int64)
    .delete(poolId int64)
    .switchPool(poolId int64)
    .restart()
```

|  значение  | описание  |
|-------------|---------|
| result |    type: Object (Result)   |

|  значение  | описание  |
|-------------|---------|
| stats |   type: метод, статистика    |
| devs |   type: метод, статистика    |
| summary |   type: метод, статистика    |
| pools |   type: метод, список серверов   |
| addPool |   type: метод, добавление сервера   |
| version |   type: метод, версия прошивки  |
| enable |   type: метод, включение в работу сервера  |
| disable |   type: метод, включение сервера  |
| delete |   type: метод, удаление сервера  |
| switchPool |   type: метод, переключние сервера  |
| restart |   type: метод, мягкий перезапуск устройства  |

----------------

### объект Result

```coffeescript
  result = {
  error: false,
  errMessage: "",
  result: ""
}
``` 
|  значение  | описание  |
|-------------|---------|
| error |    type: boolean, признак ошибки   |
| errMessage |   type: string, человекопонятное описание ошибки  |
| result | type: string, json объект в текстовом формате, если запрос завершился без ошибок |

----------------

### функция entityAction

```coffeescript
entityAction = (entityId, actionName, args)->
```
| значение   | описание               |
|-------------|-------------------|
| entityId    | type: string, id сущности отправляющего сообщение |
| actionName  | type: string, название действия, без символа '/' в верхнем регистре |
| args | Type: map[string]any |

----------------

### Обект Actor

```coffeescript
Actor.setState(entityStateParams)
```
| значение   | описание               |
|-------------|-------------------|
| setState  | type: метод, обновление статуса сущности |
| entityStateParams  | type: Object EntityStateParams |

----------------

### объект EntityStateParams

```coffeescript
  entityStateParams = {
    "new_state": ""
    "attribute_values": {}
    "settings_value": {}
    "storage_save": true
}
``` 
|  значение  | описание  |
|-------------|---------|
| new_state |    type: *string, наименование статуса   |
| attribute_values |   type: Object, новое состяние аттрибутов сущности |
| settings_value | type: Object, новое состояние настроек сущности |
| storage_save | type: boolean, признак записи нового состояния в бд |

----------------

### пример кода

```coffeescript
# cgminer
# ##################################

ifError =(res)->
  return !res || res.error || res.Error

checkStatus =->
  stats = Miner.stats()
  if ifError(stats)
    Actor.setState
      'new_state': 'ERROR'
    return
  p = JSON.parse(stats.result)
  print p

checkSum =->
  summary = Miner.summary()
  if ifError(summary)
    Actor.setState
      'new_state': 'ERROR'
    return
  p = JSON.parse(summary.result)
  print p

entityAction = (entityId, actionName, args)->
  switch actionName
    when 'CHECK' then checkStatus()
    when 'SUM' then checkSum()
```
