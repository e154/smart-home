---
title: "Actions"
linkTitle: "Actions"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/actions_window1.png" >}}

&nbsp;

&nbsp;

Действия определяют задачи, которые должны быть выполнены при выполнении сценария. Это может быть изменение состояния
устройств, отправка команды другому устройству, выполнение HTTP-запроса, отправка уведомления и другие действия.

Пример реализации обработчика:

```coffeescript
automationAction = (entityId)->
  EntityCallAction(entityId, 'ON', {})
```
