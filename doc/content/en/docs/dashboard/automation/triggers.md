
---
title: "Triggers"
linkTitle: "Triggers"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/triggers_window1.png" >}}

Triggers in the Smart Home system are key elements that initiate the execution of tasks in response to specific events or conditions. Currently, the system supports three main types of triggers, each with its advantages, providing flexibility in automation configuration.

#### 1. **Trigger "State Change":**

- **Description:**
    - Activates upon the change of the state of a specific device or group of devices in the Smart Home system.

- **Advantages:**
    - **Reacting to Changes:** Ideal for scenarios where specific reactions are needed, such as turning lights on/off, opening doors, etc.
    - **Customization Flexibility:** Allows selecting specific devices and state parameters for triggering.

#### 2. **Trigger "System":**

- **Description:**
    - Activates upon the occurrence of system events, such as system start/stop, device connection/disconnection, and other events on the system bus.

- **Advantages:**
    - **System Monitoring:** Used for monitoring the overall system state and detecting events related to the functioning of the smart home.
    - **Integration with External Systems:** Enables interaction with systems operating within the system bus.

#### 3. **Trigger "Time (Cron)":**

- **Description:**
    - Activates at specified time intervals according to a schedule defined in cron format.

- **Advantages:**
    - **Regular Execution:** Ideal for tasks that need to be executed regularly on a schedule.
    - **Energy Savings:** Can be used to optimize energy consumption, such as turning off lights or heating during nighttime.

#### Examples of Trigger Usage

1. **Trigger "State Change":**
    - **Scenario:**
        - Turn on the air conditioner when the window in the bedroom is opened.
    - **Trigger Configuration:**
        - Device: Window in the bedroom.
        - State: Open.

2. **Trigger "System":**
    - **Scenario:**
        - Send a notification upon connecting a new device to the system.
    - **Trigger Configuration:**
        - System Event: Connection of a new device.

3. **Trigger "Time (Cron)":**
    - **Scenario:**
        - Turn off the lights in all rooms after midnight.
    - **Trigger Configuration:**
        - Time Interval: "0 0 * * *" corresponding to every midnight.

#### Examples in Coffeescript

1. `TriggerAlexa`:
```coffeescript
automationTriggerAlexa = (msg) ->
  p = msg.payload
  Done p
  return false
```
The `automationTriggerAlexa` handler is called in response to a trigger from Amazon Alexa. It takes the message object `msg` and can perform specific actions related to this trigger.

2. `TriggerStateChanged`:
```coffeescript
automationTriggerStateChanged = (msg) ->
  print '---trigger---'
  p = msg.payload
  Done p.new_state.state.name
  return false
```
The `automationTriggerStateChanged` handler is called when the state of a device changes. It takes the message object
`msg` and can perform specific actions based on the new state of the device.

3. `TriggerSystem`:
```coffeescript
automationTriggerSystem = (msg) ->
  p = msg.payload
  Done p.event
  return false
```
The `automationTriggerSystem` handler is called in response to system events. It takes the message object `msg` and can perform specific actions related to this event.

4. `TriggerTime`:
```coffeescript
automationTriggerTime = (msg) ->
  p = msg.payload
  Done p
  return false
```
The `automationTriggerTime` handler is called after a specified period of time has elapsed. It takes the message object `msg` and can perform specific actions related to this time.

Each trigger handler can execute the necessary logic in response to the corresponding trigger and then return the value `false` to indicate that further processing is not required.

Example implementation:

```coffeescript
automationTriggerStateChanged = (msg)->
    #print '---trigger---'
    p = msg.payload
    if !p.new_state || !p.new_state.state
        return false
    return msg.new_state.state.name == 'DOUBLE_CLICK'
```

#### Conclusion

The use of different types of triggers in the Smart Home system allows users to create complex and intelligent automation scenarios that seamlessly respond to changes in the system and external events. This provides users with a personalized and efficient smart home experience.
