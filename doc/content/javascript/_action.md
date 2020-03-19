---
weight: 36
title: action
groups:
    - javascript
---

## Action {#action}

Действие.

```coffeescript
action = Action
action.Id
action.Name
action.Description
action.Device()
action.Node()
```

**На выходе**

**Значение**    | **Описание**
----------------|--------------
  `Id`          | type: int
  `Name`        | type: string
  `Description` | type: string
  `Device()`    | type: Object, ссылка на экземпляр [Device](#device) 
  `Node()`      | type: Object, ссылка на экземпляр [Node](#node) 
  
