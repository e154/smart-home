---
title: "Устройство"
linkTitle: "entity"
date: 2021-11-20
description: >

---

В проекте **Smart Home** объект "Entity" является основной единицей и главным объектом, который хранит состояние и
настройки устройства.

"Entity" представляет собой абстракцию конкретного устройства в умном доме. Он содержит информацию о типе устройства,
его идентификаторе, текущем состоянии и других настройках.

Основные атрибуты "Entity" включают:

1. Тип устройства: Определяет категорию или класс устройства, например, "светильник", "термостат", "дверной замок" и
   т.д.
   Это позволяет классифицировать устройства и определить доступные для них операции и настройки.

2. Идентификатор: Уникальный идентификатор, присвоенный устройству, который позволяет однозначно идентифицировать его в
   системе умного дома. Это может быть числовое значение, строка или другой формат идентификации.

3. Состояние устройства: Содержит информацию о текущем состоянии устройства, например, "включено/выключено",
   "открыто/закрыто", "температура" и т.д. Состояние может изменяться в результате операций или событий, связанных с
   устройством.

4. Настройки устройства: Включает в себя различные параметры и настройки, связанные с конкретным устройством. Например,
   это могут быть параметры освещенности для светильника, предпочитаемая температура для термостата и т.д.

"Entity" является основной абстракцией устройства в проекте **Smart Home** и обеспечивает единый интерфейс для работы с
различными типами устройств. Это позволяет удобно управлять состоянием и настройками устройств, а также создавать гибкие
сценарии автоматизации и взаимодействия между устройствами в вашем умном доме.

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

### методы объекта Entity

```coffeescript
EntitySetState(ENTITY_ID, attr)
EntitySetStateName(ENTITY_ID, stateName)
EntityGetState(ENTITY_ID)
EntitySetAttributes(ENTITY_ID, attr)
EntitySetMetric(ENTITY_ID, name, value)
EntityCallAction(ENTITY_ID, action, args)
EntityCallScene(ENTITY_ID, args)
EntityGetSettings(ENTITY_ID)
EntityGetAttributes(ENTITY_ID)
```

Абстрактный объект "Entity" в проекте **Smart Home** предоставляет следующие методы:

1. `EntitySetState(ENTITY_ID, attr)`: Этот метод используется для установки состояния конкретной сущности (`ENTITY_ID`)
   путем предоставления набора атрибутов (`attr`). С его помощью, скорее всего, можно обновлять состояние сущности
   новыми значениями атрибутов.

    ```javascript
    const attrs = {
        foo: bar,
    }
    const stateName = 'connected'
    EntitySetState(ENTITY_ID, {
        new_state: stateName,
        attribute_values: attrs,
        storage_save: true
    });
    ```

2. `EntitySetStateName(ENTITY_ID, stateName)`: Этот метод устанавливает состояние сущности (`ENTITY_ID`) в определенное
   именованное состояние (`stateName`). Он используется, когда необходимо изменить состояние сущности на
   предопределенное.

    ```javascript
    const stateName = 'connected'
    EntitySetStateName(ENTITY_ID, stateName);
    ```

3. `EntityGetState(ENTITY_ID)`: Этот метод извлекает текущее состояние указанной сущности (`ENTITY_ID`). Он позволяет
   запрашивать текущее состояние сущности.

    ```javascript
  
    const currentState = EntityGetState(ENTITY_ID);
    print(marshal(homeState))
    // out: {"name":"connected","description":"connected","image_url":null,"icon":null}
    ```

4. `EntitySetAttributes(ENTITY_ID, attr)`: Аналогично `EntitySetState`, этот метод устанавливает атрибуты для
   определенной сущности (`ENTITY_ID`). Он может использоваться, когда нужно обновить атрибуты без изменения состояния.

    ```javascript
    const attrs = {
        foo: bar,
    }
    const stateName = 'connected'
    EntitySetAttributes(ENTITY_ID, attrs);
    ```

5. `EntitySetMetric(ENTITY_ID, name, value)`: Этот метод устанавливает метрику для указанной сущности (`ENTITY_ID`).
   Метрики часто используются для измерения и записи различных аспектов поведения или производительности сущности.

    ```javascript
    const attrs = {
        foo: bar,
    }
    const metricName = 'counter'
    EntitySetMetric(ENTITY_ID, name, attrs);
    ```

6. `EntityCallAction(ENTITY_ID, action, args)`: Этот метод используется для вызова или выполнения определенного
   действия (`action`), связанного с сущностью (`ENTITY_ID`), и предоставления необходимых аргументов (`args`) для этого
   действия.

    ```javascript
    const attrs = {
        foo: bar,
    }
    const action = 'ENABLE'
    EntityCallAction(ENTITY_ID, action, attrs);
    ```

7. `EntityCallScene(ENTITY_ID, args)`: Подобно `EntityCallAction`, этот метод используется для вызова или выполнения
   сцены, связанной с сущностью (`ENTITY_ID`), и предоставления необходимых аргументов (`args`) для этой сцены.

    ```javascript
    const attrs = {
        foo: bar,
    }
    EntityCallScene(ENTITY_ID, attrs);
    ```

8. `EntityGetSettings(ENTITY_ID)`: Этот метод извлекает настройки или конфигурацию для указанной сущности (`ENTITY_ID`).
   Он позволяет получать доступ к настройкам, связанным с сущностью.

    ```javascript
    const settings = EntityGetSettings(ENTITY_ID);
    print(marshal(settings))
    // out: {"username":"zigbee2mqtt","cleanSession":false,"keepAlive":15,"direction":"in","topics":"owntracks/#","pingTimeout":10,"connectTimeout":30,"qos":0}

    ```

9. `EntityGetAttributes(ENTITY_ID)`: Этот метод извлекает все атрибуты, связанные с определенной
   сущностью (`ENTITY_ID`). Он предоставляет способ получения подробной информации о атрибутах сущности.

    ```javascript
    const attrs = EntityGetAttributes(ENTITY_ID);
    print(marshal(attrs))
    // out: {"username":"zigbee2mqtt","cleanSession":false,"keepAlive":15,"direction":"in","topics":"owntracks/#","pingTimeout":10,"connectTimeout":30,"qos":0}

    ```

Методы объекта "Entity" предоставляют удобные возможности для управления состоянием, атрибутами, метриками и выполнения
действий на устройстве. Вы можете использовать эти методы для создания логики управления устройствами в вашем проекте
**Smart Home**.



