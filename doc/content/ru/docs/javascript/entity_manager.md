---
title: "Менеджер устройств"
linkTitle: "entity manager"
date: 2021-11-20 
description: >

---

Предостовляет возможность работать напрямую с устройствами. Получение/изменение настроек, получение/изменение текущего
состояния.

{{< alert color="success" >}}Функция доступна в любом скрипте системы.{{< /alert >}}

----------------

### объект entityManager

```coffeescript
# доступные метода
entityManager
  .getEntity(id) -> Entity
  .setState(id, stateName)
  .setAttributes(id, attrs)
  .setMetric(id, name, value)
  .callAction(id, actionName, args)
```

|  значение  | описание  |
|-------------|---------|
| id | уникальному id устройства |
| getEntity | поиск устройства, возвращает [Entity](#методы-объекта-entity) |
| setState | изменить текущее состояние, **stateName** - наименования будущего состояния|
| setAttributes | изменить свойство текущего состояния, **attr** - map параметрый {key:val} |
| setMetric | обновить метрики устройства, **name** - наименование метрики, **attr** - map параметрый {key:val}|
| callAction | иcполнение команды на устройстве, **actionName** - наименование команды, **args** - map параметрый {key:val}|

### методы объекта Entity

```coffeescript
Entity
  .setState(stateName)
  .setAttributes(key, attrs)
  .getAttributes() -> attrs
  .setMetric(name, value)
  .callAction(name, value)
```

|  значение  | описание  |
|-------------|---------|
| setState | изменить текущее состояние, **stateName** - наименования будущего состояния|
| setAttributes | изменить свойство текущего состояния, **attr** - map параметрый {key:val} |
| getAttributes | текущее состояние устройства, **attrs** - map параметрый {key:val} |
| setMetric | обновить метрики устройства, **name** - наименование метрики, **attr** - map параметрый {key:val}|
| callAction | иcполнение команды на устройстве, **actionName** - наименование команды, **args** - map параметрый {key:val}|
| short | срез свойств состояния и настроек устройства |

### свойства объекта Entity
```coffeescript
{
  "id": "sensor.device1",
  "type": "sensor",
  "name": "",
  "description": "device description",
  "icon": null,
  "image_url": null,
  "actions": [
    {
      "name": "ENABLE",
      "description": "включить",
      "image_url": null,
      "icon": null
    },
    {
      "name": "DISABLE",
      "description": "выключить",
      "image_url": null,
      "icon": null
    }
  ],
  "states": [
    {
      "name": "ENABLED",
      "description": "enabled state",
      "image_url": null,
      "icon": null
    },
    {
      "name": "DISABLED",
      "description": "disabled state",
      "image_url": null,
      "icon": null
    }
  ],
  "state": {
    "name": "ENABLED",
    "description": "enabled state",
    "image_url": null,
    "icon": null
  },
  "attributes": {},
  "area": null,
  "metrics": [],
  "hidden": false
}
```

----------------

### пример кода

```coffeescript
# entity manager
# ##################################

# поиск устройства по уникальному id, возвращает [Entity](#объект-entity)
device = entityManager.getEntity('telegram.clavicus')

# изменить текущее состояние по уникальному наименованию состояния
entityManager.setState('telegram.clavicus', 'IN_WORK')

# изменить свойство текущего состояния
entityManager.setAttributes('telegram.clavicus', {'temp': 23.3})

# обновить метрики устройства
entityManager.setMetric('telegram.clavicus', 'cpu', {'all': 55})

# иcполнение команды на устройстве
entityManager.callAction('floor.light', 'ON', {'power': 30})


```
