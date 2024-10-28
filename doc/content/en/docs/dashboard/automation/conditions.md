---
title: "Conditions"
linkTitle: "Conditions"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/condition_window1.png" >}}

&nbsp;

&nbsp;

A condition defines an additional check that needs to be performed before executing the scenario. This condition can
check the current state of a device or other factors. The condition is an optional component, and if it is present, the
execution of the scenario will depend on its result.

Example implementation of a condition handler:

```coffeescript
automationCondition = (entityId)->
  entity = EntityGetState(entityId)
  if !entity || !entity.state
    return false
  if entity.state.name == 'ON'
    return true
  return false
```
