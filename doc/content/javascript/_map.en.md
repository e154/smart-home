---
weight: 30
title: log
groups:
    - javascript
---

## IC.Map* {#ic_map}

Возвращает объект карты

```coffeescript
map = IC.Map
```
  
Доступные методы приведены далее:

### .SetElementState(device, elementName, newState) {#ic_map_set_element_state}

```coffeescript
IC.Map.SetElementState(device, state)
```

**На входе**

**Значение**            | **Описание**
------------------------|--------------
  `device`              | type: Object, ссылка на объект [DeviceModel{}](#device)
  `elementName`         | type: string, наименование элемента
  `stnewStateate`       | type: string, системное наименование состояния

### .GetElement(device, elementName) {#ic_map_get_element}

```coffeescript
elements = IC.Map.GetElement(device, elementName)
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

### .GetElements(device) {#ic_map_get_element}

```coffeescript
elements = IC.Map.GetElement(device)
```

**На входе**

**Значение**    | **Описание**
----------------|--------------
  `device`      | type: Object, ссылка на объект [Device{}](#device)

**На выходе**

**Значение**    | **Описание**
----------------|--------------
 `elements`     | type: Array, массив объектов [MapElement{}](#map_element)
