---
title: "Conditions"
linkTitle: "Conditions"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/condition_window1.png" >}}

&nbsp;

&nbsp;

Условие определяет дополнительную проверку, которая должна быть выполнена перед выполнением сценария. Это условие может
проверять текущее состояние устройства или другие факторы. Условие является опциональным компонентом, и если оно
присутствует,
то выполнение сценария будет зависеть от его результата.

Пример реализации обработчика:

```coffeescript
automationCondition = (entityId)->
  entity = EntityGetState(entityId)
  if !entity || !entity.state
    return false
  if entity.state.name == 'ON'
    return true
  return false
```
