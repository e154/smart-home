
---
title: "Actions"
linkTitle: "Actions"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/actions_window1.png" >}}

&nbsp;

&nbsp;

Actions define the tasks that need to be performed when the scenario is executed. This can include changing the state of devices, sending commands to other devices, executing an HTTP request, sending notifications, and other actions.

Example implementation of an action handler:
```coffeescript
automationAction = (entityId)->
    EntityCallAction(entityId, 'ON', {})
```
