---
title: "entityManager"
linkTitle: "entityManager"
date: 2021-11-20 
description: >

---

В проекте **Smart Home** имеется менеджер устройств "Entity Manager", который обеспечивает удобный доступ к устройствам, 
их настройкам и текущему состоянию.

{{< alert color="success" >}}Функция доступна в любом скрипте системы.{{< /alert >}}

Менеджер устройств "Entity Manager" предоставляет набор методов и функций для управления устройствами и получения 
информации о них. Некоторые основные возможности "Entity Manager" включают:


1. `getEntity(id) -> Entity`: Этот метод позволяет получить объект "Entity" по его идентификатору. Вы передаете идентификатор устройства в качестве аргумента `id`, и метод возвращает соответствующий объект "Entity".

    ```javascript
    const entity = entityManager.getEntity(id);
    ```

2. `setState(id, stateName)`: Данный метод используется для установки состояния устройства с заданным идентификатором. Вы передаете идентификатор устройства в аргументе `id` и имя состояния в аргументе `stateName`. Например:

    ```javascript
    entityManager.setState(id, 'on');
    ```

3. `setAttributes(id, attrs)`: Этот метод позволяет установить атрибуты устройства с заданным идентификатором. Вы передаете идентификатор устройства в аргументе `id` и объект с новыми атрибутами в аргументе `attrs`.

    ```javascript
    const newAttributes = {
      brightness: 80,
      color: 'red',
    };
    
    entityManager.setAttributes(id, newAttributes);
    ```

4. `setMetric(id, name, value)`: Этот метод используется для установки метрики устройства с заданным идентификатором. Вы передаете идентификатор устройства в аргументе `id`, имя метрики в аргументе `name` и значение метрики в аргументе `value`.

    ```javascript
    entityManager.setMetric(id, 'temperature', 25);
    ```

5. `callAction(id, actionName, args)`: Данный метод позволяет вызвать определенное действие на устройстве с заданным идентификатором. Вы передаете идентификатор устройства в аргументе `id`, имя действия в аргументе `actionName` и аргументы, связанные с действием, в аргументе `args`.

    ```javascript
    entityManager.callAction(id, 'turnOn');
    ```

6. `callScene(id, args)`: Этот метод используется для вызова сцены с заданным идентификатором. Вы передаете идентификатор сцены в аргументе `id` и аргументы, связанные со сценой, в аргументе `args`.

    ```javascript
    entityManager.callScene(id, args);
    ```

Менеджер устройств "Entity Manager" упрощает управление устройствами, настройками и состоянием в проекте **Smart Home**.
Он предоставляет удобный интерфейс для получения информации о устройствах, управления ими и обработки связанных с ними 
событий. Это позволяет вам создавать более гибкие и функциональные приложения для вашего умного дома.

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
  .callScene(id, args)
```

|  значение  | описание  |
|-------------|---------|
| id | уникальному id устройства |
| getEntity | поиск устройства, возвращает [Entity](#методы-объекта-entity) |
| setState | изменить текущее состояние, **stateName** - наименования будущего состояния|
| setAttributes | изменить свойство текущего состояния, **attr** - map параметрый {key:val} |
| setMetric | обновить метрики устройства, **name** - наименование метрики, **attr** - map параметрый {key:val}|
| callAction | иcполнение команды на устройстве, **actionName** - наименование команды, **args** - map параметрый {key:val}|
| callScene | выполнить сцену с **id** |



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

# получение настроек
settings = device.getSettings()
headers = [{'apikey':settings['apikey']}]
url = 'https://webhook.site/2692a589-a5bc-4156-af7d-75875578798f'
res = http.headers(headers).get(url)

```
