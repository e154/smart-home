
---
title: "Triggers"
linkTitle: "Triggers"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/triggers_window1.png" >}}

&nbsp;

&nbsp;

Each type of trigger has its own handler for performing the corresponding actions.

Here are examples of different types of triggers and their respective handlers:

1. `TriggerAlexa`:
```coffeescript
automationTriggerAlexa = (msg) ->
  p = msg.payload
  Done p
  return false
```
The `automationTriggerAlexa` handler is invoked in response to a trigger from Amazon Alexa. It takes the `msg` message object and can perform specific actions related to this trigger.

2. `TriggerStateChanged`:
```coffeescript
automationTriggerStateChanged = (msg) ->
  print '---trigger---'
  p = msg.payload
  Done p.new_state.state.name
  return false
```
The `automationTriggerStateChanged` handler is called when the state of a device changes. It takes the `msg` message object and can perform specific actions based on the new state of the device.

3. `TriggerSystem`:
```coffeescript
automationTriggerSystem = (msg) ->
  p = msg.payload
  Done p.event
  return false
```
The `automationTriggerSystem` handler is invoked in response to system events. It takes the `msg` message object and can perform specific actions related to this event.

4. `TriggerTime`:
```coffeescript
automationTriggerTime = (msg) ->
  p = msg.payload
  Done p
  return false
```
The `automationTriggerTime` handler is called after a certain amount of time has passed. It takes the `msg` message object and can perform specific actions related to this time.

Each trigger handler can execute the necessary logic in response to the corresponding trigger and then return `false` to indicate that further processing is not required.

Example implementation:

```coffeescript
automationTriggerStateChanged = (msg) ->
    #print '---trigger---'
    p = msg.payload
    if !p.new_state || !p.new_state.state
        return false
    return msg.new_state.state.name == 'DOUBLE_CLICK'
```
