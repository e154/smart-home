---
title: "Entity"
linkTitle: "entity"
date: 2021-11-20
description: >

---

In the **Smart Home** project, the "Entity" object is the main unit and the primary object that stores the device's state and settings.

An "Entity" represents an abstraction of a specific device in a smart home. It contains information about the device type, its identifier, current state, and other settings.

The key attributes of an "Entity" include:

1. Device Type: Determines the category or class of the device, such as "light fixture," "thermostat," "door lock," etc. This helps classify devices and define the operations and settings available to them.

2. Identifier: A unique identifier assigned to the device, which allows it to be uniquely identified within the smart home system. This can be a numerical value, string, or other identification format.

3. Device State: Contains information about the current state of the device, such as "on/off," "open/closed," "temperature," etc. The state can change as a result of operations or events related to the device.

4. Device Settings: Includes various parameters and settings related to the specific device. For example, these could be brightness settings for a light fixture, preferred temperature for a thermostat, etc.

The "Entity" serves as the main device abstraction in the **Smart Home** project and provides a unified interface for working with different types of devices. This allows convenient management of device states and settings, as well as the creation of flexible automation scenarios and interactions between devices in your smart home.

### Entity Object Properties
```coffeescript
{
  "id": "sensor.device1",
  "type": "sensor",
  "name": "",
  "description": "device description",
  "icon": null,
  "image_url": null,
  "actions": [
    {
      "name": "ENABLE",
      "description": "enable",
      "image_url": null,
      "icon": null
    },
    {
      "name": "DISABLE",
      "description": "disable",
      "image_url": null,
      "icon": null
    }
  ],
  "states": [
    {
      "name": "ENABLED",
      "description": "enabled state",
      "image_url": null,
      "icon": null
    },
    {
      "name": "DISABLED",
      "description": "disabled state",
      "image_url": null,
      "icon": null
    }
  ],
  "state": {
    "name": "ENABLED",
    "description": "enabled state",
    "image_url": null,
    "icon": null
  },
  "attributes": {},
  "area": null,
  "metrics": [],
  "hidden": false
}
```

### Entity Object Methods

```coffeescript
EntitySetState(ENTITY_ID, attr)
EntitySetStateName(ENTITY_ID, stateName)
EntityGetState(ENTITY_ID)
EntitySetAttributes(ENTITY_ID, attr)
EntitySetMetric(ENTITY_ID, name, value)
EntityCallAction(ENTITY_ID, action, args)
EntityCallScene(ENTITY_ID, args)
EntityGetSettings(ENTITY_ID)
EntityGetAttributes(ENTITY_ID)
```

The abstract "Entity" object in the **Smart Home** project provides the following methods:

1. `EntitySetState(ENTITY_ID, attr)`: This method is used to set the state of a specific entity (`ENTITY_ID`) by providing a set of attributes (`attr`). It likely allows you to update the state of an entity with new attribute values.

2. `EntitySetStateName(ENTITY_ID, stateName)`: This method sets the state of an entity (`ENTITY_ID`) to a specific named state (`stateName`). It's used when you want to change the state of an entity to a predefined state.

3. `EntityGetState(ENTITY_ID)`: This method retrieves the current state of a specified entity (`ENTITY_ID`). It allows you to query the current state of the entity.

4. `EntitySetAttributes(ENTITY_ID, attr)`: Similar to `EntitySetState`, this method sets attributes for a specific entity (`ENTITY_ID`). It may be used when you want to update attributes without changing the state.

5. `EntitySetMetric(ENTITY_ID, name, value)`: This method sets a metric for the specified entity (`ENTITY_ID`). Metrics are often used to measure and record various aspects of an entity's behavior or performance.

6. `EntityCallAction(ENTITY_ID, action, args)`: This method is used to trigger or call a specific action (`action`) associated with the entity (`ENTITY_ID`) and provides any necessary arguments (`args`) for that action.

7. `EntityCallScene(ENTITY_ID, args)`: Similar to `EntityCallAction`, this method is used to call or trigger a scene associated with the entity (`ENTITY_ID`) and provides any required arguments (`args`) for that scene.

8. `EntityGetSettings(ENTITY_ID)`: This method retrieves the settings or configuration for the specified entity (`ENTITY_ID`). It allows you to access and query the settings associated with an entity.

9. `EntityGetAttributes(ENTITY_ID)`: This method retrieves all attributes associated with a specific entity (`ENTITY_ID`). It provides a way to obtain detailed information about the entity's attributes.

The methods of the "Entity" object provide convenient capabilities for managing the device's state, attributes, metrics, and executing actions on the device. You can use these methods to create device control logic in your **Smart Home** project.

