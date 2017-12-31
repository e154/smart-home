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

### .SetElementState(element, state) {#ic_map_set_element_state}

```coffeescript
IC.Map.SetElementState(device, state)
```

**На входе**

**Значение**    | **Описание**
----------------|--------------
  `element`      | type: Object, ссылка на объект карты
  `state`       | type: string, системное наименование состояния

### .GetElement(device) {#ic_map_get_element}

```coffeescript
element = IC.Map.GetElement(device)
```

**На входе**

**Значение**    | **Описание**
----------------|--------------
  `device`      | type: Object, ссылка на объект [Device{}](#device)


**На выходе**

**Значение**    | **Описание**
----------------|--------------
 `element`      | type: Object, ссылка на объект [MapElement{}](#map_element)
