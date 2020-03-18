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

### .SetElementState(elementName, newState) {#map_set_element_state}

```coffeescript
Map.SetElementState('dev1_light1', 'ERROR')
```

**На входе**

**Значение**            | **Описание**
------------------------|--------------
  `elementName`         | type: string, наименование элемента
  `newState`            | type: string, системное наименование состояния

### .GetElement(elementName) {#map_get_element}

```coffeescript
element = Map.GetElement('dev1_light1')
```

**На входе**

**Значение**    | **Описание**
----------------|--------------
  `elementName` | type: string, наименование элемента


**На выходе**

**Значение**    | **Описание**
----------------|--------------
 `element`      | type: Object, ссылка на объектов [MapElement{}](#map_element)

