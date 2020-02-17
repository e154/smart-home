---
weight: 30
title: log
groups:
    - javascript
---

## Map* {#map}

Возвращает объект карты

```coffeescript
map = Map
```
  
Доступные методы приведены далее:

### .SetElementState(device, elementName, newState) {#map_set_element_state}

```coffeescript
Map.SetElementState(device, 'dev1_light1', 'ERROR')
```

**На входе**

**Значение**            | **Описание**
------------------------|--------------
  `device`              | type: Object, ссылка на объект [DeviceModel{}](#device)
  `elementName`         | type: string, наименование элемента
  `newState`            | type: string, системное наименование состояния

### .GetElement(device, elementName) {#map_get_element}

```coffeescript
elements = Map.GetElement(device, 'dev1_light1')
```

**На входе**

**Значение**    | **Описание**
----------------|--------------
  `device`      | type: Object, ссылка на объект [Device{}](#device)
  `elementName` | type: string, наименование элемента


**На выходе**

**Значение**    | **Описание**
----------------|--------------
 `elements`     | type: Object, ссылка на объектов [MapElement{}](#map_element)

### .GetElements(device) {#map_get_element}

```coffeescript
elements = Map.GetElements(device)
```

**На входе**

**Значение**    | **Описание**
----------------|--------------
  `device`      | type: Object, ссылка на объект [Device{}](#device)

**На выходе**

**Значение**    | **Описание**
----------------|--------------
 `elements`     | type: Array, массив объектов [MapElement{}](#map_element)
