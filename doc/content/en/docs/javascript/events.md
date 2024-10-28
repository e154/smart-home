---
title: "Events"
linkTitle: "events"
date: 2024-01-13
description: >

---

In the Smart Home system, there is a JavaScript function called **PushSystemEvent** that plays a crucial role in dynamic
system management. This function accepts various commands for interacting with tasks, triggers, and other components of
the system.

| Command                   | Description                                               |
|---------------------------|-----------------------------------------------------------|
| `command_enable_task`     | Enables the execution of a task.                          |
| `command_disable_task`    | Disables the execution of a task.                         |
| `command_enable_trigger`  | Enables a trigger, activating the triggering capability.  |
| `command_disable_trigger` | Disables a trigger, suspending the triggering capability. |
| `event_call_trigger`      | Initiates a trigger call event.                           |
| `event_call_action`       | Initiates an action call event.                           |
| `command_load_entity`     | Loads an entity into the system.                          |
| `command_unload_entity`   | Unloads an entity from the system.                        |

#### Example of Usage:

```javascript
// Example of enabling a task
PushSystemEvent('command_enable_task', { id: 1 });

// Example of calling a trigger
PushSystemEvent('event_call_trigger', { id: 1 });

// Example of loading an entity
PushSystemEvent('command_load_entity', { id: 'sensor.entity1' });
```

These commands provide control over tasks, triggers, trigger call events, actions, as well as loading and unloading
entities in the Smart Home system. Their use enables dynamic management of system functionality and interaction with its
components.
