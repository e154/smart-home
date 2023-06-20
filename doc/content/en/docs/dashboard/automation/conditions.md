
---
title: "Conditions"
linkTitle: "Conditions"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/condition_window1.png" >}}

&nbsp;

&nbsp;

A condition defines an additional check that needs to be performed before executing the scenario. This condition can check the current state of a device or other factors. The condition is an optional component, and if it is present, the execution of the scenario will depend on its result.

Example implementation of a condition handler:
```coffeescript
automationCondition = (entityId)->
    #print '---condition---'
    entity = Condition.getEntityById('zigbee2mqtt.` + zigbeePlugId + `')
    if !entity || !entity.state 
        return false
    if entity.state.name == 'ON'
        return true
    return false
```

Please note that a condition is an optional component. If a condition is present, a check will be performed before executing the actions. If the condition returns a positive result, the actions will be executed. Otherwise, if the condition returns a negative result or is not specified, the actions will be executed without further checking.
